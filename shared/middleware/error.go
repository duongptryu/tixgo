package middleware

import (
	"errors"
	"net/http"
	"tixgo/shared/response"
	"tixgo/shared/syserr"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			handleError(c, err)
		}
	}
}

func handleError(c *gin.Context, err error) {
	var sysErr *syserr.Error
	if errors.As(err, &sysErr) {
		// statusCode := getHTTPStatusCode(sysErr.Code())
		response.NewErrorResponse(
			string(sysErr.Code()),
			sysErr.Error(),
			nil,
		).JSON(c, http.StatusOK)
		return
	}

	// Default error
	response.NewErrorResponse(
		"internal_error",
		"An error occurred",
		nil,
	).JSON(c, http.StatusInternalServerError)
}
