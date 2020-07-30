package handler

import (
	"errors"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func parse(r *http.Request) (request, error) {
	// We always need to read and close the request body.
	_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return request{}, errors.New("unable to read request body")
	}
	_ = r.Body.Close()

	log.Printf("r.Header: %+v\n", r.Header["Authorization"])
	var req request

	switch r.Method {
	case "GET":
		req = request{queries: r.URL.Query()}
	default:
		err = errors.New("only GET requests are supported")
	}

	log.Printf("parsed request: %+v\n", req)

	return req, err
}
