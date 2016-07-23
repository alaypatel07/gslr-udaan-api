package pubsub

import (
	"testing"
	"time"
)
type consumer struct {
	d string
}

func (c consumer) Consume (data interface{}) {
	c.d = data.(string)
}

func TestNewChannel(t *testing.T) {
	//Test to create new channel
	c := NewChannel("foo")
	if (*cs.channels["foo"]).name != (*c).name {
		t.Fail()
	}
	//Test trying to create already created channel
	c1 := NewChannel("foo")
	if c != c1 {
		t.Fail()
	}
}

func TestNewPublisher(t *testing.T) {
	p := NewPublisher("foo")
	if ((*cs.channels["foo"]).name != p.c.name) {
		t.Fail()
	}
}

func TestChannel_Subscribe(t *testing.T) {
	p := NewPublisher("foo")
	var c consumer
	cs.channels["foo"].Subscribe(c)
	p.Publish("bar")
	time.Sleep(time.Second * 1)
	if "bar" != c.d {
		t.Fail()
	}
	time.Sleep(time.Second * 1)
	p.Publish("world")
	if c.d == "w" {
		t.Fail()
	}
}