package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"shirinec.com/internal/dto"
	"shirinec.com/internal/errors"
	"shirinec.com/internal/models"
	"shirinec.com/internal/services"
)

type AccountHandler interface {
	Create(c *gin.Context)
	List(c *gin.Context)
	GetByID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type accountHandler struct {
	accountService services.AccountService
}

func NewAccountHandler(accountService *services.AccountService) AccountHandler {
	return &accountHandler{
		accountService: *accountService,
	}
}

func (h *accountHandler) Create(c *gin.Context) {
	var input dto.AccountCreateRequest
	err := c.ShouldBindJSON(&input)
	if err != nil {
		if errList := server_errors.AsValidatorError(err); errList != nil {
			c.JSON(server_errors.ValidationErrorBuilder(errList).Unwrap())
			return
		}
		log.Printf("[Info] - accountHandler.Create - Binding user input to dto.ItemCreateRequest: %+v\n", err)
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	id := c.GetString("user_id")

	uid, err := uuid.Parse(id)
	if err != nil {
		log.Printf("[Error] - accountHandler.Create - Parsing uuid from id string: %+v\n", err)
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	item, err := h.accountService.Create(context.Background(), &input, uid)
	if err != nil {
		log.Printf("[Error] - accountHandler.Create - Calling accountService.Create :%+v\n", err)
		c.JSON(err.(*server_errors.SError).Unwrap())
		return
	}

	result := dto.CreateResponse[models.Account]{
		Result: *item,
	}

	c.JSON(http.StatusOK, result)
}

func (h *accountHandler) List(c *gin.Context) {
	var input dto.ListRequest
	if err := c.ShouldBindQuery(&input); err != nil {
		if errList := server_errors.AsValidatorError(err); errList != nil {
			c.JSON(server_errors.ValidationErrorBuilder(errList).Unwrap())
			return
		}
		log.Printf("[Error] - accountHandler.List - Binding input query to dto.ListRequest: %+v\n", err)
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		log.Printf("[Error] - accountHandler.List - Parsing uuid from user_id string: %+v\n", err)
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	items, err := h.accountService.List(context.Background(), input.Page, input.Size, userID)
	if err != nil {
		c.JSON(err.(*server_errors.SError).Unwrap())
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *accountHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("[Error] - accountHandler.GetByID - Parsing id param: %+v\n", err)
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		log.Printf("[Error] - accountHandler.GetByID - Parsing uuid from user_id string: %+v\n", err)
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	item, err := h.accountService.GetByID(context.Background(), id, userID)
	if err != nil {
		c.JSON(err.(*server_errors.SError).Unwrap())
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *accountHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("[Error] - accountHandler.Delete - Parsing id from context: %+v\n", err)
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		log.Printf("[Error] - accountHandler.Delete - Parsing uuid from user_id string: %+v\n", err)
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	err = h.accountService.Delete(context.Background(), id, userID)
	if err != nil {
		c.JSON(err.(*server_errors.SError).Unwrap())
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "Item deleted successfully!"})
}

func (h *accountHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("[Error] - accountHandler.Update - Parsing id from param: %+v\n", err)
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		log.Printf("[Error] - accountHandler.Update - Parsing uuid from user_id string: %+v\n", err)
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	var input dto.AccountUpdateRequest
	if err = c.ShouldBindJSON(&input); err != nil {
		if errList := server_errors.AsValidatorError(err); errList != nil {
			c.JSON(server_errors.ValidationErrorBuilder(errList).Unwrap())
			return
		}
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	item, err := h.accountService.Update(context.Background(), &input, id, userID)
	if err != nil {
		c.JSON(err.(*server_errors.SError).Unwrap())
		return
	}

	c.JSON(http.StatusOK, item)
}
