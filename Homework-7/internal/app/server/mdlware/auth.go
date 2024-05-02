package mdlware

import (
	"net/http"

	"homework/internal/pkg/config"
)

// Auth is the HTTP authentication middleware
func Auth(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		reqUsername, reqPassword, ok := req.BasicAuth()
		if !ok {
			http.Error(w, "Не удалось извлечь данные авторизации из запроса.", http.StatusUnauthorized)
			return
		}

		bdPassword, ok := config.Cfg().GetUserPassword(reqUsername)
		if !ok {
			http.Error(w, "Пользователь не найден.", http.StatusUnauthorized)
			return
		}
		if reqPassword != bdPassword {
			http.Error(w, "Неправильный пароль.", http.StatusUnauthorized)
			return
		}

		handler.ServeHTTP(w, req)
	}
}
