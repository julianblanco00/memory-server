package memory

import "encoding/json"

func hGetAll(h *HashData, k string) (interface{}, error) {
	v, ok := h.values[k]
	if !ok {
		return "empty hash", nil
	}

	json, err := json.Marshal(v.data)
	if err != nil {
		return "", err
	}

	return string(json), nil
}
