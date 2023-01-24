package handler

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) handleGET(respWr http.ResponseWriter, req *http.Request) {
	words := h.db.FindAll()
	buffer, _ := json.Marshal(words)
	respWr.Write(buffer)
	h.lgr.Write("get successfully, return all content")
}
