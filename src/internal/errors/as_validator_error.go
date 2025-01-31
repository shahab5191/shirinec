package server_errors

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/go-playground/validator/v10"
)

func AsValidatorError(err error) *[]string {
	if err != nil {
		var validationErr validator.ValidationErrors
		if errors.As(err, &validationErr) {
			var errList []string
			for _, err := range validationErr {
				switch err.Tag() {
				case "email":
					errList = append(errList, fmt.Sprintf("%s is not in correct format", err.Field()))
				case "min":
					length, _ := strconv.Atoi(err.Param())
					errList = append(
						errList,
						fmt.Sprintf(
							"%s length should be longer than %d",
							err.Field(),
							length,
						),
					)
				case "len":
					length, _ := strconv.Atoi(err.Param())
					errList = append(
						errList,
						fmt.Sprintf(
							"%s length should be exactly %d",
							err.Field(),
							length,
						),
					)
				case "intLen":
					length, _ := strconv.Atoi(err.Param())
					errList = append(
						errList,
						fmt.Sprintf(
							"%s length should be exactly %d",
							err.Field(),
							length,
						),
					)
				case "jwt":
					errList = append(errList, fmt.Sprintf("%s is not a correct jwt", err.Field()))
				case "required":
					errList = append(errList, fmt.Sprintf("%s field is required", err.Field()))
				case "hexcolor":
					errList = append(errList, fmt.Sprintf("%s field should be a hex color", err.Field()))
				case "categoryCreateType":
					errList = append(errList, fmt.Sprintf("%s field should be 'income', 'expense' or 'account'", err.Field()))
				case "alphanum":
					errList = append(errList, fmt.Sprintf("%s field must contain only letters and numbers characters", err.Field()))
				case "alphaNumericSpace":
					errList = append(errList, fmt.Sprintf("%s field must contain only letters, numbers, or spaces.", err.Field()))
				case "mediaUploadBind":
					errList = append(errList, "binds_to should be 'item', 'profile' or 'category'")
				default:
					log.Printf("[Error] - AsValidatorError - Undefined error tag: %+v\n", err)
					errList = append(errList, err.Error())
				}
			}
			return &errList
		}
	}
	return nil
}
