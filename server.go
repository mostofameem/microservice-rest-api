package web

import (
	"fmt"
	"log/slog"
	"net/http"
	"order_service/config"
	"order_service/web/middlewares"
	"order_service/web/swagger"
	"sync"
)

func StartServer(wg *sync.WaitGroup) {
	manager := middlewares.NewManager()
	mux := http.NewServeMux()

	InitRouts(mux, manager)
	handler := middlewares.EnableCors(mux)

	swagger.SetupSwagger(mux, manager)

	wg.Add(1)

	go func() {
		defer wg.Done()

		conf := config.GetConfig()

		addr := fmt.Sprintf(":%d", conf.HttpPort)
		slog.Info(fmt.Sprintf("Listening at %s", addr))

		if err := http.ListenAndServe(addr, handler); err != nil {
			slog.Error(err.Error())
		}
	}()

}
