package services

import (
	"fintrack/internal/models"
	"fintrack/internal/storage"
)

type CategoryService struct {
	storage storage.Storage
}

func NewCategoryService(storage storage.Storage) *CategoryService {
	return &CategoryService{
		storage: storage,
	}
}

// GetCategoriesByType returns categories filtered by income/expense type.
func (cs *CategoryService) GetCategoriesByType(isIncome bool) ([]models.Category, error) {
	allCategories, err := cs.storage.GetCategories()
	if err != nil {
		return nil, err
	}

	var filtered []models.Category
	for _, cat := range allCategories {
		if cat.IsIncome == isIncome {
			filtered = append(filtered, cat)
		}
	}
	return filtered, nil
}

func (cs *CategoryService) AddCategory(name string, isIncome bool) error {
	category := models.Category{
		Name:     name,
		IsIncome: isIncome,
		Type:     map[bool]string{true: "income", false: "expense"}[isIncome],
	}
	return cs.storage.SaveCategory(category)
}
