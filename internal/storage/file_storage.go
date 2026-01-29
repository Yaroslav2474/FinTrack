package storage

import (
	"encoding/json"
	"fintrack/internal/models"
	"os"
	"path/filepath"
	"sync"
)

type FileStorage struct {
	transactionFile string
	categoryFile    string
	mu              sync.RWMutex
}

func NewFileStorage(transactionFile, categoryFile string) *FileStorage {
	// ensure directory for transactionFile exists
	if dir := filepath.Dir(transactionFile); dir != "" {
		_ = os.MkdirAll(dir, 0755)
	}
	if dir := filepath.Dir(categoryFile); dir != "" {
		_ = os.MkdirAll(dir, 0755)
	}

	// create transaction file if missing
	if _, err := os.Stat(transactionFile); err != nil {
		// write empty array
		_ = os.WriteFile(transactionFile, []byte("[]\n"), 0644)
	}

	// create category file if missing — write default categories
	if _, err := os.Stat(categoryFile); err != nil {
		var all []models.Category
		all = append(all, models.DefaultExpenseCategories...)
		all = append(all, models.DefaultIncomeCategories...)
		if data, err := json.MarshalIndent(all, "", "  "); err == nil {
			_ = os.WriteFile(categoryFile, append(data, '\n'), 0644)
		} else {
			_ = os.WriteFile(categoryFile, []byte("[]\n"), 0644)
		}
	}

	return &FileStorage{
		transactionFile: transactionFile,
		categoryFile:    categoryFile,
	}
}

func (fs *FileStorage) SaveTransaction(transaction models.Transaction) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	transactions, _ := fs.readTransactions()
	transactions = append(transactions, transaction)
	return fs.writeTransactions(transactions)
}

func (fs *FileStorage) GetAllTransactions() ([]models.Transaction, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	return fs.readTransactions()
}

func (fs *FileStorage) GetCategories() ([]models.Category, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	return fs.readCategories()
}

func (fs *FileStorage) SaveCategory(category models.Category) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	categories, _ := fs.readCategories()
	categories = append(categories, category)
	return fs.writeCategories(categories)
}

func (fs *FileStorage) readTransactions() ([]models.Transaction, error) {
	data, err := os.ReadFile(fs.transactionFile)
	if err != nil {
		return []models.Transaction{}, nil
	}

	var transactions []models.Transaction
	if err := json.Unmarshal(data, &transactions); err != nil {
		return []models.Transaction{}, nil
	}
	return transactions, nil
}

func (fs *FileStorage) writeTransactions(transactions []models.Transaction) error {
	data, err := json.MarshalIndent(transactions, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(fs.transactionFile, data, 0644)
}

func (fs *FileStorage) readCategories() ([]models.Category, error) {
	data, err := os.ReadFile(fs.categoryFile)
	if err != nil {
		return []models.Category{}, nil
	}

	var categories []models.Category
	if err := json.Unmarshal(data, &categories); err != nil {
		return []models.Category{}, nil
	}
	return categories, nil
}

func (fs *FileStorage) writeCategories(categories []models.Category) error {
	data, err := json.MarshalIndent(categories, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(fs.categoryFile, data, 0644)
}
