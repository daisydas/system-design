package main

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"system-design/kv/data"
)

type Handler struct {
	*data.Operation
}

func (h *Handler) GetKeyHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	val := h.Operation.GetFromKVStore(key)
	if val == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, err := w.Write([]byte(val))
	if err != nil {
		fmt.Printf("error writing response :%v\n", err)
		return
	}

}

type keyValue struct {
	Key string `json:"key"`
	Val string `json:"val"`
	TTL int    `json:"ttl"`
}

func (h *Handler) PutKeyHandler(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer body.Close()
	var kv keyValue

	jsoniter.NewDecoder(body).Decode(&kv)

	if kv.Key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := h.Operation.PutToKVStore(kv.Key, kv.Val, kv.TTL)
	if err != nil {
		fmt.Printf("error putting key value :%v\n", err)
		return
	}
	return
}

func (h *Handler) DeleteKeyHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err := h.Operation.DeleteFromKVStore(key)
	if err != nil {
		fmt.Printf("error writing response :%v\n", err)
		return
	}
}
