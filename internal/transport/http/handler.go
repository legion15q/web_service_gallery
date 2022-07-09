package http

import (
	"net/http"
	"web_app/internal/config"
	"web_app/internal/service"

	"github.com/gorilla/mux"
	//httpSwagger "github.com/swaggo/http-swagger"
	handlers "web_app/internal/transport/http/handlers"
)

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) Init(cfg *config.Config) *mux.Router {
	// Init mux handler
	router := mux.NewRouter()
	//router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	// Init router
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK); w.Write([]byte("pong")) }).Methods("GET")
	hndlr := handlers.NewHandler(h.services)
	hndlr.InitHandlers(router)
	return router
}
