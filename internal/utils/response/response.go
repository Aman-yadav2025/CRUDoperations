package response

import (
	"CRUDapp/internal/types"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	StatusOK    = "OK"
	StatusError = "error"
)

func WriteJson(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) types.Response {
	return types.Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}

func ValidationError(errs validator.ValidationErrors) types.Response {
	var errMsgs []string

	for _,err := range errs {
		switch err.ActualTag() {
		case "required" :
			errMsgs = append(errMsgs,fmt.Sprintf("fill %s is a required field", err.Field()))
		default:
			errMsgs = append(errMsgs,fmt.Sprintf("field %s is invalid",err.Field()))
		}
	}
	return types.Response{
		Status: StatusError,
		Error: strings.Join(errMsgs,","),
	}
}
