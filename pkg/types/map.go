package types

type Map[K comparable, V any] map[K]V

func (m Map[K, V]) Get(key K, fallback V) V {
	if v, ok := m[key]; ok {
		return v
	}
	return fallback
}
