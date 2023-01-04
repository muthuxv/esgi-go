package adapter

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/muthuxv/esgi-go/broadcast"
	"github.com/muthuxv/esgi-go/payment"
	"github.com/muthuxv/esgi-go/product"
)

type GinAdapter interface {
	Stream(c *gin.Context)

	CreatePayment(c *gin.Context)
	GetPayment(c *gin.Context)
	UpdatePayment(c *gin.Context)
	DeletePayment(c *gin.Context)
	GetPayments(c *gin.Context)

	UpdateProduct(c *gin.Context)
	CreateProduct(c *gin.Context)
	DeleteProduct(c *gin.Context)
	GetProduct(c *gin.Context)
	GetProducts(c *gin.Context)
}

type ginAdapter struct {
	broadcaster    broadcast.Broadcaster
	productService product.Service
	paymentService payment.Service
}

type Message struct {
	UserId string
	Text   string
}

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewGinAdapter(broadcaster broadcast.Broadcaster, productService product.Service, paymentService payment.Service) *ginAdapter {
	return &ginAdapter{
		broadcaster:    broadcaster,
		paymentService: paymentService,
		productService: productService,
	}
}

// Stream is the handler for the stream endpoint
func (adapter *ginAdapter) Stream(c *gin.Context) {

	//create a new channel to handle the stream
	listener := make(chan interface{})

	// get the broadcaster

	adapter.broadcaster.Register(listener)

	//close the channel when error message or client is gone
	defer adapter.broadcaster.Unregister(listener)

	clientGone := c.Request.Context().Done()

	c.Stream(func(w io.Writer) bool {
		select {
		case <-clientGone:
			return false
		case message := <-listener:
			serviceMsg, ok := message.(Message)
			if !ok {
				fmt.Println("not a message")
				c.SSEvent("message", message)
				return false
			}
			c.SSEvent("message", " "+serviceMsg.UserId+" → "+serviceMsg.Text)
			return true
		}
	})

	fmt.Println("stream is OK")
}

func (adapter *ginAdapter) CreatePayment(c *gin.Context) {

	productId, err := strconv.Atoi(c.PostForm("productId"))

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "invalid product id",
		})
		return
	}

	// get the product

	product, err := adapter.productService.FetchByID(productId)

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "something went wrong",
		})
		return
	}

	payment, err := adapter.paymentService.Create(product)

	b := adapter.broadcaster

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "something went wrong",
		})
		return
	}

	b.Submit(Message{
		UserId: "1",
		Text:   "payment is created",
	})

	c.JSON(http.StatusOK, &Response{
		Status:  http.StatusOK,
		Message: "payment is created",
		Data:    payment,
	})

}

func (adapter *ginAdapter) GetPayment(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "invalid id",
		})
		return
	}

	payment, err := adapter.paymentService.FetchByID(id)

	b := adapter.broadcaster

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "something went wrong",
		})
		return
	}

	b.Submit(Message{
		UserId: "1",
		Text:   "payment price is " + strconv.FormatFloat(payment.PricePaid, 'f', 2, 6),
	})

	c.JSON(http.StatusOK, &Response{
		Status:  http.StatusOK,
		Message: "payment price is " + strconv.FormatFloat(payment.PricePaid, 'f', 2, 6),
		Data:    payment,
	})

}

func (adapter *ginAdapter) UpdatePayment(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "invalid id",
		})
		return
	}

	var payment payment.Payment

	err = c.ShouldBindJSON(&payment)

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "invalid payment",
		})
		return
	}

	payment, err = adapter.paymentService.Update(id, payment)

	b := adapter.broadcaster

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "something went wrong",
		})
		return
	}

	b.Submit(Message{
		UserId: "1",
		Text:   "payment is updated",
	})

	c.JSON(http.StatusOK, &Response{
		Status:  http.StatusOK,
		Message: "payment is updated",
		Data:    payment,
	})

}

