package model

// Customer represents a customer domain mapping.
type Customer struct {
	Domain     string `json:"domain"`
	CustomerID string `json:"customer_id"`
}
