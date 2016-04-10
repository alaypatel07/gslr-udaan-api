package pubsub

import "testing"

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

func TestChannel_Subscribe(t *testing.T) {
	ch := NewChannel("foo")
	c := NewClient()
	ch.Subscribe(c)
	if len(ch.Clients) != 1 {
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
	ch.Publish("Working")
	c2Data := <- c2.C["foo"]
	c1Data := <- c1.C["foo"]
	c1StringData, _ := c1Data.(string)
	c2StringData, _ := c2Data.(string)
	if c1StringData != "working" || c2StringData != "working" {
		t.Fail()
	}
}