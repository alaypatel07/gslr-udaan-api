package data

import (
	"errors"
)

var redisClient *RedisDatabaseManager

type Participant struct {
	Team string `json:"team"`
	Name string `json:"name"`
}

type Event struct {
	Name string `json:"name"`
	Participants []Participant `json:"participants"`
}

func (e *Event)Map() (map[string]string, error) {
	m := make(map[string]string)
	for _, team := range e.Participants {
		if m[team.Team] != "" {
			return nil, errors.New("Same team event twice")
		}
		m[team.Team] = team.Name
	}
	return m, nil
}

func init() {
	redisClient = NewRedisDatabaseManager()
}

// Register Event to the Database
func (e *Event)Register() error {
	m, err := e.Map();
	if err != nil {
		return err
	}
	if _, err := redisClient.Insert(e.Name, m); err != nil {
		return err
	}
	return nil
}