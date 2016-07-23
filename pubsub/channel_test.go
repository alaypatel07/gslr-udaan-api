package pubsub

import (
	"testing"
)

func TestNewChannel(t *testing.T) {
	c := NewChannel("foo")
	if c.Name != "foo" && len(c.Clients) != 0{
		t.Fail()
	}
}

func TestNewClient(t *testing.T) {
	c := NewClient()
	if c.C == nil {
		t.Fail()

	}
}

func TestClient_Subscribe(t *testing.T) {
	ch := NewChannel("foo")
	c := NewClient()
	c.Subscribe(ch)
	if len(ch.Clients) != 1 {
		t.Fail()
	}
}

func TestChannel_Publish(t *testing.T) {
	ch := NewChannel("foo")
	c1 := NewClient()
	c2 := NewClient()
	ch.Subscribe(c1)
	ch.Subscribe(c2)
	go func() {
		c2Data := <- c2.C["foo"]
		c2StringData, _ := c2Data.(string)
		if c2StringData != "working" {
		t.Fail()
	}
	}()
	go func() {
		c1Data := <- c1.C["foo"]
		c1StringData, _ := c1Data.(string)
		if c1StringData != "working" {
			t.Fail()
		}
	}()
	ch.Publish("working")

}