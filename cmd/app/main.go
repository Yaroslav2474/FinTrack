package main

import (
	"bufio"
	"fintrack/internal/models"
	"fintrack/internal/services"
	"fintrack/internal/storage"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

const (
	Rest   = ""
	Red    = lipgloss.Color("#ff0000")
	Green  = lipgloss.Color("#00ff0d")
	Yellow = lipgloss.Color("#fffb00")
	Blue   = lipgloss.Color("#0000FF")
	Cyan   = lipgloss.Color("#3bdddd")
	White  = lipgloss.Color("#FFFFFF")
)

var (
	ColorRed    = lipgloss.NewStyle().Foreground(Red)
	ColorGreen  = lipgloss.NewStyle().Foreground(Green)
	ColorYellow = lipgloss.NewStyle().Foreground(Yellow)
	ColorBlue   = lipgloss.NewStyle().Foreground(Blue)
	ColorCyan   = lipgloss.NewStyle().Foreground(Cyan)
	ColorWhite  = lipgloss.NewStyle().Foreground(White)
)

type App struct {
	transactionService *services.TransactionService
	categoryService    *services.CategoryService
	scanner            *bufio.Scanner
}

func generateUniqueID() string {
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("tx_%d", timestamp%10000000000)
}

func NewApp() *App {
	fileStorage := storage.NewFileStorage("data", "transactions")

	transactionService := services.NewTransactionService(fileStorage)
	categoryService := services.NewCategoryService(fileStorage)

	inf, err := models.GetDefaultCategories()

	fmt.Println(inf)

	if err != nil {
		fmt.Printf("%s %s", ColorYellow.Render("Предупреждение при загрузке категорий: "), ColorYellow.Render(fmt.Sprintf("%v", err)))
	}

	scanner := bufio.NewScanner(os.Stdin)

	return &App{
		transactionService: transactionService,
		categoryService:    (*services.CategoryService)(categoryService),
		scanner:            scanner,
	}

}

func clearScreen() {

	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

}

func (app *App) DisplayMainMenu() {
	clearScreen()

	fmt.Printf("%s\n", ColorCyan.Render("=====================FinTrack====================="))
	fmt.Printf("%s\n", ColorWhite.Render("1. Добавить транзакции"))
	fmt.Printf("%s\n", ColorWhite.Render("2.Показать транзакции"))
	fmt.Printf("%s\n", ColorWhite.Render("3.Показать категории"))
	fmt.Printf("%s\n", ColorWhite.Render("0.Выход"))
	fmt.Printf("%s\n", ColorCyan.Render("=================================================="))

}

func (app *App) addTransaction() error {
	clearScreen()
	fmt.Println(ColorBlue.Render("==============Добавление транзакции==============="))
	fmt.Print(ColorCyan.Render("Введите сумму: "))
	if !app.scanner.Scan() {
		return fmt.Errorf("ошибка чтения суммы")
	}

	amountstr := strings.TrimSpace(app.scanner.Text())
	amount, err := strconv.ParseFloat(amountstr, 64)

	if err != nil {
		return fmt.Errorf("ошибка при вводе суммы")
	}

	if amount <= 0 {
		return fmt.Errorf("сумма должна быть положительной")
	}

	fmt.Print(ColorCyan.Render("\nТип транзакции (1-доход 2-расход):"))

	if !app.scanner.Scan() {
		return fmt.Errorf("ошибка чтения типа транзакции")
	}

	var isincome bool

	switch strings.TrimSpace(app.scanner.Text()) {
	case "1":
		isincome = true
	case "2":
		isincome = false
	default:
		return fmt.Errorf("неверный выбор типа транзакции. Выберите 1 или 2")
	}

	categories, err := app.categoryService.GetCategoriesByType(isincome)
	if err != nil {
		return fmt.Errorf("ошибка получения категорий: %v", err)
	}

	if len(categories) == 0 {
		return fmt.Errorf("нет доступных категорий для выбранного типа")
	}

	fmt.Println(ColorCyan.Render("Доступные категории: "))

	for i, category := range categories {
		fmt.Printf("\n%d.%s\n", i+1, category.Name)
	}

	fmt.Print(ColorCyan.Render("\nВыберите категорию(номер): "))

	if !app.scanner.Scan() {
		return fmt.Errorf("ошибка чтения категории")
	}

	categoryindex, err := strconv.Atoi(strings.TrimSpace(app.scanner.Text()))

	if err != nil || categoryindex < 1 || categoryindex > len(categories) {
		return fmt.Errorf("неверный номер категории. Выберите от 1 до %d", len(categories))
	}

	selectedCategory := categories[categoryindex-1].Name

	fmt.Print(ColorCyan.Render("\nВведите описание: "))

	if !app.scanner.Scan() {
		return fmt.Errorf("ошибка чтения описания")
	}

	descripyion := strings.TrimSpace(app.scanner.Text())

	if descripyion == "" {
		descripyion = "\nБез описания"
	}

	transactionType := "income"
	if !isincome {
		transactionType = "expense"
	}

	transaction := &models.Transaction{
		ID:          generateUniqueID(),
		Amount:      amount,
		Category:    selectedCategory,
		Description: descripyion,
		Type:        models.TransactionType(transactionType),
		Date:        time.Now(),
	}

	if err := app.transactionService.AddTransaction(amount, selectedCategory, descripyion, transactionType); err != nil {
		return fmt.Errorf("ошибка при добавлении транзакции: %v", err)
	}

	transactionTypeDisplay := "Расход"
	if isincome {
		transactionTypeDisplay = "Доход"
	}

	fmt.Println(ColorGreen.Render("\n Транзакция успешно добавлена!\n"))
	fmt.Printf("ID: %s\nСумма: %.2f\nТип: %s\nКатегория: %s\nОписание: %s\nДата: %s\n",
		transaction.ID,
		transaction.Amount,
		transactionTypeDisplay,
		transaction.Category,
		transaction.Description,
		transaction.Date.Format("02.01.2006 15:04:05"),
	)

	return nil
}

func (app *App) showTransactions() error {
	clearScreen()
	fmt.Println(ColorBlue.Render("==================== Транзакции ===================="))

	transactions, err := app.transactionService.GetAllTransactions()
	if err != nil {
		return fmt.Errorf("ошибка при получении транзакций: %v", err)
	}

	if len(transactions) == 0 {
		fmt.Println(ColorYellow.Render("Нет доступных транзакций для отображения."))
		return nil
	}

	fmt.Println()
	fmt.Printf("%-20s %-15s %-10s %-20s %-30s %-20s\n", "ID", "Сумма", "Тип", "Категория", "Описание", "Дата")
	fmt.Println(strings.Repeat("-", 70))

	totalIncome := 0.0
	totalExpense := 0.0

	for _, t := range transactions {
		amountColor := ColorGreen
		transactionType := "Доход"
		if t.Type == "expense" {
			amountColor = ColorRed
			transactionType = "Расход"
			totalExpense += t.Amount
		} else {
			totalIncome += t.Amount
		}

		displayID := t.ID
		if len(displayID) > 10 {
			displayID = t.ID[:7] + "..."
		}

		fmt.Printf("%-10s | %s%-12.2f | %-15s | %-20s | %s\n",
			displayID,
			amountColor, t.Amount, t.Category,
			t.Date.Format("02.01.2006 15:04"),
			transactionType)

	}

	fmt.Println(strings.Repeat("-", 70))

	balance := totalIncome - totalExpense
	balanceColor := lipgloss.NewStyle().Foreground(Green)
	if balance < 0 {
		balanceColor = lipgloss.NewStyle().Foreground(Red)
	}

	fmt.Printf("%s %s", ColorCyan.Render("\nИтоговый доход: "), ColorCyan.Render(fmt.Sprintf("%.2f", totalIncome)))
	fmt.Printf("%s %s", ColorCyan.Render("\nИтоговый расход: "), ColorCyan.Render(fmt.Sprintf("%.2f", totalExpense)))
	fmt.Printf("%s %s", balanceColor.Render("\nБаланс: "), ColorCyan.Render(fmt.Sprintf("%.2f", balance)))

	return nil

}

func (app *App) showCategories() error {

	clearScreen()

	fmt.Println(ColorBlue.Render("==================== Категории ===================="))

	incomeCategories, err := app.categoryService.GetCategoriesByType(true)

	if err != nil {
		return fmt.Errorf("ошибка загрузки категорий доходов: %v", err)
	}

	fmt.Println(ColorCyan.Render("Доходы:"))

	if len(incomeCategories) == 0 {
		fmt.Println(ColorYellow.Render("Нет доступных категорий доходов."))
	} else {
		for _, categories := range incomeCategories {
			fmt.Printf("  • %s — %s\n", categories.Name, categories.Type)
		}

	}

	return nil

}

func init() {
	fs := storage.NewFileStorage(
		"internal/data/transactions.json",
		"internal/data/categories.json",
	)

	transactionService := services.NewTransactionService(fs)
	categoryService := services.NewCategoryService(fs)

	err := categoryService.AddCategory("Продукты", false)
	if err != nil && err.Error() != "" {
		log.Printf("\nОшибка добавления категории: %v\n", err)
	}

	err = transactionService.AddTransaction(100.50, "Продукты", "Покупка в магазине", "expense")
	if err != nil {
		log.Printf("\nОшибка добавления транзакции: %v\n", err)
	}

	transactions, err := transactionService.GetAllTransactions()
	if err != nil {
		log.Fatalf("\nОшибка получения транзакций: %v\n", err)
	}

	fmt.Println("\n=== Все транзакции ===\n")
	for _, t := range transactions {
		fmt.Printf("ID: %s, Сумма: %.2f, Категория: %s, Тип: %s, Дата: %s\n",
			t.ID, t.Amount, t.Category, t.Type, t.Date.Format("02.01.2006 15:04"))
	}

	categories, err := categoryService.GetCategoriesByType(false)
	if err != nil {
		log.Fatalf("Ошибка получения категорий: %v\n", err)
	}

	fmt.Println("\n=== Категории расходов ===\n")
	for _, c := range categories {
		fmt.Printf("Название: %s, Тип: %s\n", c.Name, c.Type)
	}
}

func main() {
	app := NewApp()

	if app == nil {
		fmt.Println(ColorRed.Render("Ошибка при инициализации приложения."))
		return
	}

	clearScreen()
	fmt.Println(ColorGreen.Render("╔════════════════════════════════════════════════════════╗"))
	fmt.Println(ColorGreen.Render("║           Добро пожаловать в FinTrack!                 ║"))
	fmt.Println(ColorGreen.Render("║    Простой и надежный финансовый трекер на Go          ║"))
	fmt.Println(ColorGreen.Render("╚════════════════════════════════════════════════════════╝"))
	fmt.Print(ColorCyan.Render("\nНажмите Enter для начала работы..."))
	app.scanner.Scan()

	for {
		app.DisplayMainMenu()

		fmt.Print(ColorCyan.Render("\nВыберите опцию: "))

		if !app.scanner.Scan() {
			fmt.Println(ColorYellow.Render("\n\nОбнаружен сигнал завершения. Завершение работы..."))
			break
		}

		choiceStr := strings.TrimSpace(app.scanner.Text())

		if choiceStr == "" {
			continue
		}

		choice, err := strconv.Atoi(choiceStr)
		if err != nil {
			fmt.Println(ColorRed.Render("Некорректный ввод. Пожалуйста, введите число."))
			continue
		}

		switch choice {
		case 1:
			err := app.addTransaction()
			if err != nil {
				fmt.Println(ColorRed.Render("Ошибка при добавлении транзакции: " + err.Error()))
			}
		case 2:
			err := app.showTransactions()
			if err != nil {
				fmt.Println(ColorRed.Render("Ошибка при отображении транзакций: " + err.Error()))
			}
		case 3:
			err := app.showCategories()
			if err != nil {
				fmt.Println(ColorRed.Render("Ошибка при отображении категорий: " + err.Error()))
			}
		case 0:
			clearScreen()
			fmt.Println(ColorGreen.Render("╔════════════════════════════════════════════════════════╗"))
			fmt.Println(ColorGreen.Render("║        Спасибо за использование FinTrack!              ║"))
			fmt.Println(ColorGreen.Render("║                До свидания!                            ║"))
			fmt.Println(ColorGreen.Render("╚════════════════════════════════════════════════════════╝"))
			time.NewTimer(3 * time.Second)
			return
		default:
			fmt.Println(ColorRed.Render("\nНекорректный выбор. Пожалуйста, выберите опцию от 1 до 4."))
		}

		waitForEnter(app.scanner)

	}

}

func waitForEnter(scanner *bufio.Scanner) {
	fmt.Print(ColorWhite.Render("\nНажмите Enter для продолжения..."))
	scanner.Scan()
}
