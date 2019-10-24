package memory

import (
	"github.com/vkonstantin/wg/todo/common/auth"
	"github.com/vkonstantin/wg/todo/model"
	"github.com/vkonstantin/wg/todo/storage"
	"sync"
	"sync/atomic"
	"time"
)

// Config is a storage configuration
type Config struct {
	GarbageCollectionInterval time.Duration
	DeduplicationCapacitySize int
	GarbageThreshold          float64
}

// NewDefault create memory storage with default config
func NewDefault() storage.Storage {
	conf := Config{
		GarbageCollectionInterval: time.Minute,
		DeduplicationCapacitySize: 100,
		GarbageThreshold:          0.5,
	}

	return New(conf)
}

type mem struct {
	config Config
	// id for addUser
	userID uint64
	// RequestIDs
	requestIDs          map[string]int64 // map[key]time
	lastRequestIDs      []int64
	lastRequestIDsIndex int
	requestIDsLock      *sync.RWMutex
	// todos
	todoID       uint64
	todos        []*model.Item
	todosLock    *sync.RWMutex
	todosRemoved int
	todosIndex   map[uint64]int // index for fast search. map[todoID]int - index in todos array
}

// New create memory storage
func New(config Config) storage.Storage {
	s := newStorage(config)
	go s.startGarbageCollectorLoop()
	return s
}

func newStorage(config Config) *mem {
	if config.DeduplicationCapacitySize <= 0 {
		config.DeduplicationCapacitySize = 1
	}
	m := mem{
		config:         config,
		requestIDs:     make(map[string]int64, config.DeduplicationCapacitySize),
		lastRequestIDs: make([]int64, config.DeduplicationCapacitySize),
		requestIDsLock: &sync.RWMutex{},
		todos:          make([]*model.Item, 0),
		todosLock:      &sync.RWMutex{},
		todosIndex:     make(map[uint64]int, 0),
	}

	return &m
}

func (m *mem) ListTODOs() ([]*model.Item, error) {
	m.todosLock.RLock()
	defer m.todosLock.RUnlock()

	list := make([]*model.Item, 0, len(m.todos)-m.todosRemoved)
	for i := range m.todos {
		if m.todos[i] != nil {
			item := new(model.Item)
			*item = *m.todos[i] // make a first-level copy. If the structure has pointers that it is needed to do a deep copy
			list = append(list, item)
		}
	}
	return list, nil
}

func (m *mem) AddTODO(userID uint64, text string) (*model.Item, error) {
	m.todosLock.Lock()
	defer m.todosLock.Unlock()

	m.todoID++
	item := model.Item{
		ID:     m.todoID,
		Text:   text,
		UserID: userID,
	}
	m.todos = append(m.todos, &item)
	m.todosIndex[m.todoID] = len(m.todos) - 1

	return &item, nil
}

func (m *mem) RemoveTODO(itemID uint64) error {
	m.todosLock.Lock()
	defer m.todosLock.Unlock()

	if arrayPosition, ok := m.todosIndex[itemID]; ok {
		m.todos[arrayPosition] = nil
		delete(m.todosIndex, itemID)
		m.todosRemoved++
	}

	return nil
}

func (m *mem) AddUser() (*auth.User, error) {
	newUserID := atomic.AddUint64(&m.userID, 1)
	user := auth.User{ID: newUserID}
	return &user, nil
}

func (m *mem) IsDuplicate(requestID string) bool {
	// Read lock just for check
	m.requestIDsLock.RLock()
	if _, ok := m.requestIDs[requestID]; ok {
		m.requestIDsLock.RUnlock()
		return true
	}
	m.requestIDsLock.RUnlock()

	// Need to add item. Then we need write lock and double check
	m.requestIDsLock.Lock()
	defer m.requestIDsLock.Unlock()

	if _, ok := m.requestIDs[requestID]; ok {
		return true
	}

	ts := time.Now().UnixNano()
	m.requestIDs[requestID] = ts
	m.lastRequestIDs[m.lastRequestIDsIndex] = ts
	m.lastRequestIDsIndex++
	if m.lastRequestIDsIndex == m.config.DeduplicationCapacitySize {
		m.lastRequestIDsIndex = 0
	}

	return false
}

func (m *mem) startGarbageCollectorLoop() {
	for {
		time.Sleep(m.config.GarbageCollectionInterval)
		m.startGarbageCollector()
	}
}

func (m *mem) startGarbageCollector() {
	m.startGarbageCollectorOfRequestIDs()
	m.startGarbageCollectorOfTODOs()
}

func (m *mem) startGarbageCollectorOfTODOs() {
	if len(m.todos) > 0 && (float64(m.todosRemoved)/float64(len(m.todos))) > m.config.GarbageThreshold {
		m.todosLock.Lock()
		defer m.todosLock.Unlock()

		list := make([]*model.Item, 0, len(m.todos)-m.todosRemoved)
		m.todosIndex = make(map[uint64]int, len(m.todos)-m.todosRemoved)
		for _, i := range m.todos {
			if i != nil {
				list = append(list, i)
				m.todosIndex[i.ID] = len(list) - 1
			}
		}
		m.todos = list
		m.todosRemoved = 0
	}
}

func (m *mem) startGarbageCollectorOfRequestIDs() {
	if int(float64(len(m.requestIDs))*m.config.GarbageThreshold) > m.config.DeduplicationCapacitySize {
		m.requestIDsLock.Lock()
		defer m.requestIDsLock.Unlock()

		minTs := m.lastRequestIDs[m.lastRequestIDsIndex]
		for requestID, ts := range m.requestIDs {
			if ts < minTs {
				delete(m.requestIDs, requestID)
			}
		}
	}
}
