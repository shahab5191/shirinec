package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"shirinec.com/src/internal/dto"
	"shirinec.com/src/internal/errors"
	"shirinec.com/src/internal/models"
	"shirinec.com/src/internal/services"
	"shirinec.com/src/internal/utils"
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
		utils.Logger.Warnf("Undefined error binding user input to dto.ItemCreateRequest: %s", err.Error())
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	id := c.GetString("user_id")

	uid, err := uuid.Parse(id)
	if err != nil {
		utils.Logger.Errorf("Parsing uuid from id string: %s", err.Error())
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	item, err := h.itemService.Create(context.Background(), &input, uid)
	if err != nil {
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
		utils.Logger.Warnf("Undefined error while binding input query to dto.ListRequest: %s", err.Error())
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		utils.Logger.Errorf("Parsing uuid from user_id string: %s", err.Error())
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	items, err := h.itemService.List(context.Background(), input.Page, input.Size, userID)
	if err != nil {
        c.JSON(err.(*server_errors.SError).Unwrap())
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *itemHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Errorf("Parsing id param: %s", err.Error())
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		utils.Logger.Errorf("Parsing uuid from user_id string: %s", err.Error())
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	item, err := h.itemService.GetByID(context.Background(), id, userID)
	if err != nil {
		c.JSON(err.(*server_errors.SError).Unwrap())
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *itemHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Errorf("Parsing id from context: %s", err.Error())
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		utils.Logger.Errorf("Parsing uuid from user_id string: %s", err.Error())
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	if err = h.itemService.Delete(context.Background(), id, userID); err != nil {
		c.JSON(err.(*server_errors.SError).Unwrap())
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "Item deleted successfully!"})
}

func (h *itemHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Errorf("Parsing id from param: %s", err.Error())
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		utils.Logger.Errorf("Parsing uuid from user_id string: %s", err.Error())
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	var input dto.ItemUpdateRequest
	if err = c.ShouldBindJSON(&input); err != nil {
		if errList := server_errors.AsValidatorError(err); errList != nil {
			c.JSON(server_errors.ValidationErrorBuilder(errList).Unwrap())
			return
		}
		utils.Logger.Warnf("Undefined error while binding input to dto.ItemUpdateRequest: %s", err.Error())
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	item, err := h.itemService.Update(context.Background(), &input, id, userID)
	if err != nil {
		c.JSON(err.(*server_errors.SError).Unwrap())
		return
	}

	c.JSON(http.StatusOK, item)
}
