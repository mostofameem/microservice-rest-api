package web

import (
	"net/http"
	"order_service/web/handlers"
	"order_service/web/middlewares"
)

func InitRouts(mux *http.ServeMux, manager *middlewares.Manager) {
	mux.Handle(
		"POST /addorder",
		manager.With(
			http.HandlerFunc(handlers.AddToOrderList),
			//middlewares.AuthenticateJWT,
		),
	)
	mux.Handle(
		"GET /newcart",
		manager.With(
			http.HandlerFunc(handlers.NewCart),
			//middlewares.AuthenticateJWT,
		),
	)
	mux.Handle(
		"POST /updateorderstatus/{id}",
		manager.With(
			http.HandlerFunc(handlers.UpdateOrderInfo),
			//middlewares.AuthenticateJWT,
		),
	)
	mux.Handle(
		"GET /getorderdetails/{id}",
		manager.With(
			http.HandlerFunc(handlers.GetOrderDetails),
			//middlewares.AuthenticateJWT,
		),
	)
	mux.Handle(
		"POST /addproduct",
		manager.With(
			http.HandlerFunc(handlers.InsertProduct),
			//middlewares.AuthenticateJWT,
		),
	)
}
