package memory

func Exists(sd *StringData, keys []string) (int, error) {
	count := 0

	for _, k := range keys {
		if _, ok := sd.values[k]; ok {
			count++
		}
	}

	return count, nil
}

func HExists(hd *HashData, keys []string) (int, error) {
	count := 0

	for _, k := range keys {
		if _, ok := hd.values[k]; ok {
			count++
		}
	}

	return count, nil
}
