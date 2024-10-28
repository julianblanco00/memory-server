package memory

func HSet(d *HashData, k string, v []string) (string, error) {
	m := make(map[string]string)

	i := 0
	for i < len(v) {
		m[v[i]] = v[i+1]
		i += 2
	}

	d.values[k] = HashValue{
		data: m,
	}

	return "", nil
}
