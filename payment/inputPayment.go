package payment

type InputPayment struct {
	ProductID int     `json:"product_id"`
	PricePaid float64 `json:"price_paid"`
}
