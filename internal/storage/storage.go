package storage

import "fintrack/internal/models"

type Storage interface {
	SaveTransactions(models.Transaction) error
	LoadTransactions() ([]models.Transaction, error)
	GetCategories() ([]models.Category, error)
}
