package httpError

import (
	"cine-circle-api/pkg/customError"
	"cine-circle-api/pkg/logger"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"net/http"
	"runtime"
)

// FormattedJsonError Object returned by web services in case of error
type FormattedJsonError struct {
	Message   string                 `json:"message,omitempty"`
	ErrorCode string                 `json:"code"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// Package 'errors' does not expose its interface, so we have to declare it
type stackTracer interface {
	StackTrace() errors.StackTrace
}

// HandleHTTPError Write readable and formatted error in the restful.Response
func HandleHTTPError(req *restful.Request, res *restful.Response, err error) {
	// Instantiate the result object that will be displayed in HTTP response
	jsonFormattedError := FormattedJsonError{
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

	// Manage custom error typed

	if e, ok := err.(*customError.CustomError); ok {

		jsonFormattedError.ErrorCode = e.Code()

		message := fmt.Sprintf("%s\t%s : returned custom error %s", req.Request.Method, req.Request.URL.String(), e.Error())

		if e.HttpStatus() >= 500 {
			logger.Errorf(message)
		} else {
			logger.Infof(message)
		}

		res.WriteHeaderAndEntity(e.HttpStatus(), jsonFormattedError)

		return
	}

	// Manage various errors, possibly wrapped

	// Postgres

	if errors.Is(err, gorm.ErrRecordNotFound) {

		jsonFormattedError.Message = fmt.Sprintf("Entity not found in database: %s", err.Error())

		logger.Infof("%s\t%s returned Gorm error (%d) : %s", req.Request.Method, req.Request.URL.String(), http.StatusNotFound, err.Error())

		res.WriteHeaderAndEntity(http.StatusNotFound, jsonFormattedError)
		return
	}

	// Unknown errors

	jsonFormattedError.Message = fmt.Sprintf("Something went wrong : '%s'. Please report this error with the following stack to the IT crew", err.Error())

	logger.Errorf("%s\t%s : error occurs (%d) : %s", req.Request.Method, req.Request.URL.String(), http.StatusInternalServerError, err.Error())

	res.WriteHeaderAndEntity(http.StatusInternalServerError, jsonFormattedError)
}
