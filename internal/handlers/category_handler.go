package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"shirinec.com/internal/dto"
	"shirinec.com/internal/enums"
	"shirinec.com/internal/errors"
	"shirinec.com/internal/models"
	"shirinec.com/internal/services"
	"shirinec.com/internal/utils"
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
	validate        *validator.Validate
}

func NewCategoryHandler(categoryService services.CategoryService, validate *validator.Validate) CategoryHandler {
	return &categoryHandler{
		categoryService: categoryService,
		validate:        validate,
	}
}

func (h *categoryHandler) Create(c *gin.Context) {

	var input dto.CategoryCreateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		if errList := server_errors.AsValidatorError(err); errList != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errList})
			return
		}
		log.Printf("[Error] - categoryHandler.Create - Bind body %+v\n", err)
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
	entityTypeStr := string(input.Type)
	entityTypeStr = strings.ToLower(entityTypeStr)
	entityTypeStr = string(unicode.ToUpper(rune(entityTypeStr[0]))) + entityTypeStr[1:]
	entityType := enums.CategoryType(entityTypeStr)
	category.EntityType = &entityType
	log.Printf("Category Object: %+v\n", category)

	err = h.categoryService.Create(&category)
	if err != nil {
		log.Printf("[Error] - categoryHandler.Create - Calling repo create %+v\n", err)
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *categoryHandler) List(c *gin.Context) {
	var input dto.ListRequest
	if err := c.ShouldBindQuery(&input); err != nil {
		log.Printf("[Error] - categoryHandler.List - Bind query %+v\n", err)
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	categories, err := h.categoryService.ListCategories(userID, input.Page, input.Size)

	if err != nil {
		var sErr *server_errors.SError
		if errors.As(err, &sErr) {
			c.JSON(sErr.Unwrap())
		}
		log.Printf("[Error] - categoryHandler.List - impossible error: %+v\n", err)
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (h *categoryHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("[Error] - categoryHandler.GetByID - Parsing id param: %+v\n", err)
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
		log.Printf("[Error] - categoryHandler.GetByID - Getting category from service: %+v\n", err)
		var sErr *server_errors.SError
		if errors.As(err, &sErr) {
			c.JSON(sErr.Unwrap())
		} else {
			log.Printf("[Error] - categoryHandler.GetByID - Impossible Error!: %+v\n", err)
			c.JSON(server_errors.InternalError.Unwrap())
		}
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *categoryHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("[Error] - categoryHandler.Delete - parsing id param: %+v\n", err)
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	err = h.categoryService.Delete(userID, int(id))
	if err != nil {
		log.Printf("[Error] - categoryHandler.Delete - getting category from service: %+v\n", err)
		var sErr *server_errors.SError
		if errors.As(err, &sErr) {
			c.JSON(sErr.Unwrap())
		} else {
			log.Printf("[Error] - categoryHandler.Delete - Impossible Error!: %+v\n", err)
			c.JSON(server_errors.InternalError.Unwrap())
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *categoryHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("[Error] - categoryHandler.Update - parsing id param: %+v\n", err)
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	var input dto.CategoryUpdateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		if errList := server_errors.AsValidatorError(err); errList != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errList})
			return
		}
		log.Printf("[Error] - categoryHandler.Update - Bind body %+v\n", err)
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

	category, err := h.categoryService.Update(&userID, id, &input)
	if err != nil {
		log.Printf("[Error] - categoryHandler.Update - Calling category service update %+v\n", err)
		var seError *server_errors.SError
		if errors.As(err, &seError) {
			c.JSON(seError.Unwrap())
		} else {
			c.JSON(server_errors.InternalError.Unwrap())
		}
		return
	}

	c.JSON(http.StatusOK, category)
}
