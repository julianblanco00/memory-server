package memory

type StringValue struct {
	expIn *int64
	data  string
}
type StringData struct {
	values map[string]StringValue
	amount int32
}

func (d *StringData) exists(keys []string) (int, error) {
	return Exists(d, keys)
}

func (d *StringData) append(keys []string) (int, error) {
	return Append(d, keys)
}

func (d *StringData) mset(vals []string) (string, error) {
	return MSet(d, vals)
}

func (d *StringData) set(k, v string, opts []string) (string, error) {
	return Set(d, k, v, opts)
}

func (d *StringData) get(k string) (interface{}, error) {
	return Get(k, d)
}

func (d *StringData) del(keys []string) (string, error) {
	return Del(keys, d)
}

func NewStringData() *StringData {
	return &StringData{
		amount: 0,
		values: make(map[string]StringValue),
	}
}
