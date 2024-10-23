package memory

import (
	"fmt"
	"strconv"
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

func (d *Data) getValue(k string) (string, error) {
	v, ok := d.values_map[k]
	if !ok {
		return "", fmt.Errorf("value not found for key %s", k)
	}
	return v.data, nil
}

func (d *Data) delValue(keys []string) (string, error) {
	delCount := 0
	for _, k := range keys {
		if _, ok := d.values_map[k]; ok {
			delCount++
		}
		delete(d.values_map, k)
	}
	return strconv.Itoa(delCount), nil
}

func NewData() *Data {
	return &Data{
		amount:     0,
		values_map: make(map[string]Value),
	}
}

func parseCommand(command string, data *Data) (string, error) {
	parts := strings.Split(command, " ")
	cmd := parts[0]
	key := parts[1]

	switch strings.TrimSpace(cmd) {
	case "GET":
		return data.getValue(key)
	case "SET":
		return data.setValue(key, parts[2], parts[3:])
	case "DEL":
		keys := parts[1:]
		return data.delValue(keys)
	default:
		return "", fmt.Errorf("invalid command %s \n", cmd)
	}
}
