package webService

import (
	"cine-circle/internal/typedErrors"
	"cine-circle/pkg/logger"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"net/http"
	"runtime"
)

// FormattedJsonError Object returned by web services in case of error
type FormattedJsonError struct {
	Kind     string                 `json:"kind"`
	Message  string                 `json:"message,omitempty"`
	Code     string                 `json:"code"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

func (err FormattedJsonError) Error() string {
	return fmt.Sprintf("kind:%s error:%s", err.Kind, err.Message)
}

// Package 'errors' does not expose its interface, so we have to declare it
type stackTracer interface {
	StackTrace() errors.StackTrace
}

// HandleHTTPError Write readable and formatted error in the restful.Response
func HandleHTTPError(req *restful.Request, res *restful.Response, err error) {

	// Instantiate the result object that will be displayed in HTTP response
	jsonFormattedError := FormattedJsonError{
		Kind:     "",
		Message:  err.Error(),
		Metadata: make(map[string]interface{}, 1),
	}

	// Stack trace builder
	if e, ok := err.(stackTracer); ok {

		// A stack is attached, let's show it
		jsonFormattedError.Metadata["Stack"] = fmt.Sprintf("%+v", e.StackTrace()) // Get all frames

	} else {

		// If no stack attached, this will just return the calls from webservice
		buf := make([]byte, 1<<16)
		stackLength := runtime.Stack(buf, false)
		jsonFormattedError.Metadata["Stack"] = string(buf[:stackLength])
	}

	// Manage DPC error typed

	if e, ok := err.(typedErrors.TypedError); ok {

		jsonFormattedError.Kind = e.Type()
		jsonFormattedError.Code = e.BusinessCode()

		err = res.WriteHeaderAndEntity(e.HttpStatus(), jsonFormattedError)
		if err != nil {
			logger.Sugar.Errorf("Cannot send response : %+v", err)
		}

		return
	}

	// Manage various errors, possibly wrapped

	// Postgres

	if errors.Is(err, gorm.ErrRecordNotFound) {

		logger.Sugar.Debugf("Entity not found in database: %+v", err)

		jsonFormattedError.Kind = "Postgres"

		err = res.WriteHeaderAndEntity(http.StatusNotFound, jsonFormattedError)
		if err != nil {
			logger.Sugar.Errorf("Cannot send response : %+v", err)
		}

		jsonFormattedError.Message = fmt.Sprintf("Entity not found in database: %+v", err)
		jsonFormattedError.Kind = "Postgres"

		logger.Sugar.Errorf(jsonFormattedError.Message)

		err = res.WriteHeaderAndEntity(http.StatusNotFound, jsonFormattedError)
		if err != nil {
			logger.Sugar.Errorf("Cannot send response : %+v", err)
		}
		return
	}

	// Unknown errors

	jsonFormattedError.Kind = "Internal server error"
	jsonFormattedError.Message = fmt.Sprintf("Something went wrong : '%s'. Please report this error with the following stack to the IT crew", err.Error())

	logger.Sugar.Errorf("%s - %s : error occurs : %+v", req.Request.Method, req.Request.URL.String(), err)

	err = res.WriteHeaderAndEntity(http.StatusInternalServerError, jsonFormattedError)
	if err != nil {
		logger.Sugar.Errorf("Cannot send response : %+v", err)
	}
}
