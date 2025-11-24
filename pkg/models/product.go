package models

import "time"

type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	ImageURL    string    `json:"image_url"`
	Sellers     []Seller  `json:"sellers"`
	LastUpdated time.Time `json:"last_updated"`
}

type Seller struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Rating      float64 `json:"rating"`
	ReviewCount int     `json:"review_count"`
	IsOfficial  bool    `json:"is_official"`
}

type PriceAnalysis struct {
	ProductID         string          `json:"product_id"`
	MinPrice          float64         `json:"min_price"`
	MaxPrice          float64         `json:"max_price"`
	AvgPrice          float64         `json:"avg_price"`
	MedianPrice       float64         `json:"median_price"`
	OptimalPrice      float64         `json:"optimal_price"`
	DumpingSellers    []DumpingSeller `json:"dumping_sellers"`
	TotalSellers      int             `json:"total_sellers"`
	PriceDistribution map[string]int  `json:"price_distribution"`
	Timestamp         time.Time       `json:"timestamp"`
}

type DumpingSeller struct {
	SellerID       string  `json:"seller_id"`
	SellerName     string  `json:"seller_name"`
	Price          float64 `json:"price"`
	DumpingAmount  float64 `json:"dumping_amount"`
	DumpingPercent float64 `json:"dumping_percent"`
}

type PriceHistory struct {
	ProductID string              `json:"product_id"`
	History   []PriceHistoryEntry `json:"history"`
}

type PriceHistoryEntry struct {
	Timestamp    time.Time `json:"timestamp"`
	MinPrice     float64   `json:"min_price"`
	AvgPrice     float64   `json:"avg_price"`
	OptimalPrice float64   `json:"optimal_price"`
	SellerCount  int       `json:"seller_count"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}
