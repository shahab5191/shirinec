package handler

import (
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
	"shirinec.com/internal/utils"
)

type ExpenseCategoryHandler interface {
	Create(c *gin.Context)
	List(c *gin.Context)
	GetByID(c *gin.Context)
	Delete(c *gin.Context)
	Update(c *gin.Context)
}

type expenseCategoryHandler struct {
	expenseCategoryService services.ExpenseCategoryService
}

func NewExpenseCategoryHandler(expenseCategoryService services.ExpenseCategoryService) ExpenseCategoryHandler {
	return &expenseCategoryHandler{expenseCategoryService: expenseCategoryService}
}

func (h *expenseCategoryHandler) Create(c *gin.Context) {

	var input dto.CreateIncomeCategoryRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("[Error] - expenseCategoryHandler.Create - Bind body %+v\n", err)
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	if !utils.IsValidHexColor(input.Color) {
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	var category models.ExpenseCategory
	category.UserID = userID
	category.Name = &input.Name
	category.Color = &input.Color

	err = h.expenseCategoryService.Create(&category)
	if err != nil {
		log.Printf("[Error] - expenseCategoryHandler.Create - Calling repo create %+v\n", err)
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *expenseCategoryHandler) List(c *gin.Context) {
	var input dto.ListRequest
	if err := c.ShouldBindQuery(&input); err != nil {
		log.Printf("[Error] - expenseCategoryHandler.List - Bind query %+v\n", err)
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	categories, err := h.expenseCategoryService.ListCategories(userID, input.Page, input.Size)

	if err != nil {
		var sErr *server_errors.SError
		if errors.As(err, &sErr) {
			c.JSON(sErr.Unwrap())
		} else {
			log.Printf("[Error] - expenseCategoryHandler.List - impossible error: %+v\n", err)
			c.JSON(server_errors.InternalError.Unwrap())
		}
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (h *expenseCategoryHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("[Error] - expenseCategoryHandler.GetByID - Parsing id param: %+v\n", err)
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	category, err := h.expenseCategoryService.GetByID(userID, int(id))
	if err != nil {
		log.Printf("[Error] - expenseCategoryHandler.GetByID - Getting category from service: %+v\n", err)
		var sErr *server_errors.SError
		if errors.As(err, &sErr) {
			c.JSON(sErr.Unwrap())
		} else {
			log.Printf("[Error] - expenseCategoryHandler.GetByID - Impossible Error!: %+v\n", err)
			c.JSON(server_errors.InternalError.Unwrap())
		}
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *expenseCategoryHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("[Error] - expenseCategoryHandler.Delete - parsing id param: %+v\n", err)
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	err = h.expenseCategoryService.Delete(userID, int(id))
	if err != nil {
		log.Printf("[Error] - expenseCategoryHandler.Delete - getting category from service: %+v\n", err)
		var sErr *server_errors.SError
		if errors.As(err, &sErr) {
			c.JSON(sErr.Unwrap())
		} else {
			log.Printf("[Error] - expenseCategoryHandler.Delete - Impossible Error!: %+v\n", err)
			c.JSON(server_errors.InternalError.Unwrap())
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *expenseCategoryHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("[Error] - expenseCategoryHandler.Update - parsing id param: %+v\n", err)
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	var input dto.UpdateIncomeCategoryRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("[Error] - expenseCategoryHandler.Create - Bind body %+v\n", err)
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	if input.Color != nil && !utils.IsValidHexColor(*input.Color) {
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	category, err := h.expenseCategoryService.Update(&userID, id, &input)
	if err != nil {
		log.Printf("[Error] - expenseCategoryHandler.Create - Calling repo create %+v\n", err)
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	c.JSON(http.StatusOK, category)
}
