package main

import (
	"esgi-go/payment"
	"esgi-go/product"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "user:password@tcp(127.0.0.1:3306)/go-esgi?charset=utf8mb4&parseTime=True&loc=Local"
	}

	db, err := gorm.Open(mysql.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	// Migrate the schema
	db.AutoMigrate(&product & Product{})
	db.AutoMigrate(&payment & Payment{})

	productRepository := product.NewProductRepository(db)
	productService := product.NewProductService(productRepository)

	paymentRepository := payment.NewPaymentRepository(db)
	paymentService := payment.NewPaymentService(paymentRepository)

	r := gin.Default()
	api := r.Group("/api")
	api.GET("ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

}
