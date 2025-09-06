package requests

import (
    "encoding/json"
    "net/http"
    "github.com/go-playground/validator/v10"
)

var validate = validator.New()

// Validate request body and bind it to a struct
func Validate(r *http.Request, requestStruct interface{}) error {
    if err := json.NewDecoder(r.Body).Decode(requestStruct); err != nil {
        return err
    }
    return validate.Struct(requestStruct)
}