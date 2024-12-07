package validators

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"shirinec.com/internal/enums"
)

func categoryCreateTypeValidator(fl validator.FieldLevel) bool {
    value := fl.Field().String()
    validValues := []enums.CategoryType{enums.Income, enums.Account, enums.Expense}

    for _,v := range validValues {
        if value == strings.ToLower(string(v)){
            return true
        }
    }
    return false
}
