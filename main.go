package main

import (
	"log"
	"os"

	broadcast "github.com/muthuxv/esgi-go/channels"
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

	b := broadcast.NewBroadcaster(20)

	productRepository := product.NewProductRepository(db)
	productService := product.NewProductService(productRepository)
	productHAndler := handler.NewProductHandler(productService)

	paymentRepository := payment.NewPaymentRepository(db)
	paymentService := payment.NewPaymentService(paymentRepository)
	paymentHandler := handler.NewPaymentHandler(paymentService, b)

	r := gin.Default()
	api := r.Group("/api")

	api.POST("/createProduct", productHAndler.Create)
	api.GET("/products", productHAndler.FetchAll)
	api.GET("/product/:id", productHAndler.FetchById)
	api.PUT("/updateProduct/:id", productHAndler.Update)
	api.DELETE("/deleteProduct/:id", productHAndler.Delete)

	api.POST("/createPayment", paymentHandler.Create)
	api.GET("/payments", paymentHandler.FetchAll)
	api.GET("/payment/:id", paymentHandler.FetchById)
	api.PUT("/updatePayment/:id", paymentHandler.Update)
	api.DELETE("/deletePayment/:id", paymentHandler.Delete)

	api.GET("/stream/payment", paymentHandler.Stream)

	r.Run(":8080")

}
