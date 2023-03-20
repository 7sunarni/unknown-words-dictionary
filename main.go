package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/dictionary", dictionary)
	http.HandleFunc("/remember", remember)
	http.HandleFunc("/history", searched)
	http.ListenAndServe(":18080", nil)
}

func dictionary(rw http.ResponseWriter, r *http.Request) {
	word := r.URL.Query().Get("word")
	if word == "" {
		rw.Write([]byte(""))
		return
	}
	dictionary, err := getWord(word)
	if err != nil {
		rw.Write([]byte(err.Error()))
		return
	}

	if err := Put(dictionary); err != nil {
		rw.Write([]byte(err.Error()))
		return
	}

	WrapHTML(rw, dictionary)
}

func remember(rw http.ResponseWriter, r *http.Request) {
	ret, err := TodayNeedRemember()
	if err != nil {
		rw.Write([]byte(err.Error()))
		return
	}
	WrapHTML(rw, ret)
}

func searched(rw http.ResponseWriter, r *http.Request) {
	ret, err := TodySearched()
	if err != nil {
		rw.Write([]byte(err.Error()))
		return
	}
	WrapHTML(rw, ret)
}
