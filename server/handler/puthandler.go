package handler

import (
	"encoding/json"
	"goquickstart/database"
	"io/ioutil"
	"net/http"
)

func (h *Handler) handlePUT(respWr http.ResponseWriter, req *http.Request) {

	buffer, err := ioutil.ReadAll(req.Body)

	if err != nil {
		h.lgr.WriteS("PUT return 400 as body read err: ", err.Error())
		respWr.WriteHeader(http.StatusBadRequest)
		respWr.Write([]byte{})
		return
	}
	word := database.Words{}
	err = json.Unmarshal(buffer, &word)
	if err != nil {
		h.lgr.WriteS("PUT return 400 as unmarshal failed: ", err.Error())
		respWr.WriteHeader(http.StatusBadRequest)
		respWr.Write([]byte{})
		return
	}
	h.db.UpdateWord(word.Word, word)
	if err != nil {
		h.lgr.WriteS("PUT return 400 as insert failed: ", err.Error())
		respWr.WriteHeader(http.StatusBadRequest)
		respWr.Write([]byte{})
		return
	}
	respWr.WriteHeader(http.StatusOK)
	respWr.Write([]byte{})

	h.lgr.WriteS("put word %s successfully", word.Word)
}
