package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type StatusError struct {
	Code int
	Err  error
	Msg  string
}

type Error interface {
	error
	Status() int
}

func (se StatusError) Error() string {
	return se.Err.Error()
}

func (se StatusError) Status() int {
	return se.Code
}

type Response interface{}

type Res struct {
	Data string `json:"data"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func NewRes(data string) Response {
	return Res{Data: data, Code: 0, Msg: ""}
}

func errorHandler(w http.ResponseWriter, status int, msgOpt ...string) {
	msg := http.StatusText(status)
	w.WriteHeader(status)

	if len(msgOpt) > 0 {
		msg = strings.Join(msgOpt, ",")
	}
	res := Res{"", status, msg}
	jsonRes, _ := json.Marshal(res)
	w.Write(jsonRes)
}

type Handler struct {
	View func(w http.ResponseWriter, r *http.Request) error
}

func (h Handler) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.View(w, r)
	if err != nil {
		switch e := err.(type) {
		case Error:
			log.Printf("HTTP %d - %s", e.Status(), e)
			http.Error(w, e.Error(), e.Status())
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
