package memory

type HashValue struct {
	expIn *int64
	data  map[string]string
}
type HashData struct {
	values map[string]HashValue
	amount int32
}

func (h *HashData) hset(k string, vals []string) (string, error) {
	return hSet(h, k, vals)
}

func (h *HashData) hget(k, f string) (interface{}, error) {
	return hGet(h, k, f)
}

func (h *HashData) hgetall(k string) (interface{}, error) {
	return hGetAll(h, k)
}

func (h *HashData) hdel(k string, fs []string) (string, error) {
	return hDel(h, k, fs)
}

func NewHashData() *HashData {
	return &HashData{
		amount: 0,
		values: make(map[string]HashValue),
	}
}
