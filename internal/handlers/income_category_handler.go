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

type IncomeCategoryHandler interface {
    Create(c *gin.Context)
    List(c *gin.Context)
    GetByID(c *gin.Context)
    Delete(c *gin.Context)
    Update(c *gin.Context)
}

type incomeCategoryHandler struct {
    incomeCategoryService services.IncomeCategoryService
}

func NewIncomeCategoryHandler(incomeCategoryService services.IncomeCategoryService) IncomeCategoryHandler {
    return &incomeCategoryHandler{incomeCategoryService: incomeCategoryService}
}

func (h *incomeCategoryHandler) Create(c *gin.Context) {

    var input dto.CreateIncomeCategoryRequest
    if err := c.ShouldBindJSON(&input); err != nil {
        log.Printf("[Error] - incomeCategoryHandler.Create - Bind body %+v\n", err)
        c.JSON(server_errors.InvalidInput.Unwrap())
        return
    }

    if !utils.IsValidHexColor(input.Color){
        c.JSON(server_errors.InvalidInput.Unwrap())
        return
    }

    userID, err := uuid.Parse(c.GetString("user_id"))
    if err != nil {
        c.JSON(server_errors.InternalError.Unwrap())
        return
    }

    var category models.IncomeCategory
    category.UserID = userID
    category.Name = &input.Name
    category.Color = &input.Color

    err = h.incomeCategoryService.Create(&category)
    if err != nil {
        log.Printf("[Error] - incomeCategoryHandler.Create - Calling repo create %+v\n", err)
        c.JSON(server_errors.InternalError.Unwrap())
        return
    }

    c.JSON(http.StatusOK, category)
}

func (h *incomeCategoryHandler) List(c *gin.Context) {
    var input dto.ListRequest
    if err := c.ShouldBindQuery(&input); err != nil {
        log.Printf("[Error] - incomeCategoryHandler.List - Bind query %+v\n", err)
        c.JSON(server_errors.InvalidInput.Unwrap())
        return
    }
    log.Printf("Input: %+v\n", input)

    userID, err := uuid.Parse(c.GetString("user_id"))
    if err != nil {
        c.JSON(server_errors.InternalError.Unwrap())
        return
    }

    categories, err := h.incomeCategoryService.ListCategories(userID, input.Page, input.Size)
    if err != nil {
        var sErr *server_errors.SError
        if errors.As(err, &sErr){
            c.JSON(sErr.Unwrap())
        }else{
            log.Printf("[Error] - icomeCategoryHandler.List - impossible error: %+v\n", err)   
            c.JSON(server_errors.InternalError.Unwrap())
        }
        return
    }

    c.JSON(http.StatusOK, categories)
}

func (h * incomeCategoryHandler)GetByID (c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil{
        log.Printf("[Error] - IncomeCategoryHandler.GetByID - parsing id param: %+v\n", err)
        c.JSON(server_errors.InvalidInput.Unwrap())
        return
    }

    userID, err := uuid.Parse(c.GetString("user_id"))
    if err != nil {
        c.JSON(server_errors.InternalError.Unwrap())
        return
    }

    category, err := h.incomeCategoryService.GetByID(userID, int(id))
    if err != nil {
        log.Printf("[Error] - IncomeCategoryHandler.GetByID - getting category from service: %+v\n", err)
        var sErr *server_errors.SError
        if errors.As(err, &sErr) {
            c.JSON(sErr.Unwrap())
        }else{
            log.Printf("[Error] - incomeCategoryHandler.GetByID - Impossible Error!: %+v\n", err)
            c.JSON(server_errors.InternalError.Unwrap())
        }
        return
    }

    c.JSON(http.StatusOK, category)
}


func (h * incomeCategoryHandler)Delete (c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil{
        log.Printf("[Error] - IncomeCategoryHandler.Delete - parsing id param: %+v\n", err)
        c.JSON(server_errors.InvalidInput.Unwrap())
        return
    }

    userID, err := uuid.Parse(c.GetString("user_id"))
    if err != nil {
        c.JSON(server_errors.InternalError.Unwrap())
        return
    }

    err = h.incomeCategoryService.Delete(userID, int(id))
    if err != nil {
        log.Printf("[Error] - IncomeCategoryHandler.Delete - getting category from service: %+v\n", err)
        var sErr *server_errors.SError
        if errors.As(err, &sErr) {
            c.JSON(sErr.Unwrap())
        }else{
            log.Printf("[Error] - incomeCategoryHandler.Delete - Impossible Error!: %+v\n", err)
            c.JSON(server_errors.InternalError.Unwrap())
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *incomeCategoryHandler) Update(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil{
        log.Printf("[Error] - IncomeCategoryHandler.Update - parsing id param: %+v\n", err)
        c.JSON(server_errors.InvalidInput.Unwrap())
        return
    }

    var input dto.UpdateIncomeCategoryRequest
    if err := c.ShouldBindJSON(&input); err != nil {
        log.Printf("[Error] - incomeCategoryHandler.Create - Bind body %+v\n", err)
        c.JSON(server_errors.InvalidInput.Unwrap())
        return
    }

    if input.Color != nil && !utils.IsValidHexColor(*input.Color){
        c.JSON(server_errors.InvalidInput.Unwrap())
        return
    }

    userID, err := uuid.Parse(c.GetString("user_id"))
    if err != nil {
        c.JSON(server_errors.InternalError.Unwrap())
        return
    }


    category, err := h.incomeCategoryService.Update(&userID, id, &input)
    if err != nil {
        log.Printf("[Error] - incomeCategoryHandler.Create - Calling repo create %+v\n", err)
        c.JSON(server_errors.InternalError.Unwrap())
        return
    }

    c.JSON(http.StatusOK, category)
}
