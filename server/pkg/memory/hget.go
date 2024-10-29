package memory

func hGet(h *HashData, k, f string) (interface{}, error) {
	if m, ok := h.values[k]; ok {
		if v, ok := m.data[f]; ok {
			return v, nil
		}
	}

	return nil, nil
}
