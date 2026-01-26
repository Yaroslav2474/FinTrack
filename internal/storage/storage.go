package storage

import "fintrack/internal/models"

type Storage interface {
	SaveTransaction(transaction models.Transaction) error
	GetAllTransactions() ([]models.Transaction, error)
	GetCategories() ([]models.Category, error)
	SaveCategory(category models.Category) error
}
