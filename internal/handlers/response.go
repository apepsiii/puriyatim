package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// JSONResponse adalah struktur standar untuk semua JSON response dari handler.
type JSONResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// JSONOk mengirim response JSON 200 OK dengan data opsional.
func JSONOk(c echo.Context, message string, data ...interface{}) error {
	resp := JSONResponse{Success: true, Message: message}
	if len(data) > 0 {
		resp.Data = data[0]
	}
	return c.JSON(http.StatusOK, resp)
}

// JSONCreated mengirim response JSON 201 Created.
func JSONCreated(c echo.Context, message string, data ...interface{}) error {
	resp := JSONResponse{Success: true, Message: message}
	if len(data) > 0 {
		resp.Data = data[0]
	}
	return c.JSON(http.StatusCreated, resp)
}

// JSONBadRequest mengirim response JSON 400 Bad Request.
func JSONBadRequest(c echo.Context, message string) error {
	return c.JSON(http.StatusBadRequest, JSONResponse{Success: false, Message: message})
}

// JSONNotFound mengirim response JSON 404 Not Found.
func JSONNotFound(c echo.Context, message string) error {
	return c.JSON(http.StatusNotFound, JSONResponse{Success: false, Message: message})
}

// JSONInternalError mengirim response JSON 500 Internal Server Error.
func JSONInternalError(c echo.Context, message string) error {
	return c.JSON(http.StatusInternalServerError, JSONResponse{Success: false, Message: message})
}

// JSONUnauthorized mengirim response JSON 401 Unauthorized.
func JSONUnauthorized(c echo.Context, message string) error {
	return c.JSON(http.StatusUnauthorized, JSONResponse{Success: false, Message: message})
}

// JSONWithFields mengirim response JSON 200 OK dengan map fields tambahan.
// Berguna untuk response yang memiliki banyak field berbeda.
func JSONWithFields(c echo.Context, fields map[string]interface{}) error {
	if _, ok := fields["success"]; !ok {
		fields["success"] = true
	}
	return c.JSON(http.StatusOK, fields)
}
