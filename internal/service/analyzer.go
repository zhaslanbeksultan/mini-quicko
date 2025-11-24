package service

import (
	"fmt"
	"math"
	"mini-quicko/internal/storage"
	"mini-quicko/pkg/models"
	"sort"
	"time"
)

const (
	DumpingThreshold = 0.15 // 15% below average = dumping
	OptimalMargin    = 0.05 // 5% below average for optimal price
)

type Analyzer struct {
	store storage.Storage
	kaspi *KaspiService
}

func NewAnalyzer(store storage.Storage, kaspi *KaspiService) *Analyzer {
	return &Analyzer{
		store: store,
		kaspi: kaspi,
	}
}

func (a *Analyzer) AnalyzeProduct(productID string) (*models.PriceAnalysis, error) {
	product, err := a.kaspi.GetProduct(productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	if len(product.Sellers) == 0 {
		return nil, fmt.Errorf("no sellers found for product")
	}

	analysis := &models.PriceAnalysis{
		ProductID:    productID,
		TotalSellers: len(product.Sellers),
		Timestamp:    time.Now(),
	}

	prices := make([]float64, 0, len(product.Sellers))
	for _, seller := range product.Sellers {
		prices = append(prices, seller.Price)
	}

	sort.Float64s(prices)

	analysis.MinPrice = prices[0]
	analysis.MaxPrice = prices[len(prices)-1]
	analysis.AvgPrice = calculateAverage(prices)
	analysis.MedianPrice = calculateMedian(prices)

	dumpingThreshold := analysis.AvgPrice * (1 - DumpingThreshold)
	dumpingSellers := make([]models.DumpingSeller, 0)

	for _, seller := range product.Sellers {
		if seller.Price < dumpingThreshold {
			dumpingAmount := dumpingThreshold - seller.Price
			dumpingPercent := (dumpingAmount / analysis.AvgPrice) * 100

			dumpingSellers = append(dumpingSellers, models.DumpingSeller{
				SellerID:       seller.ID,
				SellerName:     seller.Name,
				Price:          seller.Price,
				DumpingAmount:  math.Round(dumpingAmount*100) / 100,
				DumpingPercent: math.Round(dumpingPercent*100) / 100,
			})
		}
	}

	sort.Slice(dumpingSellers, func(i, j int) bool {
		return dumpingSellers[i].DumpingPercent > dumpingSellers[j].DumpingPercent
	})

	analysis.DumpingSellers = dumpingSellers
	analysis.OptimalPrice = a.calculateOptimalPrice(prices, analysis.AvgPrice)
	analysis.PriceDistribution = a.calculatePriceDistribution(prices)

	a.saveToHistory(productID, analysis)

	return analysis, nil
}

func (a *Analyzer) calculateOptimalPrice(prices []float64, avgPrice float64) float64 {
	optimalPrice := avgPrice * (1 - OptimalMargin)

	competitiveCount := 0
	for _, price := range prices {
		if price <= optimalPrice {
			competitiveCount++
		}
	}

	if float64(competitiveCount)/float64(len(prices)) > 0.3 {
		optimalPrice = avgPrice * 0.98
	}

	return math.Round(optimalPrice*100) / 100
}

func (a *Analyzer) calculatePriceDistribution(prices []float64) map[string]int {
	dist := make(map[string]int)

	if len(prices) == 0 {
		return dist
	}

	minPrice := prices[0]
	maxPrice := prices[len(prices)-1]
	rangeSize := (maxPrice - minPrice) / 5

	dist["very_low"] = 0
	dist["low"] = 0
	dist["medium"] = 0
	dist["high"] = 0
	dist["very_high"] = 0

	for _, price := range prices {
		switch {
		case price < minPrice+rangeSize:
			dist["very_low"]++
		case price < minPrice+2*rangeSize:
			dist["low"]++
		case price < minPrice+3*rangeSize:
			dist["medium"]++
		case price < minPrice+4*rangeSize:
			dist["high"]++
		default:
			dist["very_high"]++
		}
	}

	return dist
}

func (a *Analyzer) GetHistory(productID string) (*models.PriceHistory, error) {
	return a.store.GetHistory(productID)
}

func (a *Analyzer) saveToHistory(productID string, analysis *models.PriceAnalysis) {
	entry := models.PriceHistoryEntry{
		Timestamp:    analysis.Timestamp,
		MinPrice:     analysis.MinPrice,
		AvgPrice:     analysis.AvgPrice,
		OptimalPrice: analysis.OptimalPrice,
		SellerCount:  analysis.TotalSellers,
	}

	a.store.AddHistoryEntry(productID, entry)
}

func calculateAverage(prices []float64) float64 {
	if len(prices) == 0 {
		return 0
	}
	sum := 0.0
	for _, p := range prices {
		sum += p
	}
	return math.Round((sum/float64(len(prices)))*100) / 100
}

func calculateMedian(prices []float64) float64 {
	n := len(prices)
	if n == 0 {
		return 0
	}
	if n%2 == 0 {
		return math.Round(((prices[n/2-1]+prices[n/2])/2)*100) / 100
	}
	return prices[n/2]
}
