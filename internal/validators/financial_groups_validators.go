package validators

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"shirinec.com/internal/enums"
)

func financialRoleValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	validValues := []enums.FinancialGroupRole{enums.FinancialGroupOwner, enums.FinancialGroupMember}

	for _, v := range validValues {
		if value == strings.ToLower(string(v)) {
			return true
		}
	}
	return false
}
