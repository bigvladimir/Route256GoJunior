package main

import (
	"context"
	"homework/internal/app/core"
	"homework/internal/app/orders"
	"homework/internal/app/pvz"
	"homework/internal/app/server"
	"homework/internal/app/server/mdlware"
	"homework/internal/pkg/db"
	"log"
	"net/http"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDb(ctx, bdCfgPath)
	if err != nil {
		log.Fatal("Ошибка при подключении к базе данных:", err)
	}
	defer database.GetPool(ctx).Close()

	// http вариант приложения работает только с пвз, но пока не буду переписывать инициализацию сервиса
	// чтобы не дописывать ещё один интерфейс без заказов
	// можно будет потом расширить функционал http приложения для обработки заказов
	packVariants := map[string]orders.PackageVariant{
		orders.BagVariantName:  orders.NewBagPackage(),
		orders.BoxVariantName:  orders.NewBoxPackage(),
		orders.FilmVariantName: orders.NewFilmPackage(),
	}
	orderService := orders.NewService(orders.NewOrderStorage(database), packVariants)
	pvzService := pvz.NewService(pvz.NewPvzStorage(database))
	service := core.NewCoreService(orderService, pvzService)
	serv := server.NewServer(service)

	http.Handle("/", mdlware.Logger(mdlware.Auth(server.CreateRouter(serv))))

	// если будет не настроен сертификат и будет ошибка, то программа не закроется и можно будет обращаться по http
	go func() {
		if err := http.ListenAndServeTLS(httpsPort, crtPath, keyPath, nil); err != nil {
			log.Println("Ошибка при запуске обработчика https:", err, "\n Можно пользоваться только http запросами.")
		}
	}()

	if err := http.ListenAndServe(httpPort, nil); err != nil {
		log.Fatal("Ошибка при запуске обработчика http:", err)
	}
}
