package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"shirinec.com/internal/dto"
	"shirinec.com/internal/errors"
	"shirinec.com/internal/services"
	"shirinec.com/internal/utils"
)

type FinancialGroupHandler interface {
	Create(c *gin.Context)
	AddUser(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
	Delete(c *gin.Context)
}

type financialGroupHandler struct {
	financialGroupService services.FinancialGroupService
}

func NewFinancialGroupHandler(financialGroupService *services.FinancialGroupService) FinancialGroupHandler {
	return &financialGroupHandler{
		financialGroupService: *financialGroupService,
	}
}

func (h *financialGroupHandler) Create(c *gin.Context) {
	var input dto.FinancialGroupCreateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		if errList := server_errors.AsValidatorError(err); errList != nil {
			utils.Logger.Infof("Error is ValidatorError: %s", err.Error())
			c.JSON(server_errors.ValidationErrorBuilder(errList).Unwrap())
			return
		}

		utils.Logger.Errorf("financialGroupHandler.Create - Binding request body to dto.FinancialGroupCreateRequest: %s", err.Error())
		c.JSON(server_errors.InvalidAuthorizationHeader.Unwrap())
		return
	}

	uid, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		utils.Logger.Errorf("Parsing uuid from user_id string: %s", err.Error())
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	item, err := h.financialGroupService.Create(context.Background(), &input, uid)
	if err != nil {
		c.JSON(err.(*server_errors.SError).Unwrap())
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *financialGroupHandler) AddUser(c *gin.Context) {
	financialGroupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Errorf("financialGroupHandler.AddUser - Parsing id: %s", err.Error())
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	memberID, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		utils.Logger.Errorf("financialGroupHandler.AddUser - Getting userID param: %s", err.Error())
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		utils.Logger.Errorf("financialGroupHandler.AddUser - Parsing uuid from user_id: %s", err.Error())
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	if err = h.financialGroupService.AddUserToGroup(context.Background(), financialGroupID, memberID, userID); err != nil {
		c.JSON(err.(*server_errors.SError).Unwrap())
		return
	}

	c.Status(http.StatusOK)
}

func (h *financialGroupHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	financialGroup, err := h.financialGroupService.GetByID(context.Background(), id, userID)
	if err != nil {
		c.JSON(err.(*server_errors.SError).Unwrap())
		return
	}

	result := dto.CreateResponse[dto.FinancialGroup]{Result: *financialGroup}

	c.JSON(http.StatusOK, result)
}

func (h *financialGroupHandler) List(c *gin.Context) {
	var input dto.FinancialGroupListRequest
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	groupsList, err := h.financialGroupService.List(context.Background(), input, userID)
	if err != nil {
		c.JSON(err.(*server_errors.SError).Unwrap())
		return
	}

	c.JSON(http.StatusOK, groupsList)
}

func (h *financialGroupHandler) Delete(c *gin.Context) {

}
