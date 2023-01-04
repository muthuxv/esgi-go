package handler

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	broadcast "github.com/muthuxv/esgi-go/channels"
	"github.com/muthuxv/esgi-go/payment"
)

type paymentHandler struct {
	paymentService payment.Service
	broadcast      broadcast.Broadcaster
}

type Message struct {
	Text string
}

func NewPaymentHandler(paymentService payment.Service, broadcast broadcast.Broadcaster) *paymentHandler {
	return &paymentHandler{paymentService, broadcast}
}

func (ph *paymentHandler) Create(c *gin.Context) {
	// Get json body
	var input payment.InputPayment
	err := c.ShouldBindJSON(&input)
	if err != nil {
		response := &Response{
			Success: false,
			Message: "Cannot extract JSON body",
			Data:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	newPayment, err := ph.paymentService.Create(input)
	if err != nil {
		response := &Response{
			Success: false,
			Message: "Something went wrong",
			Data:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	ph.broadcast.Submit(Message{Text: "New payment created"})
	response := &Response{
		Success: true,
		Message: "New payment created",
		Data:    newPayment,
	}
	c.JSON(http.StatusCreated, response)
}

func (ph *paymentHandler) FetchAll(c *gin.Context) {
	payments, err := ph.paymentService.FetchAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, &Response{
			Success: false,
			Message: "Something went wrong",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &Response{
		Success: true,
		Data:    payments,
	})
}

func (ph *paymentHandler) FetchById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Message: "Wrong id parameter",
			Data:    err.Error(),
		})
		return
	}

	payment, err := ph.paymentService.FetchByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Message: "Something went wrong",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &Response{
		Success: true,
		Data:    payment,
	})
}

func (ph *paymentHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Message: "Wrong id parameter",
			Data:    err.Error(),
		})
		return
	}

	// Get json body
	var input payment.InputPayment
	err = c.ShouldBindJSON(&input)
	if err != nil {
		response := &Response{
			Success: false,
			Message: "Cannot extract JSON body",
			Data:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	uPayment, err := ph.paymentService.Update(id, input)
	if err != nil {
		response := &Response{
			Success: false,
			Message: "Something went wrong",
			Data:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	ph.broadcast.Submit(Message{Text: "New payment updated"})

	response := &Response{
		Success: true,
		Message: "New payment created",
		Data:    uPayment,
	}
	c.JSON(http.StatusCreated, response)
}

func (ph *paymentHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Message: "Wrong id parameter",
			Data:    err.Error(),
		})
		return
	}

	err = ph.paymentService.Delete(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Message: "Something went wrong",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &Response{
		Success: true,
		Message: "Payment successfully deleted",
	})
}

func (ph *paymentHandler) Stream(c *gin.Context) {
	listener := make(chan interface{})
	ph.broadcast.Register(listener)
	defer ph.broadcast.Unregister(listener)

	clientGone := c.Request.Context().Done()
	c.Stream(func(w io.Writer) bool {
		select {
		case <-clientGone:
			return false
		case message := <-listener:
			serviceMsg, ok := message.(Message)
			if !ok {
				c.SSEvent("message", message)
				return false
			}
			c.SSEvent("message", serviceMsg.Text)
			return true
		}
	})
}
