package validators

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"shirinec.com/src/internal/enums"
)

func accountTypeValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	validValues := []enums.AccountType{enums.AccountTypeSelf, enums.AccoutTypeExternal}

	for _, v := range validValues {
		if value == strings.ToLower(string(v)) {
			return true
		}
	}

	return false
}
