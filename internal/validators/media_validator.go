package validators

import (
	"log"

	"github.com/go-playground/validator/v10"
	"shirinec.com/internal/enums"
)

func mediaUploadBindValidator(fl validator.FieldLevel) bool {
    log.Println("Validating media")
    value := fl.Field().String()
    log.Printf("%+v\n", value)
    validValues := []enums.MediaUploadBind{enums.BindToItem, enums.BindToProfile, enums.BindToCategory}

    for _,v := range validValues{
        if value == string(v){
            return true
        }
    }
    return false
}
