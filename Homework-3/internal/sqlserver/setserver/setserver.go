package setserver

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"homework/internal/sqlserver/db"
	"homework/internal/sqlserver/mdlware"
	"homework/internal/sqlserver/repository"
	"homework/internal/sqlserver/server"
)

func SetServer() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDb(ctx)
	if err != nil {
		return err
	}
	defer database.GetPool(ctx).Close()
	serv := server.NewServer(repository.NewPvzStore(database))

	http.Handle("/", mdlware.Logger(mdlware.Auth(server.CreateRouter(serv))))

	// если будет не настроен сертификат и будет ошибка, то программа не закроется и можно будет обращаться по http
	go func() {
		if err := http.ListenAndServeTLS(httpsPort, crtPath, keyPath, nil); err != nil {
			log.Println("Ошибка при запуске обработчика https:", err, "\n Можно пользоваться только http запросами.")
		}
	}()

	if err := http.ListenAndServe(httpPort, nil); err != nil {
		err = fmt.Errorf("Ошибка при запуске обработчика http: %w", err)
		return err
	}

	return nil
}
