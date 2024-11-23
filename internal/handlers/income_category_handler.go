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
	"shirinec.com/internal/services"
)

type IncomeCategoryHandler interface {
    List(c *gin.Context)
    GetByID(c *gin.Context)
}

type incomeCategoryHandler struct {
    incomeCategoryService services.IncomeCategoryService
}

func NewIncomeCategoryHandler(incomeCategoryService services.IncomeCategoryService) IncomeCategoryHandler {
    return &incomeCategoryHandler{incomeCategoryService: incomeCategoryService}
}

func (h *incomeCategoryHandler) List(c *gin.Context) {
    var input dto.ListIncomeCategoreisRequest
    if err := c.ShouldBindQuery(&input); err != nil {
        log.Printf("[Error] - incomeCategoryHandler.List - Bind query %+v\n", err)
        c.JSON(server_errors.InvalidInput.Unwrap())
        return
    }

    userID, err := uuid.Parse(c.GetString("user_id"))
    if err != nil {
        c.JSON(server_errors.InternalError.Unwrap())
        return
    }

    categories, err := h.incomeCategoryService.ListCategories(userID, input.Limit, input.Offset)
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
