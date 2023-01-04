package main

import (
	"log"
	"os"

	"github.com/muthuxv/esgi-go/handler"
	"github.com/muthuxv/esgi-go/payment"
	"github.com/muthuxv/esgi-go/product"

	"github.com/gin-gonic/gin"
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
	db.AutoMigrate(&product.Product{})
	db.AutoMigrate(&payment.Payment{})

	productRepository := product.NewProductRepository(db)
	productService := product.NewProductService(productRepository)
	productHAndler := handler.NewProductHandler(productService)

	paymentRepository := payment.NewPaymentRepository(db)
	paymentService := payment.NewPaymentService(paymentRepository)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	r := gin.Default()
	api := r.Group("/api")

	api.POST("/product", productHAndler.Create)
	api.GET("/products", productHAndler.FetchAll)
	api.GET("/task/:id", productHAndler.FetchById)
	api.PUT("/task/:id", productHAndler.Update)
	api.DELETE("/task/:id", productHAndler.Delete)

	api.POST("/payment", paymentHandler.Create)
	api.GET("/payments", paymentHandler.FetchAll)
	api.GET("/payment/:id", paymentHandler.FetchById)
	api.PUT("/payment/:id", paymentHandler.Update)
	api.DELETE("/payment/:id", paymentHandler.Delete)

	r.Run(":8080")

}
