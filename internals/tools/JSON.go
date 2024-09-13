package tools

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

// ReadStringToJSON takes a string and embeds it to a given structure
func ReadStringToJSON(payload string, dst interface{}) error {

	err := json.Unmarshal([]byte(payload), &dst)
	if err != nil {
		return err
	}

	return nil
}

func ReadWebsocketJSON(conn *websocket.Conn, data any) error {

	decoder := json.NewDecoder(conn.Request().Body)
	err := decoder.Decode(&data)
	if err != nil {
		return err
	}

	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only have a single json value")
	}
	return nil
}

// ReadJSON reads the body of an HTTP request and decodes it into a structure
func ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := int64(1048576)

	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(data)
	if err != nil {
		return err
	}

	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only have a single json value")
	}
	return nil
}

// WriteJSON writes any payload in to the responsewriter
func WriteJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)

	return nil

}
