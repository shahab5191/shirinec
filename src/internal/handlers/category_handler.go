package handler

import (
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

type CategoryHandler interface {
	Create(c *gin.Context)
	List(c *gin.Context)
	GetByID(c *gin.Context)
	Delete(c *gin.Context)
	Update(c *gin.Context)
}

type categoryHandler struct {
	categoryService services.CategoryService
}

func NewCategoryHandler(categoryService services.CategoryService) CategoryHandler {
	return &categoryHandler{
		categoryService: categoryService,
	}
}

func (h *categoryHandler) Create(c *gin.Context) {

	var input dto.CategoryCreateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		if errList := server_errors.AsValidatorError(err); errList != nil {
			c.JSON(server_errors.ValidationErrorBuilder(errList).Unwrap())
			return
		}
		utils.Logger.Errorf("[Error] - categoryHandler.Create - Bind body %s", err.Error())
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

	var category models.Category
	category.UserID = userID
	category.Name = &input.Name
	category.Color = &input.Color
	category.IconID = input.IconID
	category.EntityType = &input.Type

	err = h.categoryService.Create(&category)
	if err != nil {
        c.JSON(err.(*server_errors.SError).Unwrap())
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *categoryHandler) List(c *gin.Context) {
	var input dto.ListRequest
	if err := c.ShouldBindQuery(&input); err != nil {
		if errList := server_errors.AsValidatorError(err); errList != nil {
			c.JSON(server_errors.ValidationErrorBuilder(errList).Unwrap())
			return
		}
		utils.Logger.Errorf("Bind query %s", err.Error())
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
        utils.Logger.Errorf("Parse uuid from user_id string: %s", err.Error())
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	categories, err := h.categoryService.ListCategories(userID, input.Page, input.Size)

	if err != nil {
        c.JSON(err.(*server_errors.SError).Unwrap())
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (h *categoryHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Errorf("Parsing id param: %s", err.Error())
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	category, err := h.categoryService.GetByID(userID, int(id))
	if err != nil {
        c.JSON(err.(*server_errors.SError).Unwrap())
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *categoryHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Errorf("Parsing id param: %s", err.Error())
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
        utils.Logger.Errorf("Parsing user_id param: %s", err.Error())
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	err = h.categoryService.Delete(userID, int(id))
	if err != nil {
        c.JSON(err.(*server_errors.SError).Unwrap())
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *categoryHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Errorf("Parsing id param: %s", err.Error())
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	var input dto.CategoryUpdateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		if errList := server_errors.AsValidatorError(err); errList != nil {
			c.JSON(server_errors.ValidationErrorBuilder(errList).Unwrap())
			return
		}
		utils.Logger.Errorf("Bind body %s", err.Error())
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	if input.Color != nil && !utils.IsValidHexColor(*input.Color) {
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
        utils.Logger.Errorf("Getting user_id: %s", err.Error())
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	category, err := h.categoryService.Update(&userID, id, &input)
	if err != nil {
        c.JSON(err.(*server_errors.SError).Unwrap())
		return
	}

	c.JSON(http.StatusOK, category)
}
