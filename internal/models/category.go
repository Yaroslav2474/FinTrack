package models

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type сategoryType string

const (
	Income   сategoryType = "income"
	Expense  сategoryType = "expense"
	filename              = "./data/categories.json"
)

type Category struct {
	ID          string
	Title       string
	Type        сategoryType
	Edit        bool
	Description string
	// Budget float64
}

var (
	defaultExpenseCategories = []Category{
		{ID: "sys_exp_food", Title: "★Продукты", Type: Expense, Edit: true},
		{ID: "sys_exp_transport", Title: "★Транспорт", Type: Expense, Edit: true},
		{ID: "sys_exp_housing", Title: "★Жилье", Type: Expense, Edit: true},
		{ID: "sys_exp_entertainment", Title: "★Развлечения", Type: Expense, Edit: true},
	}

	defaultIncomeCategories = []Category{
		{ID: "sys_inc_salary", Title: "★Зарплата", Type: Income, Edit: true},
		{ID: "sys_inc_gifts", Title: "★Подарки", Type: Income, Edit: true},
		{ID: "sys_inc_other", Title: "★Прочие доходы", Type: Income, Edit: true},
	}

	allCategories []Category
)

func loadFromFile() {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&defaultExpenseCategories)
	if err != nil {
		fmt.Println("Ошибка загрузки данных, начнём заново.")
		defaultExpenseCategories = []Category{}
	}

	decoder = json.NewDecoder(file)
	err = decoder.Decode(&defaultIncomeCategories)
	if err != nil {
		fmt.Println("Ошибка загрузки данных, начнём заново.")
		defaultIncomeCategories = []Category{}
	}
}

func saveToFile() {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Не могу сохранить файл.")
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(defaultExpenseCategories)
	if err != nil {
		fmt.Println("Ошибка при сохранении.")
	}

	encoder = json.NewEncoder(file)
	err = encoder.Encode(defaultIncomeCategories)
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
		for _, v := range defaultExpenseCategories {
			fmt.Printf("ID: %s Название: %s  Тип транзакции: %s  Подлежит редактированию: %v", v.ID, v.Title, v.Type, v.Edit)
		}
		for _, v := range defaultIncomeCategories {
			fmt.Printf("ID: %s Название: %s  Тип транзакции: %s  Подлежит редактированию: %v", v.ID, v.Title, v.Type, v.Edit)
		}

	}
	var inf []Category

	if err := json.Unmarshal(fileInfo, &inf); err != nil {
		fmt.Printf("Ошибка парсинга JSON: %v\n", err)
		return nil, err
	}

	for _, v := range inf {
		fmt.Printf("ID: %s Название: %s  Тип транзакции: %s  Подлежит редактированию: %v", v.ID, v.Title, v.Type, v.Edit)
	}

	return inf, nil
}

func FindCategoryByID(id string) (*Category, error) {
	defaultCategories, err := GetDefaultCategories()

	if err != nil {
		fmt.Print("Не получилось запарсить JSON")
		return nil, err
	}

	for _, category := range defaultCategories {
		if category.ID == id {
			return &category, nil
		}
	}

	return nil, fmt.Errorf("не получилось найти '%s' ", id)

}

func IsEditCategory(id string) (bool, string) {

	for _, cat := range allCategories {
		if cat.ID == id {
			return cat.Edit, "Системная категория"
		}
	}
	return false, "Не системная категория"

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
	category.Title = strings.TrimSpace(category.Title)

	if len(category.Title) == 0 {
		return fmt.Errorf("название не должно быть пустым")
	}

	if category.Type != Expense || category.Type != Income {
		return fmt.Errorf("некорректный тип категории: должен быть Income или Expense")
	}

	if category.Edit == true {
		return fmt.Errorf("системные категории нельзя изменять")
	}

	return nil
}

func CategoryExists(categories *Category, name string) bool {
	trimCat := strings.ToLower(strings.TrimSpace(categories.Title))
	trimName := strings.ToLower(strings.TrimSpace(name))

	for _, catName := range trimCat {
		if trimName == string(catName) {
			return true
		}
	}
	return false
}
