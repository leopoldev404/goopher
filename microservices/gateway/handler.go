package main

import (
	"log"
	"net/http"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (handler *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/customers/{id}/orders", handler.HandleCreateOrder)
}

func (handler *Handler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	log.Print(r)
}
