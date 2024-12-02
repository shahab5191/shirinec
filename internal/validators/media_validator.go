package validators

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"shirinec.com/internal/enums"
)

func mediaUploadBindValidator(fl validator.FieldLevel) bool {
    value := fl.Field().String()
    validValues := []enums.MediaUploadBind{enums.BindToItem, enums.BindToProfile, enums.BindToCategory}

    for _,v := range validValues{
        if value == string(v){
            return true
        }
    }
    return false
}

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
