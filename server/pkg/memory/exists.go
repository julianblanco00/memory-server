package memory

func sExists(sd *StringData, keys []string) (int, error) {
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

func Exists(keys []string) (int, error) {
	sc, err := sData.exists(keys)
	if err != nil {
		return 0, err
	}
	hc, err := hData.exists(keys)
	if err != nil {
		return 0, err
	}
	return sc + hc, nil
}
