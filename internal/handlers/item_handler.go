package handler

import (
	"context"
	"errors"
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

type ItemHandler interface {
	Create(c *gin.Context)
	List(c *gin.Context)
	GetByID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type itemHandler struct {
	itemService services.ItemService
}

func NewItemHandler(itemService *services.ItemService) ItemHandler {
	return &itemHandler{itemService: *itemService}
}

func (h *itemHandler) Create(c *gin.Context) {
	var input dto.ItemCreateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		if errList := server_errors.AsValidatorError(err); errList != nil {
			c.JSON(server_errors.ValidationErrorBuilder(errList).Unwrap())
			return
		}
		log.Printf("[Warning] - itemHandler.Create - Undefined error binding user input to dto.ItemCreateRequest: %+v\n", err)
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	id := c.GetString("user_id")

	uid, err := uuid.Parse(id)
	if err != nil {
		log.Printf("[Error] - itemHandler.Create - Parsing uuid from id string: %+v\n", err)
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	item, err := h.itemService.Create(context.Background(), &input, uid)
	if err != nil {
		log.Printf("[Error] - itemHandler.Create - Calling itemService.Create :%+v\n", err)
		c.JSON(err.(*server_errors.SError).Unwrap())
		return
	}

	result := dto.CreateResponse[models.Item]{
		Result: *item,
	}

	c.JSON(http.StatusOK, result)
}

func (h *itemHandler) List(c *gin.Context) {
	var input dto.ListRequest
	if err := c.ShouldBindQuery(&input); err != nil {
		if errList := server_errors.AsValidatorError(err); errList != nil {
			c.JSON(server_errors.ValidationErrorBuilder(errList).Unwrap())
			return
		}
		log.Printf("[Warning] - itemHandler.List - Undefined error while binding input query to dto.ListRequest: %+v\n", err)
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		log.Printf("[Error] - itemHandler.List - Parsing uuid from user_id string: %+v\n", err)
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	items, err := h.itemService.List(context.Background(), input.Page, input.Size, userID)
	if err != nil {
		if sErr, ok := err.(*server_errors.SError); ok {
			c.JSON(sErr.Unwrap())
		}
		log.Printf("[Error] - itemHandler.List - impossible error: %+v\n", err)
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *itemHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("[Error] - itemHandler.GetByID - Parsing id param: %+v\n", err)
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		log.Printf("[Error] - itemHandler.GetByID - Parsing uuid from user_id string: %+v\n", err)
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	item, err := h.itemService.GetByID(context.Background(), id, userID)
	if err != nil {
		var sErr *server_errors.SError
		if errors.As(err, &sErr) {
			c.JSON(sErr.Unwrap())
			return
		}
		log.Printf("[Error] - itemHandler.GetByID - impossible Error!: %+v\n", err)
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *itemHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("[Error] - itemHandler.Delete - Parsing id from context: %+v\n", err)
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		log.Printf("[Error] - itemHandler.Delete - Parsing uuid from user_id string: %+v\n", err)
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	if err = h.itemService.Delete(context.Background(), id, userID); err != nil {
		var sErr *server_errors.SError
		if errors.As(err, &sErr) {
			c.JSON(sErr.Unwrap())
			return
		}
		log.Printf("[Error] - itemService.Delete - Impossible error: %+v\n", err)
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "Item deleted successfully!"})
}

func (h *itemHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("[Error] - itemHandler.Update - Parsing id from param: %+v\n", err)
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		log.Printf("[Error] - itemHandler.Update - Parsing uuid from user_id string: %+v\n", err)
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	var input dto.ItemUpdateRequest
	if err = c.ShouldBindJSON(&input); err != nil {
		if errList := server_errors.AsValidatorError(err); errList != nil {
			c.JSON(server_errors.ValidationErrorBuilder(errList).Unwrap())
			return
		}
		log.Printf("[Warning] - itemHandler.Update - Undefined error while binding input to dto.ItemUpdateRequest: %+v\n", err)
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	item, err := h.itemService.Update(context.Background(), &input, id, userID)
	if err != nil {
		var sErr *server_errors.SError
		if errors.As(err, &sErr) {
			c.JSON(err.(*server_errors.SError).Unwrap())
			return
		}
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	c.JSON(http.StatusOK, item)
}
