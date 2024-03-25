package mdlware

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

// middleware логирование
func Logger(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Printf("Метод: %s, путь: %s\n", req.Method, req.URL.Path)

		body, err := io.ReadAll(req.Body)
		if err != nil {
			log.Println("Ошибка при чтении тела запроса.", err)
		} else if len(body) != 0 {
			log.Printf("Тело запроса: %s\n", body)
		}
		req.Body = io.NopCloser(bytes.NewBuffer(body))
		handler.ServeHTTP(w, req)
	}
}
