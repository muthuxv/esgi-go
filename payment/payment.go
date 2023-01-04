package payment

import (
	"github.com/muthuxv/esgi-go/product"
)

type Payment struct {
	ID        int `json:"id"`
	ProductID int
	Product   product.Product
	PricePaid float64 `json:"price_paid"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}
