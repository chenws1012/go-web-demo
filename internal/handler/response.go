package handler

import "github.com/gin-gonic/gin"

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ListResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Meta    *MetaInfo   `json:"meta,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

type MetaInfo struct {
	Page      int `json:"page"`
	PageSize  int `json:"page_size"`
	Total     int `json:"total,omitempty"`
	TotalPage int `json:"total_page,omitempty"`
}

func SuccessResponse(c *gin.Context, code int, data interface{}) {
	c.JSON(code, Response{
		Success: true,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, code int, errCode, message string) {
	c.JSON(code, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    errCode,
			Message: message,
		},
	})
}

func ListSuccessResponse(c *gin.Context, code int, data interface{}, meta *MetaInfo) {
	c.JSON(code, ListResponse{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}
