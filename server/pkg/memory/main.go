package memory

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func parseRESPString(input string) ([]string, error) {
	var result []string

	i := 0
	for i < len(input) {
		if input[i] != '$' {
			return nil, errors.New("length definition should start with a $")
		}
		i++
		lenStart := i
		for i < len(input) && input[i] != ' ' && input[i] != '\n' {
			i++
		}
		lenStr := input[lenStart:i] // -> i is the end of the length line

		length, err := strconv.Atoi(lenStr)
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

func parseCommand(command string) (interface{}, error) {
	cmd, err := parseRESPString(command)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("wrong number of arguments")
	}

	switch strings.TrimSpace(cmd[0]) {
	case "GET":
		return sData.get(cmd[1])
	case "SET":
		return sData.set(cmd[1], cmd[2], cmd[3:])
	case "DEL":
		return sData.del(cmd[1:])
	case "MSET":
		return sData.mset(cmd[1:])
	case "HSET":
		return hData.hset(cmd[1], cmd[2:])
	case "HGET":
		return hData.hget(cmd[1], cmd[2])
	case "HDEL":
		return hData.hdel(cmd[1], cmd[2:])
	case "HGETALL":
		return hData.hgetall(cmd[1])
	default:
		return "", fmt.Errorf("invalid command %s \n", cmd)
	}
}
