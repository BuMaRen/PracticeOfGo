package handler

import (
	"encoding/json"
	"goquickstart/database"
	"io/ioutil"
	"net/http"
)

func (h *Handler) handlePOST(respWr http.ResponseWriter, req *http.Request) {
	buffer, err := ioutil.ReadAll(req.Body)

	// TODO: 为什么这个读出来的EOF？？？我不理解，淦
	// _, err := req.Body.Read(buffer)
	if err != nil {
		h.lgr.Write("POST return 400 as body read err:" + err.Error())
		respWr.WriteHeader(http.StatusBadRequest)
		respWr.Write([]byte{})
		return
	}
	word := database.Words{}
	err = json.Unmarshal(buffer, &word)
	if err != nil {
		h.lgr.WriteS("POST return 400 as unmarshal failed: %s", err.Error())
		respWr.WriteHeader(http.StatusBadRequest)
		respWr.Write([]byte{})
		return
	}
	err = h.db.Insert(word.NewBson())
	if err != nil {
		h.lgr.Write("POST return 400 as insert failed")
		respWr.WriteHeader(http.StatusBadRequest)
		respWr.Write([]byte(err.Error()))
		return
	}
	respWr.WriteHeader(http.StatusOK)
	respWr.Write([]byte{})

	h.lgr.WriteS("post word %s successfully", word.Word)
}
