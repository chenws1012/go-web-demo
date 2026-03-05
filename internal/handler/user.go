package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go-web-demo/internal/repository"
)

type UserHandler struct {
	userService UserService
}

type UserService interface {
	CreateUser(ctx context.Context, username, email string) (*repository.User, error)
	GetUser(ctx context.Context, id int64) (*repository.User, error)
	ListUsers(ctx context.Context, page, pageSize int) ([]*repository.User, error)
	UpdateUser(ctx context.Context, id int64, username, email string) (*repository.User, error)
	DeleteUser(ctx context.Context, id int64) error
}

type UserData struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email" binding:"omitempty,email"`
}

func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), req.Username, req.Email)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "CREATE_ERROR", err.Error())
		return
	}

	SuccessResponse(c, http.StatusCreated, toUserData(user))
}

func (h *UserHandler) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid user ID")
		return
	}

	user, err := h.userService.GetUser(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "user not found" {
			ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "User not found")
			return
		}
		ErrorResponse(c, http.StatusInternalServerError, "GET_ERROR", err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, toUserData(user))
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	users, err := h.userService.ListUsers(c.Request.Context(), page, pageSize)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "LIST_ERROR", err.Error())
		return
	}

	response := make([]*UserData, len(users))
	for i, user := range users {
		response[i] = toUserData(user)
	}

	ListSuccessResponse(c, http.StatusOK, response, &MetaInfo{
		Page:     page,
		PageSize: pageSize,
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid user ID")
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), id, req.Username, req.Email)
	if err != nil {
		if err.Error() == "user not found" {
			ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "User not found")
			return
		}
		ErrorResponse(c, http.StatusInternalServerError, "UPDATE_ERROR", err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, toUserData(user))
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid user ID")
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), id); err != nil {
		if err.Error() == "user not found" {
			ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "User not found")
			return
		}
		ErrorResponse(c, http.StatusInternalServerError, "DELETE_ERROR", err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func toUserData(user *repository.User) *UserData {
	return &UserData{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}
