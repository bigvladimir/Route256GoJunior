package mdlware

import (
	"fmt"
	"homework/internal/pkg/jsonutil"
	"net/http"
)

func Auth(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		reqUsername, reqPassword, ok := req.BasicAuth()
		if !ok {
			http.Error(w, "Не удалось извлечь данные авторизации из запроса.", http.StatusUnauthorized)
			return
		}

		data, err := jsonutil.ReadJSON(make(map[string]interface{}), usersPAth)
		if err != nil {
			err = fmt.Errorf("Ошибка при чтении базы из данных пользователей: %w", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		bdPassword, ok := data.(map[string]interface{})[reqUsername]
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
