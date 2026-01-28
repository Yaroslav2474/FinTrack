package services

import (
	"fintrack/internal/models"
	"fintrack/internal/storage"
	"fmt"
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
	// Использование timestamp в нанасекундах для уникальности
	return fmt.Sprintf("tx_%d", time.Now().UnixNano())
}

func validateTransaction(transaction models.Transaction) error {
	if transaction.Amount <= 0 {
		return fmt.Errorf("сумма не может быть <= 0")
	}

	if transaction.Category == "" {
		return fmt.Errorf("категория не может быть пустой")
	}

	if transaction.Description == "" {
		return fmt.Errorf("описание не может быть пустым")
	}

	return nil
}

func (ts *TransactionService) AddTransaction(amount float64, category string, description string, transactionType string) error {
	if amount <= 0 {
		return fmt.Errorf("сумма не может быть <= 0")
	}

	if strings.TrimSpace(category) == "" {
		return fmt.Errorf("категория не может быть пустой")
	}

	if strings.TrimSpace(description) == "" {
		return fmt.Errorf("описание не может быть пустым")
	}

	if transactionType != string(models.Expense) && transactionType != string(models.Income) {
		return fmt.Errorf("неизвестный тип транзакции: %s", transactionType)
	}

	categories, err := ts.storage.GetCategories()
	if err != nil {
		return fmt.Errorf("ошибка получения категорий: %w", err)
	}

	categoryFound := false
	for _, cat := range categories {
		if strings.EqualFold(cat.Name, category) {
			categoryFound = true
			if string(cat.Type) != transactionType {
				return fmt.Errorf("несоответствие типа категории")
			}
			break
		}
	}

	if !categoryFound {
		return fmt.Errorf("категория не найдена")
	}

	newTransaction := models.Transaction{
		ID:          generateUniqueID(),
		Amount:      amount,
		Category:    category,
		Description: description,
		Type:        models.TransactionType(transactionType),
		Date:        time.Now(),
	}

	if err := validateTransaction(newTransaction); err != nil {
		return err
	}

	return ts.storage.SaveTransaction(newTransaction)
}

func (ts *TransactionService) GetAllTransactions() ([]models.Transaction, error) {
	transactions, err := ts.storage.GetAllTransactions()
	if err != nil {
		return nil, fmt.Errorf("не удалось получить транзакции: %w", err)
	}
	return transactions, nil
}
