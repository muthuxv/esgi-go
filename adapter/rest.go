package adapter

import (
	"io"
	"net/http"

	"github.com/muthuxv/esgi-go/service"

	"github.com/gin-gonic/gin"
)

type GinAdapter interface {
	Stream(c *gin.Context)
	PostPayment(c *gin.Context)
	DeletePayment(c *gin.Context)
}

type messageInput struct {
	UserId string `json:"userId"`
	Data   string `json:"data"`
}

type response struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
	Data    string `json:"data,omitempty"`
}

type ginAdapter struct {
	paymentManager service.Manager
}

func NewGinAdapter(rm service.Manager) *ginAdapter {

	return &ginAdapter{rm}
}

// Stream godoc
// @Summary      Stream messages
// @Description  Stream messages from a payment
// @Tags         chat
// @Produce      text/event-stream
// @Param        id   path      string  true  "Payment ID"
// @Failure      400  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /stream/{id} [get]
func (ga *ginAdapter) Stream(c *gin.Context) {
	paymentid := c.Param("paymentid")
	listener := ga.paymentManager.OpenListener(paymentid)
	defer ga.paymentManager.CloseListener(paymentid, listener)

	clientGone := c.Request.Context().Done()
	c.Stream(func(w io.Writer) bool {
		select {
		case <-clientGone:
			return false
		case message := <-listener:
			serviceMsg, ok := message.(*service.Message)
			if !ok {
				c.SSEvent("message", message)
				return false
			}
			c.SSEvent("message", " "+serviceMsg.UserId+" → "+serviceMsg.Text)
			return true
		}
	})
}

// PostPayment godoc
// @Summary      Post to payment
// @Description  Post a message to a payment
// @Tags         chat
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Payment ID"
// @Param        messageInput   body      messageInput  true  "Message body"
// @Success      200  {object}  response
// @Failure      400  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /payment/{id} [post]
func (ga *ginAdapter) PostPayment(c *gin.Context) {
	paymentid := c.Param("paymentid")
	var input messageInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		NewError(c, http.StatusBadRequest, err)
		return
	}
	ga.paymentManager.Submit(input.UserId, paymentid, input.Data)

	c.JSON(http.StatusOK, &response{
		Success: true,
	})
}

func NewError(c *gin.Context, i int, err error) {
	panic("unimplemented")
}

// DeletePayment godoc
// @Summary      Delete a payment
// @Description  Delete the payment
// @Tags         chat
// @Produce      json
// @Param        id   path      string  true  "Payment ID"
// @Success      200  {object}  response
// @Failure      500  {object}  HTTPError
// @Router       /payment/{id} [delete]
func (ga *ginAdapter) DeletePayment(c *gin.Context) {
	paymentid := c.Param("paymentid")
	ga.paymentManager.DeleteBroadcast(paymentid)
	c.JSON(http.StatusOK, &response{
		Success: true,
	})
}
