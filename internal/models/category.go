package models

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type CategoryType string

const (
	Income   CategoryType = "income"
	Expense  CategoryType = "expense"
	Filename              = "./data/categories.json"
)

type Category struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	IsIncome bool   `json:"is_income"`
	Edit     bool   `json:"edit"`
}

var (
	DefaultExpenseCategories = []Category{
		{ID: "1", Name: "★Продукты", IsIncome: false, Type: "expense"},
		{ID: "2", Name: "★Транспорт", IsIncome: false, Type: "expense"},
		{ID: "3", Name: "★Жилье", IsIncome: false, Type: "expense"},
		{ID: "4", Name: "★Развлечения", IsIncome: false, Type: "expense"},
	}

	DefaultIncomeCategories = []Category{
		{ID: "5", Name: "★Зарплата", IsIncome: true, Type: "income"},
		{ID: "6", Name: "★Подарки", IsIncome: true, Type: "income"},
		{ID: "7", Name: "★Прочие доходы", IsIncome: true, Type: "income"},
	}

	AllCategories []Category
)

func loadFromFile() {
	file, err := os.Open(Filename)
	if err != nil {
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&DefaultExpenseCategories)
	if err != nil {
		fmt.Println("Ошибка загрузки данных, начнём заново.")
		DefaultExpenseCategories = []Category{}
	}

	decoder = json.NewDecoder(file)
	err = decoder.Decode(&DefaultIncomeCategories)
	if err != nil {
		fmt.Println("Ошибка загрузки данных, начнём заново.")
		DefaultIncomeCategories = []Category{}
	}
}

func saveToFile() {
	file, err := os.Create(Filename)
	if err != nil {
		fmt.Println("Не могу сохранить файл.")
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(DefaultExpenseCategories)
	if err != nil {
		fmt.Println("Ошибка при сохранении.")
	}

	encoder = json.NewEncoder(file)
	err = encoder.Encode(DefaultIncomeCategories)
	if err != nil {
		fmt.Println("Ошибка при сохранении.")
	}
}

func GetDefaultCategories() ([]Category, error) {

	fileInfo, err := os.ReadFile("./data/categories.json")

	if err != nil {
		fmt.Printf("Возникла ошабка при инициализации базы данных\n%s", err)
	}

	if fileInfo == nil {
		for _, v := range DefaultExpenseCategories {
			fmt.Printf("ID: %s Название: %s  Тип транзакции: %s  Подлежит редактированию: %v\n", v.ID, v.Name, v.IsIncome, v.Edit)
		}
		for _, v := range DefaultIncomeCategories {
			fmt.Printf("ID: %s Название: %s  Тип транзакции: %s  Подлежит редактированию: %v\n", v.ID, v.Name, v.IsIncome, v.Edit)
		}

	}
	var inf []Category

	if err := json.Unmarshal(fileInfo, &inf); err != nil {
		fmt.Printf("Ошибка парсинга JSON: %v\n", err)
		return nil, err
	}

	for _, v := range inf {
		fmt.Printf("ID: %s Название: %s  Тип транзакции: %s  Подлежит редактированию: %v\n", v.ID, v.Name, v.IsIncome, v.Edit)
	}

	return inf, nil
}

func FindCategoryByID(id string) (*Category, error) {
	defaultCategories, err := GetDefaultCategories()

	if err != nil {
		fmt.Print("Не получилось запарсить JSON")
		return nil, err
	}

	var intID int
	_, err = fmt.Sscanf(id, "%d", &intID)
	if err != nil {
		return nil, fmt.Errorf("неверный формат ID: %v", err)
	}

	for _, category := range defaultCategories {
		if category.ID == id {
			return &category, nil
		}
	}

	return nil, fmt.Errorf("не получилось найти '%s' ", id)

}

func IsEditCategory(id string) (bool, string) {
	var intID int
	_, err := fmt.Sscanf(id, "%d", &intID)
	if err != nil {
		return false, "Некорректный формат ID"
	}

	for _, cat := range AllCategories {
		if cat.ID == id {
			return cat.Edit, "Системная категория"
		}
	}
	return false, "Не системная категория"

}

func GetCategoriesByType(typ bool) string {

	if typ {
		return string(Income)
	} else {
		return string(Expense)
	}

}

func GetCategoryTypeByName(name string) (string, bool) {
	nameLower := strings.ToLower(strings.TrimSpace(name))

	expenseCategories := []string{"продукты", "транспорт", "жилье", "развлечения"}
	for _, categoryName := range expenseCategories {
		if nameLower == categoryName {
			return "Expense", true
		}
	}

	incomeCategories := []string{"зарплата", "подарки", "прочие доходы"}
	for _, categoryName := range incomeCategories {
		if nameLower == categoryName {
			return "Income", true
		}
	}

	return "", false

}

func ValidateCategory(category *Category) error {
	category.Name = strings.TrimSpace(category.Name)

	if len(category.Name) == 0 {
		return fmt.Errorf("название не должно быть пустым")
	}

	if CategoryType(category.Type) != Expense && CategoryType(category.Type) != Income {
		return fmt.Errorf("некорректный тип категории: должен быть Income или Expense")
	}

	if category.Edit == true {
		return fmt.Errorf("системные категории нельзя изменять")
	}

	return nil
}

func CategoryExists(categories *Category, name string) bool {
	trimCat := strings.ToLower(strings.TrimSpace(categories.Name))
	trimName := strings.ToLower(strings.TrimSpace(name))

	for _, catName := range trimCat {
		if trimName == string(catName) {
			return true
		}
	}
	return false
}
