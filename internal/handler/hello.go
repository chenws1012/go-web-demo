package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HelloHandler struct {
	helloService HelloService
}

type HelloService interface {
	SayHello(ctx context.Context) (string, error)
}

func NewHelloHandler(helloService HelloService) *HelloHandler {
	return &HelloHandler{
		helloService: helloService,
	}
}

type HelloResponse struct {
	Message string `json:"message"`
}

func (h *HelloHandler) SayHello(c *gin.Context) {
	message, err := h.helloService.SayHello(c.Request.Context())
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "HELLO_ERROR", err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, HelloResponse{
		Message: message,
	})
}
