package utils

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// ResponseBody defines the common structure for JSON responses.
type ResponseBody struct {
	Message    string `json:"message"`
	StatusCode int    `json:"Status_code"`
	// DevMessage error       `json:"-"`
	DevMessage string      `json:"dev_message,omitempty"`
	Body       interface{} `json:"body,omitempty"`
}

type MiddlewareResponseBody struct {
	Message    string
	StatusCode int
	Body       struct{}
}

// Generic function to handle JSON responses
func sendJSONResponse(c *fiber.Ctx, statusCode int, message string, body interface{}) error {

	response := ResponseBody{
		Message:    message,
		StatusCode: statusCode,
		Body:       body,
	}

	return c.Status(statusCode).JSON(response)
}

// Generic function to abort with JSON response in middleware
func abortWithJSON(c *fiber.Ctx, statusCode int, message string, body interface{}) {
	response := ResponseBody{
		StatusCode: statusCode,
		Message:    message,
		Body:       body,
	}
	c.Status(statusCode).JSON(response)

}

// Success Response (200 OK)
func SuccessResponse(c *fiber.Ctx, message string, body interface{}) error {

	return sendJSONResponse(c, http.StatusOK, message, body)
}

// Created Response (201 Created)
func CreatedResponse(c *fiber.Ctx, message string, body interface{}) {
	sendJSONResponse(c, http.StatusCreated, message, body)
}

// Accepted Response (202 Accepted)
func AcceptedResponse(c *fiber.Ctx, message string, body interface{}) {
	sendJSONResponse(c, http.StatusAccepted, message, body)
}

// No Content Response (204 No Content)
func NoContentResponse(c *fiber.Ctx, message string) {
	sendJSONResponse(c, http.StatusNoContent, message, nil)
}

// No Content Abort Response (204 No Content)
func NoContentAbortWithJSON(c *fiber.Ctx, message string) {
	abortWithJSON(c, http.StatusNoContent, message, nil)
}

// Validation Error Response (422 Unprocessable Entity)
func ValidationResponse(c *fiber.Ctx, message string) error {

	return sendJSONResponse(c, http.StatusUnprocessableEntity, message, map[string]interface{}{})
}

// Validation Abort Response (422 Unprocessable Entity)
func ValidationAbortWithJSON(c *fiber.Ctx, message string) {
	abortWithJSON(c, http.StatusUnprocessableEntity, message, map[string]interface{}{})
}

// Bad Request Response (400 Bad Request)
func BadRequestResponse(c *fiber.Ctx, message string) {
	sendJSONResponse(c, http.StatusBadRequest, message, nil)
}

// Bad Request Abort Response (400 Bad Request)
func BadRequestAbortWithJSON(c *fiber.Ctx, message string) {
	abortWithJSON(c, http.StatusBadRequest, message, nil)
}

// Unauthorized Response (401 Unauthorized)
func UnauthorizedResponse(c *fiber.Ctx, message string) {
	sendJSONResponse(c, http.StatusUnauthorized, message, nil)
}

// Unauthorized Abort Response (401 Unauthorized)
func UnauthorizedAbortWithJSON(c *fiber.Ctx, message string) {
	abortWithJSON(c, http.StatusUnauthorized, message, nil)
}

// Forbidden Response (403 Forbidden)
func ForbiddenResponse(c *fiber.Ctx, message string) {
	sendJSONResponse(c, http.StatusForbidden, message, nil)
}

// Forbidden Abort Response (403 Forbidden)
func ForbiddenAbortWithJSON(c *fiber.Ctx, message string) {
	abortWithJSON(c, http.StatusForbidden, message, nil)
}

// Not Found Response (404 Not Found)
func NotFoundResponse(c *fiber.Ctx, message string) error {

	return sendJSONResponse(c, http.StatusNotFound, message, map[string]interface{}{})
}

// Not Found Abort Response (404 Not Found)
func NotFoundAbortWithJSON(c *fiber.Ctx, message string) {
	abortWithJSON(c, http.StatusNotFound, message, map[string]interface{}{})
}

// Internal Server Error Response (500 Internal Server Error)
func InternalServerErrorResponse(c *fiber.Ctx, err error) {
	response := ResponseBody{
		StatusCode: http.StatusInternalServerError,
		Message:    "Internal Server Error",
		DevMessage: err.Error(),
	}
	c.Status(http.StatusInternalServerError).JSON(response)
}

// Internal Server Error Abort Response (500 Internal Server Error)
func InternalServerErrorAbortWithJSON(c *fiber.Ctx, err error) {
	response := ResponseBody{
		StatusCode: http.StatusInternalServerError,
		Message:    "Internal Server Error",
		DevMessage: err.Error(),
	}
	c.Status(http.StatusInternalServerError).JSON(response)
}

// Service Unavailable Response (503 Service Unavailable)
func ServiceUnavailableResponse(c *fiber.Ctx, message string) {
	sendJSONResponse(c, http.StatusServiceUnavailable, message, nil)
}

// Service Unavailable Abort Response (503 Service Unavailable)
func ServiceUnavailableAbortWithJSON(c *fiber.Ctx, message string) {
	abortWithJSON(c, http.StatusServiceUnavailable, message, nil)
}

// Gateway Timeout Response (504 Gateway Timeout)
func GatewayTimeoutResponse(c *fiber.Ctx, message string) {
	sendJSONResponse(c, http.StatusGatewayTimeout, message, nil)
}

// Gateway Timeout Abort Response (504 Gateway Timeout)
func GatewayTimeoutAbortWithJSON(c *fiber.Ctx, message string) {
	abortWithJSON(c, http.StatusGatewayTimeout, message, nil)
}

// Internal Server Error with custom message (500 Internal Server Error)
func InternalServerErrorWithMessage(c *fiber.Ctx, message string) error {
	return c.Status(http.StatusInternalServerError).JSON(ResponseBody{
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Body:       map[string]interface{}{},
	})
}

// Internal Server Error Abort with custom message (500 Internal Server Error)
func InternalServerErrorAbortWithMessage(c *fiber.Ctx, message string) {
	abortWithJSON(c, http.StatusInternalServerError, message, map[string]interface{}{})
}

// Too Many Requests Response (429 Too Many Requests)
func TooManyRequestsResponse(c *fiber.Ctx, message string) {
	sendJSONResponse(c, http.StatusTooManyRequests, message, nil)
}

// Too Many Requests Abort Response (429 Too Many Requests)
func TooManyRequestsAbortWithJSON(c *fiber.Ctx, message string) {
	abortWithJSON(c, http.StatusTooManyRequests, message, nil)
}

func RollbackWithError(tx *gorm.DB, errMsg string, err error) error {
	if tx != nil {
		tx.Rollback()
	}
	logrus.Errorf("%s: %v", errMsg, err)
	return fmt.Errorf("%s: %w", errMsg, err)
}

func RespondWithError(c *fiber.Ctx, code int, message string) {
	response := MiddlewareResponseBody{
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}

	c.Status(http.StatusUnauthorized).JSON(response)
}

func BadRequestWithObject(c *fiber.Ctx, body interface{}) {

	c.Status(http.StatusBadRequest).JSON(fiber.Map{
		"status":  false,
		"message": body,
	})
}
