package data

import (
	"testing"
)

func TestEvent_Map(t *testing.T) {
	e := Event{
		Name: "gs",
		Participants: []Participant{
			Participant{
				Name: "foo",
				Team: "gs",
			},
			Participant{
				Name: "bar",
				Team: "gs",
			},
		},
	}
	_, err := e.Map()
	if err == nil {
		t.Fail()
	}
	e.Participants[1].Team = "lr"
	m, err := e.Map()
	if err != nil || m[e.Participants[0].Team] != e.Participants[0].Name ||
	m[e.Participants[1].Team] != e.Participants[1].Name {
		t.Fail()
	}
}

func TestEvent_Register(t *testing.T) {
	e := Event{
		Name: "gs",
		Participants: []Participant{
			Participant{
				Name: "foo",
				Team: "gs",
			},
			Participant{
				Name: "bar",
				Team: "lr",
			},
		},
	}
	err := e.Register()
	if err != nil {
		t.Fail()
	}
	event, err := redisClient.Search(e.Name)
	_, ok := event.(map[string]string);
	if !ok {
		t.Fail()
	}
}
