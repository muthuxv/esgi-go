package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/muthuxv/esgi-go/adapter"
	"github.com/muthuxv/esgi-go/broadcast"
	"github.com/muthuxv/esgi-go/payment"
	"github.com/muthuxv/esgi-go/product"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	router := gin.Default()
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

	paymentRepository := payment.NewPaymentRepository(db)
	paymentService := payment.NewPaymentService(paymentRepository)

	productRepository := product.NewProductRepository(db)
	productService := product.NewProductService(productRepository)

	// get the broadcaster
	b := broadcast.NewBroadcaster(20)

	ginAdapter := adapter.NewGinAdapter(b, productService, paymentService)

	router.GET("/stream", ginAdapter.Stream)

	router.POST("/createPayment", ginAdapter.CreatePayment)
	router.GET("/payment/:id", ginAdapter.GetPayment)
	router.PUT("/updatePayment/:id", ginAdapter.UpdatePayment)
	router.DELETE("/deletePayment/:id", ginAdapter.DeletePayment)
	router.GET("/payments", ginAdapter.GetPayments)

	router.POST("/createProduct", ginAdapter.CreateProduct)
	router.PUT("/updateProduct/:id", ginAdapter.UpdateProduct)
	router.DELETE("/deleteProduct/:id", ginAdapter.DeleteProduct)
	router.GET("/product/:id", ginAdapter.GetProduct)
	router.GET("/products", ginAdapter.GetProducts)

	router.Run(fmt.Sprintf(":%v", 8084))

	router.Run() // listen and serve on

}