func (adapter *ginAdapter) DeletePayment(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "invalid id",
		})
		return
	}

	err = adapter.paymentService.Delete(id)

	b := adapter.broadcaster

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "something went wrong",
		})
		return
	}

	b.Submit(Message{
		UserId: "1",
		Text:   "payment is deleted",
	})

	c.JSON(http.StatusOK, &Response{
		Status:  http.StatusOK,
		Message: "payment is deleted",
	})

}

func (adapter *ginAdapter) GetPayments(c *gin.Context) {

	payments, err := adapter.paymentService.FetchAll()

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "something went wrong",
		})
		return
	}

	b := adapter.broadcaster

	b.Submit(Message{
		UserId: "1",
		Text:   "Payments are fetched",
	})

	c.JSON(http.StatusOK, &Response{
		Status:  http.StatusOK,
		Message: "Payments are fetched",
		Data:    payments,
	})

}

func (adapter *ginAdapter) GetProducts(c *gin.Context) {

	products, err := adapter.productService.FetchAll()

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "something went wrong",
		})
		return
	}

	b := adapter.broadcaster

	b.Submit(Message{
		UserId: "1",
		Text:   "Products are fetched",
	})

	c.JSON(http.StatusOK, &Response{
		Status:  http.StatusOK,
		Message: "Products are fetched",
		Data:    products,
	})

}

func (adapter *ginAdapter) CreateProduct(c *gin.Context) {

	name := c.PostForm("name")

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		fmt.Println(name)
		return
	}

	price := c.PostForm("price")
	priceFloat, err := strconv.ParseFloat(price, 64)

	if price == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid price"})
		return
	}

	fmt.Println("create product", c.PostForm("name"))
	// get the broadcaster
	b := adapter.broadcaster

	// save the payment
	product, err := adapter.productService.Create(name, priceFloat)

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "invalid product",
			Data:    err.Error(),
		})
		return
	}

	b.Submit(Message{
		UserId: "1",
		Text:   product.Name + " is created",
	})

	response := &Response{
		Status:  http.StatusOK,
		Message: "Product is created",
		Data:    product,
	}

	c.JSON(http.StatusOK, response)

}

func (adapter *ginAdapter) UpdateProduct(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "invalid id",
			Data:    err.Error(),
		})
		return
	}

	var product product.Product

	err = c.ShouldBindJSON(&product)

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "invalid product",
			Data:    err.Error(),
		})
		return
	}

	updatedProduct, err := adapter.productService.Update(id, product)

	b := adapter.broadcaster

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	b.Submit(Message{
		UserId: "1",
		Text:   updatedProduct.Name + " is updated to price " + strconv.FormatFloat(updatedProduct.Price, 'f', 2, 64),
	})

	c.JSON(http.StatusOK, &Response{
		Status:  http.StatusOK,
		Message: "product updated",
		Data:    updatedProduct,
	})

}

func (adapter *ginAdapter) DeleteProduct(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "invalid id",
			Data:    err.Error(),
		})
		return
	}

	err = adapter.productService.Delete(id)

	b := adapter.broadcaster

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "product not found",
			Data:    err.Error(),
		})
		return
	}

	b.Submit(Message{
		UserId: "1",
		Text:   "Product is deleted",
	})

	c.JSON(http.StatusOK, &Response{
		Status:  http.StatusOK,
		Message: "product deleted",
		Data:    "product deleted",
	})

}

func (adapter *ginAdapter) GetProduct(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "invalid id",
			Data:    err.Error(),
		})
		return
	}

	product, err := adapter.productService.FetchByID(id)

	b := adapter.broadcaster

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "product not found",
			Data:    err.Error(),
		})
		return
	}

	b.Submit(Message{
		UserId: "1",
		Text:   product.Name + " is found",
	})

	c.JSON(http.StatusOK, &Response{
		Status:  http.StatusOK,
		Message: "product found",
		Data:    product,
	})

}
