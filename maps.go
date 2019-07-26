package maps

import (
	"encoding/json"
	"strconv"
	"sync"
	"time"
)

// Map - тип данных, хеш-таблица со строковым ключем и интерфейсным значением
type Map struct {
	data map[string]interface{}
	*sync.RWMutex
}

// New - функция возвращает указатель на данные Map
func New() *Map {
	return &Map{
		data:    make(map[string]interface{}),
		RWMutex: &sync.RWMutex{},
	}
}

// Copy - функция создает новый указатель на данные Map и наполняет их структуру данными data
func Copy(data map[string]interface{}) *Map {
	m := New()
	m.AddMap(data)
	return m
}

// Set - функция устанавливает значение value с именем key
func (m *Map) Set(key string, value interface{}) {
	m.Lock()
	m.data[key] = value
	m.Unlock()
}

// AddMap - функция добавляет или затирает (в случае наличия) данные
func (m *Map) AddMap(data map[string]interface{}) {
	if data == nil {
		return
	}
	for k, v := range data {
		m.Lock()
		m.data[k] = v
		m.Unlock()
	}
}

// Get - фукнция получает интерфейсное значение из набора данных
func (m *Map) Get(key string) interface{} {
	m.RLock()
	value, ok := m.data[key]
	m.RUnlock()
	if !ok {
		return nil
	}
	return value
}

// Del - фукнция удаляет из набора данные с клюечом key
func (m *Map) Del(key string) {
	m.Lock()
	delete(m.data, key)
	m.Unlock()
}

// GetMap - функция возвращает данные подчиненного ключа, как указатель на Map
func (m *Map) GetMap(key string) *Map {
	m.RLock()
	value, ok := m.data[key]
	m.RUnlock()

	if !ok {
		return New()
	}
	mapValue, ok := value.(map[string]interface{})
	if !ok {
		return New()
	}

	result := Copy(mapValue)
	return result
}

func (m *Map) GetGoMap() map[string]interface{} {
	result := make(map[string]interface{})
	m.RLock()
	for k, v := range m.data {
		result[k] = v
	}
	m.RUnlock()
	return result
}

// MarshalJSON - функция для кастомного маршалинга в JSON
func (m *Map) MarshalJSON() (data []byte, err error) {
	m.RLock()
	defer m.RUnlock()
	return json.Marshal(m.data)
}

// UnmarshalJSON - функция для кастомного размаршалинга набора байт
func (m *Map) UnmarshalJSON(data []byte) (err error) {
	m.Lock()
	defer m.Unlock()
	return json.Unmarshal(data, &m.data)
}

// GetFloat64 - возвращает значение по ключу, как float 64-битное
func (m *Map) GetFloat64(key string) float64 {
	value := m.Get(key)
	if value == nil {
		return 0
	}
	switch v := value.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case string:
		i, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0
		}
		return i
	default:
		return 0
	}
}

// GetInt64 - возвращает значение по ключу, как целое 64-битное число
func (m *Map) GetInt64(key string) int64 {
	value := m.Get(key)
	if value == nil {
		return 0
	}
	switch v := value.(type) {
	case int:
		return int64(v)
	case int64:
		return v
	case float64:
		return int64(v)
	case string:
		i, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			return i
		}
		bv, err := strconv.ParseBool(v)
		if err == nil {
			if bv {
				return 1
			}
		}
		return 0
	case bool:
		if v {
			return 1
		}
		return 0
	default:
		return 0
	}
}

// GetInt - возвращает значение по ключу, как целое число
func (m *Map) GetInt(key string) int {
	value := m.Get(key)
	if value == nil {
		return 0
	}
	switch v := value.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case string:
		i, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			return int(i)
		}
		bv, err := strconv.ParseBool(v)
		if err == nil {
			if bv {
				return 1
			}
		}
		return 0
	case bool:
		if v {
			return 1
		}
		return 0
	default:
		return 0
	}
}

// GetString - функция возвращает значение по ключу, как строку
func (m *Map) GetString(key string) string {
	value := m.Get(key)
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case string:
		return v
	case int64:
		return strconv.FormatInt(v, 64)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		return strconv.Itoa(v)
	case time.Time:
		return v.Format(time.RFC3339)
	case time.Duration:
		return v.String()
	case error:
		return v.Error()
	default:
		return ""
	}
}

// GetTime - функция возвращает значение по ключу key, как время
func (m *Map) GetTime(key string) time.Time {
	value := m.Get(key)
	if value == nil {
		return time.Unix(0, 0)
	}

	switch v := value.(type) {
	case time.Time:
		return v
	case string:
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return time.Unix(0, 0)
		}
		return t
	case int64:
		return time.Unix(v, 0)
	case float64:
		return time.Unix(int64(v), 0)
	case int:
		return time.Unix(int64(v), 0)
	default:
		return time.Unix(0, 0)
	}
}

// GetJSTime - функция вернет целое число по ключу key, попытаясь превратить его во время, сперва
func (m *Map) GetJSTime(key string) int64 {
	t := m.GetTime(key)
	return t.Unix() * 1000
}

// SetJSTime - функция устанавливает время из целого числа, принятого в JS
func (m *Map) SetJSTime(key string, value int64) {
	t := time.Unix(value/1000, 0)
	m.Set(key, t)
}

// GetBool - возвращает значение по ключу, как булевое значение
func (m *Map) GetBool(key string) bool {
	value := m.Get(key)
	if value == nil {
		return false
	}
	switch v := value.(type) {
	case bool:
		return v
	case int:
		if v == 0 {
			return false
		}
		return true
	case int64:
		if v == 0 {
			return false
		}
		return true
	case float64:
		if v == 0 {
			return false
		}
		return true
	case string:
		b, err := strconv.ParseBool(v)
		if err != nil {
			return false
		}
		return b
	default:
		return false
	}
}

// Length - возвращает размерность данных
func (m *Map) Length() int {
	m.RLock()
	defer m.RUnlock()
	return len(m.data)
}
