package app

import (
	"encoding/json"
	"net/http"
	handler "stocker/app/handlers"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Router - router struct
type Router struct {
	*mux.Router
}

// NewRouter - new router instance
func NewRouter() *Router {
	return &Router{mux.NewRouter()}
}

// InitializeRoutes ...
func (r *Router) InitializeRoutes() http.Handler {
	api := (*r)

	api.HandleFunc("/start", panicRecover(handler.CreatePlayer)).
		Methods(http.MethodPost)

	api.HandleFunc("/balance/{player}", panicRecover(handler.GetWalletBalance)).
		Methods(http.MethodGet)

	api.HandleFunc("/buy", panicRecover(handler.BuyBitcoin)).
		Methods(http.MethodPost)

	api.HandleFunc("/sell", panicRecover(handler.SellBitcoin)).
		Methods(http.MethodPost)

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowCredentials: true,
	}).Handler(api)
	return handler

}

// panicRecover - recover endpoint from panic
func panicRecover(restart func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				handler.ErrorLogger.Println(err)
				jsonBody, _ := json.Marshal(map[string]string{
					"error": "There was an internal server error",
				})
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(jsonBody)
			}
		}()
		restart(w, r)
	}
}
