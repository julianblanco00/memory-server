package memory

import "fmt"

func Append(sd *StringData, key, value string) (int, error) {
	v, err := sd.get(key)
	if err != nil {
		return 0, err
	}

	if v == 0 {
		_, err := sd.set(key, value, []string{})
		if err != nil {
			return 0, err
		}
		return len(value), nil
	}

	current := v.(string)
	newVal := fmt.Sprintf("%s%s", current, value)

	_, err = sd.set(key, newVal, []string{})
	if err != nil {
		return 0, err
	}

	return len(newVal), nil
}
