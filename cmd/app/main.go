package main

import (
	"bufio"
	"fintrack/internal/models"
	"fintrack/internal/services"
	"fintrack/internal/storage"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	ColorRest   = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
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
		fmt.Printf(ColorYellow+"Предупреждение при загрузке категорий: %w"+ColorRest, err)
	}

	scanner := bufio.NewScanner(os.Stdin)

	return &App{
		transactionService: transactionService,
		categoryService:    (*services.CategoryService)(categoryService),
		scanner:            scanner,
	}

}

func clearScreen() {

}

func (app *App) DisplayMainMenu() {
	clearScreen()

	fmt.Println(ColorBlue + "=====================FinTrack=====================" + ColorRest)
	fmt.Println("1.Добавить транзакцию")
	fmt.Println("2.Показать транзакции")
	fmt.Println("3.Показать категории")
	fmt.Println("0.Выход")
	fmt.Println(ColorCyan + "==================================================" + ColorRest)

}

func (app *App) addTransaction() error {
	clearScreen()
	fmt.Println(ColorCyan + "==============Добавление транзакции===============" + ColorRest)
	fmt.Print(ColorCyan + "Введите сумму: " + ColorRest)
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

	fmt.Print(ColorCyan + "Тип транзакции (1-доход 2-расход):" + ColorRest)

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

	fmt.Println(ColorYellow + "Доступные категории: " + ColorRest)

	for i, category := range categories {
		fmt.Printf("%d.%s\n", i+1, category.Name)
	}

	fmt.Print(ColorCyan + "Выберите категорию(номер): " + ColorRest)

	if !app.scanner.Scan() {
		return fmt.Errorf("ошибка чтения категории")
	}

	categoryindex, err := strconv.Atoi(strings.TrimSpace(app.scanner.Text()))

	if err != nil || categoryindex < 1 || categoryindex > len(categories) {
		return fmt.Errorf("неверный номер категории. Выберите от 1 до %d", len(categories))
	}

	selectedCategory := categories[categoryindex-1].Name

	fmt.Print(ColorCyan + "Введите описание: " + ColorRest)

	if !app.scanner.Scan() {
		return fmt.Errorf("ошибка чтения описания")
	}

	descripyion := strings.TrimSpace(app.scanner.Text())

	if descripyion == "" {
		descripyion = "Без описания"
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

	fmt.Println(ColorGreen + "\n Транзакция успешно добавлена!" + ColorRest)
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
	fmt.Println(ColorBlue + "==================== Транзакции ====================" + ColorRest)

	transactions, err := app.transactionService.GetAllTransactions()
	if err != nil {
		return fmt.Errorf("ошибка при получении транзакций: %v", err)
	}

	if len(transactions) == 0 {
		fmt.Println(ColorYellow + "Нет доступных транзакций для отображения." + ColorRest)
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

		fmt.Printf("%-10s | %s%-12.2f%s | %-15s | %-20s | %s\n",
			displayID,
			amountColor, t.Amount, ColorRest,
			t.Category,
			t.Date.Format("02.01.2006 15:04"),
			transactionType)

	}

	fmt.Println(strings.Repeat("-", 70))

	balance := totalIncome - totalExpense
	balanceColor := ColorGreen
	if balance < 0 {
		balanceColor = ColorRed
	}

	fmt.Printf(ColorCyan+"Итоговый доход: %s%.2f%s\n"+ColorRest, ColorGreen, totalIncome, ColorRest)
	fmt.Printf(ColorCyan+"Итоговый расход: %s%.2f%s\n"+ColorRest, ColorRed, totalExpense, ColorRest)
	fmt.Printf(ColorCyan+"Баланс: %s%.2f%s\n"+ColorRest, balanceColor, balance, ColorRest)

	return nil

}

func (app *App) showCategories() error {

	clearScreen()

	fmt.Println(ColorBlue + "==================== Категории ====================" + ColorRest)

	incomeCategories, err := app.categoryService.GetCategoriesByType(true)

	if err != nil {
		return fmt.Errorf("ошибка загрузки категорий доходов: %v", err)
	}

	fmt.Println(ColorGreen + "\nДоходы:" + ColorRest)

	if len(incomeCategories) == 0 {
		fmt.Println(ColorYellow + "Нет доступных категорий доходов." + ColorRest)
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

	// Добавляем категории
	err := categoryService.AddCategory("Продукты", false)
	if err != nil && err.Error() != "" {
		log.Printf("Ошибка добавления категории: %v\n", err)
	}

	// Добавляем транзакцию
	err = transactionService.AddTransaction(100.50, "Продукты", "Покупка в магазине", "expense")
	if err != nil {
		log.Printf("Ошибка добавления транзакции: %v\n", err)
	}

	// Получаем все транзакции
	transactions, err := transactionService.GetAllTransactions()
	if err != nil {
		log.Fatalf("Ошибка получения транзакций: %v\n", err)
	}

	fmt.Println("=== Все транзакции ===")
	for _, t := range transactions {
		fmt.Printf("ID: %s, Сумма: %.2f, Категория: %s, Тип: %s, Дата: %s\n",
			t.ID, t.Amount, t.Category, t.Type, t.Date.Format("02.01.2006 15:04"))
	}

	// Получаем категории по типу
	categories, err := categoryService.GetCategoriesByType(false)
	if err != nil {
		log.Fatalf("Ошибка получения категорий: %v\n", err)
	}

	fmt.Println("\n=== Категории расходов ===")
	for _, c := range categories {
		fmt.Printf("Название: %s, Тип: %s\n", c.Name, c.Type)
	}
}

func main() {
	app := NewApp()

	if app == nil {
		fmt.Println(ColorRed + "Ошибка при инициализации приложения." + ColorRest)
		return
	}

	clearScreen()
	fmt.Println(ColorGreen + "╔════════════════════════════════════════════════════════╗" + ColorRest)
	fmt.Println(ColorGreen + "║           Добро пожаловать в FinTrack!                 ║" + ColorRest)
	fmt.Println(ColorGreen + "║    Простой и надежный финансовый трекер на Go          ║" + ColorRest)
	fmt.Println(ColorGreen + "╚════════════════════════════════════════════════════════╝" + ColorRest)
	fmt.Print(ColorCyan + "\nНажмите Enter для начала работы..." + ColorRest)
	app.scanner.Scan()

	for {
		app.DisplayMainMenu()

		fmt.Print(ColorCyan + "Выберите опцию: " + ColorRest)

		if !app.scanner.Scan() {
			fmt.Println(ColorYellow + "\n\nОбнаружен сигнал завершения. Завершение работы..." + ColorRest)
			break
		}

		choiceStr := strings.TrimSpace(app.scanner.Text())

		if choiceStr == "" {
			continue
		}

		choice, err := strconv.Atoi(choiceStr)
		if err != nil {
			fmt.Println(ColorRed + "Некорректный ввод. Пожалуйста, введите число." + ColorRest)
			continue
		}

		switch choice {
		case 1:
			err := app.addTransaction()
			if err != nil {
				fmt.Println(ColorRed + "Ошибка при добавлении транзакции: " + err.Error() + ColorRest)
			}
		case 2:
			err := app.showTransactions()
			if err != nil {
				fmt.Println(ColorRed + "Ошибка при отображении транзакций: " + err.Error() + ColorRest)
			}
		case 3:
			err := app.showCategories()
			if err != nil {
				fmt.Println(ColorRed + "Ошибка при отображении категорий: " + err.Error() + ColorRest)
			}
		case 0:
			clearScreen()
			fmt.Println(ColorGreen + "╔════════════════════════════════════════════════════════╗" + ColorRest)
			fmt.Println(ColorGreen + "║        Спасибо за использование FinTrack!              ║" + ColorRest)
			fmt.Println(ColorGreen + "║                До свидания!                            ║" + ColorRest)
			fmt.Println(ColorGreen + "╚════════════════════════════════════════════════════════╝" + ColorRest)
			return
		default:
			fmt.Println(ColorRed + "\nНекорректный выбор. Пожалуйста, выберите опцию от 1 до 4." + ColorRest)
		}

		waitForEnter(app.scanner)

	}

}

func waitForEnter(scanner *bufio.Scanner) {
	fmt.Print(ColorWhite + "\nНажмите Enter для продолжения..." + ColorRest)
	scanner.Scan()
}
