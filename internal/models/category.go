package models

type сategoryType string

const (
	Income  сategoryType = "income"
	Expense сategoryType = "expense"
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
		{ID: "sys_exp_food", Title: "Продукты", Type: Expense, Edit: true},
		{ID: "sys_exp_transport", Title: "Транспорт", Type: Expense, Edit: true},
		{ID: "sys_exp_housing", Title: "Жилье", Type: Expense, Edit: true},
		{ID: "sys_exp_entertainment", Title: "Развлечения", Type: Expense, Edit: true},
	}

	defaultIncomeCategories = []Category{
		{ID: "sys_inc_salary", Title: "Зарплата", Type: Income, Edit: true},
		{ID: "sys_inc_gifts", Title: "Подарки", Type: Income, Edit: true},
		{ID: "sys_inc_other", Title: "Прочие доходы", Type: Income, Edit: true},
	}
)
