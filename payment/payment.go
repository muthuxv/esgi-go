package payment

type Payment struct {
	ID        int    `json:"id"`
	ProductID int    `json:"product_id"`
	PricePaid int    `json:"price_paid"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
