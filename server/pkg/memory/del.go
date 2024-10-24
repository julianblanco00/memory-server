package memory

import "strconv"

func Del(keys []string, d *Data) (string, error) {
	delCount := 0
	for _, k := range keys {
		if _, ok := d.values_map[k]; ok {
			delCount++
		}
		delete(d.values_map, k)
	}
	return strconv.Itoa(delCount), nil
}
