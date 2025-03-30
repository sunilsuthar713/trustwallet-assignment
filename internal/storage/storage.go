package storage

type Storage struct {
	data map[string]interface{}
}

func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]interface{}),
	}
}

func (s *Storage) Save(key string, value interface{}) {
	s.data[key] = value
}

func (s *Storage) Get(key string) interface{} {
	return s.data[key]
}