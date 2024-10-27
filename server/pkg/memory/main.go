package memory

import (
	"errors"
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

func (d *Data) multiSet(vals []string) (string, error) {
	return MultiSet(d, vals)
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

func parseRESPString(input string) ([]string, error) {
	var result []string

	i := 0
	for i < len(input) {
		if input[i] != '$' {
			return nil, errors.New("length definition should start with a $")
		}
		i++
		lengthStart := i
		for i < len(input) && input[i] != ' ' && input[i] != '\n' {
			i++
		}
		lengthStr := input[lengthStart:i]

		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return nil, fmt.Errorf("error getting length: %v", err)
		}

		i++

		if i+length > len(input) {
			return nil, errors.New("provided length exceeds input length")
		}

		item := input[i : i+length]
		result = append(result, item)
		i += length

		for i < len(input) && (input[i] == ' ' || input[i] == '\n') {
			i++
		}
	}

	return result, nil
}

func parseCommand(command string, data *Data) (interface{}, error) {
	parsed, err := parseRESPString(command)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	cmd := parsed[0]

	switch strings.TrimSpace(cmd) {
	case "GET":
		return data.getValue(parsed[1])
	case "SET":
		return data.setValue(parsed[1], parsed[2], parsed[3:])
	case "DEL":
		return data.delValue(parsed[1:])
	// case "MSET":
	// 	return data.multiSet(parts)
	default:
		return "", fmt.Errorf("invalid command %s \n", cmd)
	}
}
