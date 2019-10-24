package model

import (
	"encoding/json"
	"github.com/vkonstantin/wg/todo/message"
)

// Item is internal struct of todo item
type Item struct {
	ID     uint64
	Text   string
	UserID uint64
}

// Items is collection of internal todo item
type Items []*Item

// MarshalJSON implementation of json Marshal interface
func (i Item) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Rest())
}

func (i Item) Rest() message.Item {
	m := message.Item{
		ID:   i.ID,
		Text: i.Text,
	}
	return m
}

// MarshalJSON implementation of json Marshal interface
func (is Items) MarshalJSON() ([]byte, error) {
	list := make([]message.Item, 0, len(is))
	for i := range is {
		if is[i] != nil {
			list = append(list, is[i].Rest())
		}
	}
	return json.Marshal(list)
}
