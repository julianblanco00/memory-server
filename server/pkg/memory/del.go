package memory

import "strconv"

func Del(keys []string, d *StringData) (string, error) {
	delCount := 0
	for _, k := range keys {
		if _, ok := d.values[k]; ok {
			delCount++
		}
		delete(d.values, k)
	}
	return strconv.Itoa(delCount), nil
}
