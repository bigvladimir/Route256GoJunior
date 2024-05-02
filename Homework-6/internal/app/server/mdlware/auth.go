package mdlware

import (
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

func Auth(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		reqUsername, reqPassword, ok := req.BasicAuth()
		if !ok {
			http.Error(w, "Не удалось извлечь данные авторизации из запроса.", http.StatusUnauthorized)
			return
		}

		yamlFile, err := os.ReadFile("./config/users.yaml")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		users := make(map[string]string)
		if err := yaml.Unmarshal(yamlFile, &users); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		bdPassword, ok := users[reqUsername]
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
