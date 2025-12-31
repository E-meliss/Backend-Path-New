package response

import (
	"net/http"
)

func Error(w http.ResponseWriter, r *http.Request, status int, code, message string) {
	JSON(w, status, map[string]any{
		"error": map[string]any{
			"code":    code,
			"message": message,
		},
	})
}
