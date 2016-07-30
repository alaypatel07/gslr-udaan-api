package pubsub

import (
	"net/http"
	"encoding/json"
	"log"
	"testing"
	"bytes"
	"io/ioutil"
	"time"
)

var p Publisher

func publishHandler(rw http.ResponseWriter, r *http.Request) {
	var data map[string]string
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&data); err != nil {
		log.Fatalln("Unmarshelling Request Publish Handler", err)
	}
	p = NewPublisher(data["channel"])
	p.Publish(data["data"])
	rw.Write([]byte("PUBLISHED"))
}

func subscribeHandler(rw http.ResponseWriter, r *http.Request) {
	e := json.NewDecoder(r.Body)
	var data map[string]string
	err := e.Decode(data)
	if err != nil {
		log.Fatalln(err)
	}
	ch := NewChannel(data["channel"])
	con := NewHttpCon(rw)
	ch.Subscribe(con)
}

func startPublishServer() {
	http.HandleFunc("/publish", publishHandler)
	http.HandleFunc("/subscribe", subscribeHandler)
	log.Fatalln(http.ListenAndServe(":8081", nil))
}


func TestHttpCon_Consume(t *testing.T) {
	go startPublishServer()
	time.Sleep(1 * time.Second)
	req := make(map[string]string)
	req["channel"] = "foo"
	req["data"] = "bar"
	r, _ := json.Marshal(req)
	contentReader := bytes.NewReader(r)
	request, err := http.NewRequest("POST", "http://127.0.0.1:8081/publish", contentReader)
	if err != nil {
		log.Fatalln(err)
	}
	request.Header.Set("Content-Type", "application/json")
	c := http.Client{}
	resp, err := c.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
	rb, _ := ioutil.ReadAll(resp.Body)
	if string(rb) != "PUBLISHED" {
		t.Fail()
	}
}