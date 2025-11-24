package service

import (
	"encoding/json"
	"fmt"
	"mini-quicko/pkg/models"
	"os"
	"time"
)

type KaspiService struct {
	mockDataPath string
	cache        map[string]*models.Product
}

func NewKaspiService(mockDataPath string) *KaspiService {
	return &KaspiService{
		mockDataPath: mockDataPath,
		cache:        make(map[string]*models.Product),
	}
}

func (k *KaspiService) GetProduct(productID string) (*models.Product, error) {
	if product, exists := k.cache[productID]; exists {
		return product, nil
	}

	data, err := os.ReadFile(k.mockDataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read mock data: %w", err)
	}

	var products map[string]*models.Product
	if err := json.Unmarshal(data, &products); err != nil {
		return nil, fmt.Errorf("failed to parse mock data: %w", err)
	}

	product, exists := products[productID]
	if !exists {
		return nil, fmt.Errorf("product not found: %s", productID)
	}

	product.LastUpdated = time.Now()
	k.cache[productID] = product

	return product, nil
}

func (k *KaspiService) ClearCache() {
	k.cache = make(map[string]*models.Product)
}
