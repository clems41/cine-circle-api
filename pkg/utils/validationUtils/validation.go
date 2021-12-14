package validationUtils

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/go-playground/validator/v10"
	"strings"
)

var (
	validate *validator.Validate
)

func init() {
	validate = validator.New()
}

// ReadEntityAndValidateStruct read request body and unmarshall it into specified out interface.
// Then, we check that none of required fields are missing.
func ReadEntityAndValidateStruct(req *restful.Request, out interface{}) (err error) {
	err = req.ReadEntity(out)
	if err != nil {
		return err
	}
	err = validate.Struct(out)
	if err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			var incorrectFields []string
			for _, validationErr := range validationErrors {
				incorrectFields = append(incorrectFields, validationErr.Field())
			}
			return fmt.Errorf("structure is not correct, some fields are incorrects : %s", strings.Join(incorrectFields, ","))
		}
	}
	return
}