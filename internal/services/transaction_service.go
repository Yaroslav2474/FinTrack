package services

import (
	"fintrack/internal/models"
	"fintrack/internal/storage"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type TransactionService struct {
	storage storage.Storage
}

func NewTransactionService(storage storage.Storage) *TransactionService {
	return &TransactionService{
		storage: storage,
	}
}

func generateUniqueID() string {
	rand.Seed(time.Now().UnixNano())
	newID := strconv.Itoa(rand.Intn(9000000) + 1000000)
	return newID
}

func validateTransaction(transaction models.Transaction) error {
	if transaction.Amount <= 0 {
		return fmt.Errorf("сумма не может быть < 0")
	}

	if transaction.Category == "" {
		return fmt.Errorf("категория не может быть пустой")
	}

	if transaction.Description == "" {
		return fmt.Errorf("описание не может быть пустым")
	}

	return nil

}

func (ts TransactionService) AddTransaction(amount float64, category string, description string, transactionType string) error {

	newTransaction := models.Transaction{
		ID:          generateUniqueID(),
		Amount:      amount,
		Category:    category,
		Description: description,
		Type:        models.TransactionType(transactionType),
		Date:        time.Now(),
	}

	if amount <= 0 {
		return fmt.Errorf("сумма не может быть < 0")
	}

	if strings.TrimSpace(category) == "" {
		return fmt.Errorf("категория не может быть пустой")
	}

	if transactionType != string(models.Expense) || transactionType != string(models.Income) {
		return fmt.Errorf("неизвестный тип транзакции: %s. Должен быть 'income' или 'expence'", transactionType)
	}

	categories, err := ts.storage.GetCategories()

	if err != nil {
		return fmt.Errorf("ошибка в получении категорий")
	}

	categoryFound := false
	autoDeterminedType := ""

	for _, cat := range categories {

		if strings.EqualFold(cat.Title, category) {
			categoryFound = true

			if transactionType == "" {
				autoDeterminedType = string(cat.Type)
			} else {

				if string(cat.Type) != transactionType {
					return fmt.Errorf("несоответствие типа категории")
				}
			}
			break

		}

	}

	if !categoryFound {
		return fmt.Errorf("категория не найдена")
	}

	if autoDeterminedType != "" {
		newTransaction.Type = models.TransactionType(autoDeterminedType)
	}

	if err := validateTransaction(newTransaction); err != nil {
		return fmt.Errorf("валидация транзакции не пройдена")
	}

	err = ts.storage.SaveTransactions(newTransaction)
	if err != nil {
		return fmt.Errorf("не удалось сохранить тразакцию")
	}

	return nil

}
