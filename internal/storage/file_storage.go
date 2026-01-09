package storage

import (
	"encoding/json"
	"fintrack/internal/models"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type FileStorage struct {
	dataDir     string
	mutex       sync.Mutex
	filePerms   os.FileMode
	dirPerms    os.FileMode
	autoCreater bool
	encoder     *json.Encoder
}

func NewFileStorage(dataDir string) (*FileStorage, error) {

	return &FileStorage{dataDir: dataDir}, nil
}

func (fs FileStorage) SaveTransactions(transactions []models.Transaction) error {

	jsonData, err := json.MarshalIndent(transactions, "", " ")

	if err != nil {
		return fmt.Errorf("не удалось преобразовать транзакции в JSON: %w", err)
	}

	filePath := filepath.Join(fs.dataDir, "transaction.json")

	if err = os.MkdirAll(fs.dataDir, 0755); err != nil {
		return fmt.Errorf("не удалось создать директорию %s: %w", fs.dataDir, err)
	}

	if err = os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("не удалось записать файл %s: %w", filePath, err)
	}

	return nil

}

func (fs FileStorage) LoadTransactions() ([]models.Transaction, error) {

	var transaction []models.Transaction

	filePath := filepath.Join(fs.dataDir, "transaction.json")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []models.Transaction{}, nil
	}

	filesize, err := os.ReadFile(filePath)

	if err != nil {
		return []models.Transaction{}, fmt.Errorf("не удалось прочитать файл по пути %s: %w", filePath, err)
	}

	err = json.Unmarshal(filesize, &transaction)

	if err != nil {
		return []models.Transaction{}, fmt.Errorf("не удалось запарсить файл по пути %s: %w", filePath, err)
	}

	return transaction, nil

}

func (fs FileStorage) ensureDataDir() error {

	filePath := filepath.Join(fs.dataDir, "transaction.json")

	_, err := os.Stat(filePath)

	if err != nil {
		err = os.MkdirAll(filePath, 0755)
	}

	return os.MkdirAll(fs.dataDir, fs.dirPerms)
}

func (fs FileStorage) getFilePath(filename string) string {

	return filepath.Join(fs.dataDir, filename)

}
