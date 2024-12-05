package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"shirinec.com/internal/dto"
	"shirinec.com/internal/errors"
	"shirinec.com/internal/models"
	"shirinec.com/internal/services"
	"shirinec.com/internal/utils"
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
            utils.Logger.Infof("Error is ValidatorError: %s", err.Error())
			c.JSON(server_errors.ValidationErrorBuilder(errList).Unwrap())
			return
		}
		utils.Logger.Infof("Binding user input to dto.ItemCreateRequest: %s", err.Error())
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

	item, err := h.accountService.Create(context.Background(), &input, uid)
	if err != nil {
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
		utils.Logger.Errorf("Binding input query to dto.ListRequest: %s", err)
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		utils.Logger.Errorf("Parsing uuid from user_id string: %s", err)
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
		utils.Logger.Errorf("Parsing id param: %s", err)
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		utils.Logger.Errorf("Parsing uuid from user_id string: %s", err)
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
		utils.Logger.Errorf("Parsing id from context: %s", err)
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		utils.Logger.Errorf("Parsing uuid from user_id string: %s", err)
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
		utils.Logger.Errorf("Parsing id from param: %s", err)
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		utils.Logger.Errorf("Parsing uuid from user_id string: %s", err)
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
