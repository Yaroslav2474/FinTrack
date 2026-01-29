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
	Filename              = "internal/data/categories.json"
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
	// legacy function kept for compatibility — not used by GetDefaultCategories
	// but left here in case of direct calls elsewhere.
	file, err := os.Open(Filename)
	if err != nil {
		return
	}
	defer file.Close()

	var cats []Category
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cats); err != nil {
		return
	}
	// split into expense/income defaults if needed
	DefaultExpenseCategories = []Category{}
	DefaultIncomeCategories = []Category{}
	for _, c := range cats {
		if c.IsIncome {
			DefaultIncomeCategories = append(DefaultIncomeCategories, c)
		} else {
			DefaultExpenseCategories = append(DefaultExpenseCategories, c)
		}
	}
}

func saveToFile() {
	// Save combined categories as a single JSON array
	file, err := os.Create(Filename)
	if err != nil {
		fmt.Println("Не могу сохранить файл.")
		return
	}
	defer file.Close()

	var all []Category
	all = append(all, DefaultExpenseCategories...)
	all = append(all, DefaultIncomeCategories...)

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(all); err != nil {
		fmt.Println("Ошибка при сохранении.")
	}
}

func GetDefaultCategories() ([]Category, error) {

	data, err := os.ReadFile(Filename)
	if err != nil {
		// If the file doesn't exist or cannot be read, return default categories without error
		var defaults []Category
		defaults = append(defaults, DefaultExpenseCategories...)
		defaults = append(defaults, DefaultIncomeCategories...)
		return defaults, nil
	}

	var cats []Category
	if err := json.Unmarshal(data, &cats); err != nil {
		// On parse error, return defaults
		var defaults []Category
		defaults = append(defaults, DefaultExpenseCategories...)
		defaults = append(defaults, DefaultIncomeCategories...)
		return defaults, nil
	}

	return cats, nil
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
