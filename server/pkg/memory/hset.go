package memory

import "strconv"

func hSet(d *HashData, k string, v []string) (string, error) {
	m := make(map[string]string)

	c := 0
	i := 0
	for i < len(v) {
		if _, ok := m[v[i]]; !ok {
			m[v[i]] = v[i+1]
			c++
		}
		i += 2
	}

	d.values[k] = HashValue{
		data: m,
	}

	return strconv.Itoa(c), nil
}
