package response

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseBody struct {
	StatusCode int
	Message    string
	DevMessage string
	Body       any
}

type importResponseBody struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Body       interface{} `json:"body"`
	DevMessage error       `json:"dev_message"`
}

func ValidationResponse(c *gin.Context, Message string) {
	response := ResponseBody{
		StatusCode: http.StatusUnprocessableEntity,
		Message:    Message,
		Body:       map[string]interface{}{},
	}
	c.JSON(http.StatusUnprocessableEntity, response)
}

func BadRequestResponse(c *gin.Context, Message, err string) {
	response := ResponseBody{
		StatusCode: http.StatusBadRequest,
		Message:    Message,
		Body:       nil,
		DevMessage: err,
	}
	c.JSON(http.StatusBadRequest, response)
}

func InternalServerErrorResponse(c *gin.Context, Message, err string) {
	response := ResponseBody{
		StatusCode: http.StatusInternalServerError,
		Message:    Message,
		Body:       nil,
		DevMessage: err,
	}
	c.JSON(http.StatusInternalServerError, response)
}

func ValidationResponseWithError(c *gin.Context, Message, err string) {
	response := ResponseBody{
		StatusCode: http.StatusUnprocessableEntity,
		Message:    Message,
		Body:       nil,
		DevMessage: err,
	}
	c.JSON(http.StatusUnprocessableEntity, response)
}

func InternalServerErrorImportResponse(c *gin.Context, Error error, message string) {
	errorResponse := importResponseBody{
		StatusCode: 400,
		Message:    message,
		DevMessage: Error,
		Body:       map[string]interface{}{},
	}
	c.JSON(http.StatusUnprocessableEntity, errorResponse)
}

func SuccessResponse(c *gin.Context, Message string, Body interface{}) {
	response := ResponseBody{
		StatusCode: http.StatusOK,
		Message:    Message,
		Body:       Body,
	}
	c.JSON(http.StatusOK, response)
}

func InternalServerErrorResponseErr(c *gin.Context, err error) {
	response := ResponseBody{
		StatusCode: http.StatusInternalServerError,
		Message:    "Internal Server Error",
		Body:       nil,
		DevMessage: err.Error(),
	}
	c.JSON(http.StatusInternalServerError, response)
}

func RequestBodyLogger(c *gin.Context) string {
	requestBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return ""
	}
	// Reset Request.Body for downstream handlers
	c.Request.Body = io.NopCloser(bytes.NewReader(requestBody))
	return string(requestBody)
}

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	s := buf.String()
	return s
}

func UnAuthorizedResponse(c *gin.Context, err error, Message string) error {
	errorResponse := ResponseBody{
		DevMessage: err.Error(),
		StatusCode: 403,
		Message:    Message,
		Body:       map[string]interface{}{},
	}
	c.JSON(http.StatusUnauthorized, errorResponse)
	return fmt.Errorf("%s", Message)
}
