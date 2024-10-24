package memory

import (
	"fmt"
	"strings"
)

type Value struct {
	expIn *int64
	data  string
}

type Data struct {
	values_map map[string]Value
	amount     int32
}

func (d *Data) setValue(k, v string, opts []string) (string, error) {
	return Set(d, k, v, opts)
}

func (d *Data) getValue(k string) (interface{}, error) {
	return Get(k, d)
}

func (d *Data) delValue(keys []string) (string, error) {
	return Del(keys, d)
}

func NewData() *Data {
	return &Data{
		amount:     0,
		values_map: make(map[string]Value),
	}
}

func parseCommand(command string, data *Data) (interface{}, error) {
	parts := strings.Split(command, " ")
	cmd := parts[0]
	key := parts[1]

	switch strings.TrimSpace(cmd) {
	case "GET":
		return data.getValue(key)
	case "SET":
		return data.setValue(key, parts[2], parts[3:])
	case "DEL":
		return data.delValue(parts[1:])
	default:
		return "", fmt.Errorf("invalid command %s \n", cmd)
	}
}
