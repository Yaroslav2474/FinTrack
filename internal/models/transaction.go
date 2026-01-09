package models

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type TransactionType string

var Pending, Completed, Canceled TransactionType

type Transaction struct {
	ID          string
	Amount      float64
	Type        TransactionType
	Category    string
	Date        time.Time
	Description string
}

func (tran Transaction) IsValid() (bool, error) {

	if tran.Amount <= 0 {
		return false, fmt.Errorf("сумма должна быть положительной")
	}

	if tran.Type != TransactionType(Income) || tran.Type != TransactionType(Expense) {
		return false, fmt.Errorf("неизветсный тип транзакции")
	}

	if len(tran.Category) == 0 {
		return false, fmt.Errorf("категория не может быть пустой")
	}

	if len(tran.Description) == 0 {
		return false, fmt.Errorf("описание не должно быть пустым")
	}

	return true, fmt.Errorf("")
}

func (tran Transaction) GetFormattedAmount() string {

	wholeAmount := int(tran.Amount)
	fractionAmount := tran.Amount - float64(wholeAmount)

	formattedAmount := ""
	// temp := strconv.Itoa(wholeAmount)

	correctAmount := formattedAmount + "." + fmt.Sprintf("%02d", int(fractionAmount*100))

	if tran.Type == TransactionType(Expense) {
		correctAmount = "(" + correctAmount + " ₽" + ")"
		return correctAmount
	} else {
		correctAmount = "+" + correctAmount + " ₽"
		return correctAmount
	}

}

func (tran Transaction) GetDisplayDate() (int, error) {
	return fmt.Printf("%d.%d.%dг.", tran.Date.Day(), tran.Date.Month(), tran.Date.Year())
}

func NewTransaction(amount float64, tType TransactionType, category, description string) (*Transaction, error) {
	rand.Seed(time.Now().UnixNano())
	newID := strconv.Itoa(rand.Intn(9000000) + 1000000)

	if category == string(Expense) {
		if amount < 0 {
			amount = math.Round(amount*100) / 100
		} else {
			return &Transaction{Category: " ", ID: " ", Amount: 0.0, Type: " ", Date: time.Now(), Description: " "}, fmt.Errorf("ошибка при введении суммы, она должна быть больше 0")
		}
	} else {
		if amount > 0 {
			amount = math.Round(amount*100) / 100
		} else {
			return &Transaction{Category: " ", ID: " ", Amount: 0.0, Type: " ", Date: time.Now(), Description: " "}, fmt.Errorf("ошибка при введении суммы, она должна быть больше 0")
		}
	}

	if tType != "income" || tType != "expense" {
		return &Transaction{Category: " ", ID: " ", Amount: 0.0, Type: " ", Date: time.Now(), Description: " "}, fmt.Errorf("ошибка при введении типа транзакции")
	}

	if len(category) == 0 {
		return &Transaction{Category: " ", ID: " ", Amount: 0.0, Type: " ", Date: time.Now(), Description: " "}, fmt.Errorf("ошибка при введении категории")
	} else {
		category = strings.ToLower(strings.TrimSpace(category))
	}

	if len(description) >= 256 {
		return &Transaction{Category: " ", ID: " ", Amount: 0.0, Type: " ", Date: time.Now(), Description: " "}, fmt.Errorf("ошибка при введении описания")
	}

	return &Transaction{ID: newID, Amount: amount, Type: tType, Category: category, Date: time.Now(), Description: description}, fmt.Errorf("")
}

func ParseTransactionType(s string) (TransactionType, error) {
	userType := strings.TrimSpace(strings.ToLower(s))

	switch userType {
	case "income":
		return TransactionType(Income), fmt.Errorf("")
	case "expense":
		return TransactionType(Expense), fmt.Errorf("")
	default:
		return TransactionType(""), fmt.Errorf("такого типа не существует")
	}
}

func AddTrannsaction() {
	var choise int

	fmt.Print("----------------------Меню----------------------")
	fmt.Print("1. Добавить транзакцию\n2. Показать транзакции\n0. Выход")

	fmt.Scan(&choise)

	for {
		if choise == 1 {

			fmt.Println("Выберете вид транзакции: \n1.Доходы\n2.Расходы")
			fmt.Scan(&choise)

			switch choise {
			case 1:

			case 2:

			default:
				fmt.Print("Только два вида транзакций.")
			}

		} else if choise == 2 {
			GetDefaultCategories()
		} else if choise == 0 {
			fmt.Println("До свидания!")
			break
		} else {
			fmt.Println("Только 3 действия в меню.")
		}
	}
}
