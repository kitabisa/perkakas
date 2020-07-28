package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"

	"github.com/kitabisa/perkakas/v2/structs"
)

type HttpHandlerContext struct {
	M structs.Meta
	E map[error]*structs.ErrorResponse
}

func NewContextHandler(meta structs.Meta) HttpHandlerContext {
	var errMap map[error]*structs.ErrorResponse = map[error]*structs.ErrorResponse{
		// register general error here, so if there are new general error you must add it here
		structs.ErrInvalidHeader:          structs.ErrInvalidHeader,
		structs.ErrUnauthorized:           structs.ErrUnauthorized,
		structs.ErrInvalidHeaderSignature: structs.ErrInvalidHeaderSignature,
		structs.ErrInvalidHeaderTime:      structs.ErrInvalidHeaderTime,
	}

	return HttpHandlerContext{
		M: meta,
		E: errMap,
	}
}

func (hctx HttpHandlerContext) AddError(key error, value *structs.ErrorResponse) {
	hctx.E[key] = value
}

func (hctx HttpHandlerContext) AddErrorMap(errMap map[error]*structs.ErrorResponse) {
	for k, v := range errMap {
		hctx.E[k] = v
	}
}

type CustomWriter struct {
	C HttpHandlerContext
}

func (c *CustomWriter) Write(w http.ResponseWriter, data interface{}, nextPage *string) {
	var successResp structs.SuccessResponse
	voData := reflect.ValueOf(data)
	arrayData := []interface{}{}

	if voData.Kind() != reflect.Slice {
		if voData.IsValid() {
			arrayData = []interface{}{data}
		}
		successResp.Data = arrayData
	} else {
		if voData.Len() != 0 {
			successResp.Data = data
		} else {
			successResp.Data = arrayData
		}
	}

	successResp.ResponseCode = "000000"
	successResp.Next = nextPage
	successResp.Meta = c.C.M

	writeSuccessResponse(w, successResp)
}

// WriteError sending error response based on err type
func (c *CustomWriter) WriteError(w http.ResponseWriter, err error) {
	if len(c.C.E) > 0 {
		errorResponse := LookupError(c.C.E, err)
		if errorResponse == nil {
			errorResponse = structs.ErrUnknown
		}

		errorResponse.Meta = c.C.M
		writeErrorResponse(w, errorResponse)
	} else {
		var errorResponse *structs.ErrorResponse = &structs.ErrorResponse{}
		if errors.As(err, &errorResponse) {
			errorResponse.Meta = c.C.M
			writeErrorResponse(w, errorResponse)
		} else {
			errorResponse = structs.ErrUnknown
			errorResponse.Meta = c.C.M
			writeErrorResponse(w, errorResponse)
		}
	}
}

func writeResponse(w http.ResponseWriter, response interface{}, contentType string, httpStatus int) {
	res, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to unmarshal"))
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(httpStatus)
	w.Write(res)
}

func writeSuccessResponse(w http.ResponseWriter, response structs.SuccessResponse) {
	writeResponse(w, response, "application/json", http.StatusOK)
}

func writeErrorResponse(w http.ResponseWriter, errorResponse *structs.ErrorResponse) {
	writeResponse(w, errorResponse, "application/json", errorResponse.HttpStatus)
}

// LookupError will get error message based on error type, with variables if you want give dynamic message error
func LookupError(lookup map[error]*structs.ErrorResponse, err error) (res *structs.ErrorResponse) {
	if msg, ok := lookup[err]; ok {
		res = msg
	}

	return
}
