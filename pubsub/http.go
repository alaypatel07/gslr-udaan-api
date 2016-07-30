package pubsub

import (
	"net/http"
	"encoding/json"
)

type HttpCon struct {
	http.ResponseWriter
}

func NewHttpCon(rw http.ResponseWriter) *HttpCon{
	return &HttpCon{rw}
}

func (hc HttpCon)Consume(data interface{})  {
	e := json.NewEncoder(hc)
	e.Encode(data)
}
