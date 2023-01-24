package handler

import (
	"goquickstart/database"
	"goquickstart/logger"
	"net/http"
)

type Handler struct {
	db  database.DataBase
	lgr logger.Logger
}

func NewHandler(l logger.Logger, d database.DataBase) *Handler {
	return &Handler{lgr: l, db: d}
}

func (h *Handler) ServeHTTP(respWr http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		h.lgr.Write("process GET method")
		h.handleGET(respWr, req)
	case http.MethodPost:
		h.lgr.Write("process POST method")
		h.handlePOST(respWr, req)
	case http.MethodPut:
		h.lgr.Write("process PUT method")
		h.handlePUT(respWr, req)
	case http.MethodDelete:
	default:
	}
}
