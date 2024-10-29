package memory

import (
	"strconv"
)

func hDel(h *HashData, k string, fs []string) (string, error) {
	c := 0
	hv, ok := h.values[k]

	if !ok {
		return "0", nil
	}

	for _, f := range fs {
		if _, ok := hv.data[f]; ok {
			delete(hv.data, f)
			c++
		}
	}

	if len(hv.data) == 0 {
		delete(h.values, k)
	}

	return strconv.Itoa(c), nil
}
