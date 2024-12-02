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
                    errList = append(errList, "Email is not in correct format")
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
                case "jwt":
                    errList = append(errList, fmt.Sprintf("%s is not a correct jwt", err.Field()))
                case "required":
                    errList = append(errList, fmt.Sprintf("%s field is required", err.Field()))
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
