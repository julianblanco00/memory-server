package memory

import (
	"fmt"
	"reflect"
)

func Get(k string, d *StringData) (interface{}, error) {
	d.mutex.RLock()
	v, ok := d.values[k]
	d.mutex.RUnlock()
	if !ok {
		return nil, nil
	}

	if reflect.TypeOf(v.data).Kind() != reflect.String {
		return nil, fmt.Errorf("value in key is not a string")
	}

	return v.data, nil
}
