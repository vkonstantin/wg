package memory

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTODOlist(t *testing.T) {
	conf := Config{
		GarbageCollectionInterval: time.Second,
		DeduplicationCapacitySize: 10,
		GarbageThreshold:          0,
	}
	m := newStorage(conf)
	i, err := m.AddTODO(1, "text1")
	assert.Nil(t, err)
	assert.Equal(t, "text1", i.Text)

	m.AddTODO(1, "text2")
	m.AddTODO(1, "text3")

	list, _ := m.ListTODOs()
	assert.Equal(t, 3, len(list))
	assert.Equal(t, "text3", list[2].Text)

	m.RemoveTODO(list[1].ID)
	list, _ = m.ListTODOs()
	assert.Equal(t, 2, len(list))
	assert.Equal(t, "text3", list[1].Text)

	// internal state
	assert.Equal(t, 3, len(m.todos))
	m.startGarbageCollectorOfTODOs()
	assert.Equal(t, 2, len(m.todos))

	list, _ = m.ListTODOs()
	assert.Equal(t, 2, len(list))
	assert.Equal(t, "text3", list[1].Text)
}
