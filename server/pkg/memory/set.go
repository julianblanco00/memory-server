package memory

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type Opt string

type Options map[Opt]*int64

const (
	nx      Opt = "NX"                          // Only set the key if it does not already exist.
	xx      Opt = "XX"                          // Only set the key if it already exists.
	get     Opt = "GET"                         // Return the old string stored at key, or nil if key did not exist. An error is returned and SET aborted if the value stored at key is not a string.
	keepttl Opt = "KEEPTTL"                     // Retain the time to live associated with the key.
	ex      Opt = "EX seconds"                  // Set the specified expire time, in seconds (a positive integer).
	px      Opt = "PX milliseconds"             // Set the specified expire time, in milliseconds (a positive integer).
	exat    Opt = "EXAT timestamp-seconds"      // Set the specified Unix time at which the key will expire, in seconds (a positive integer).
	pxat    Opt = "PXAT timestamp-milliseconds" // Set the specified Unix time at which the key will expire, in milliseconds (a positive integer).
)

func isExpOpt(opt Opt) bool {
	return opt == "EX" || opt == "PX" || opt == "EXAT" || opt == "PXAT"
}

func getOldValue(options Options, data *StringData, k string) (interface{}, error) {
	var oldValue interface{}
	if _, ok := options[get]; ok {
		if v, ok := data.values[k]; ok {
			if reflect.TypeOf(v).Kind() != reflect.String {
				return "", fmt.Errorf("value stored at %s is not a string", k)
			}
			oldValue = v
		}
	}
	return oldValue, nil
}

func getExpAt(options Options, data *StringData, k string) (*int64, error) {
	count := 0
	var expAt *int64

	for o, v := range options {
		if isExpOpt(o) || o == "KEEPTTL" {
			count++
			if count > 1 {
				return nil, fmt.Errorf("more than 1 expiry option specified")
			}
			if o == "KEEPTTL" {
				if old, ok := data.values[k]; ok {
					expAt = old.expIn
				}
			}
			if o == "EX" {
				d := time.Duration(*v) * time.Millisecond // seconds to millis
				n := int64(d)
				expAt = &n
			}
			if o == "PX" {
				expAt = v
			}
			if o == "EXAT" {
				diff := time.Until(time.Unix(*v, 0)).Milliseconds()
				expAt = &diff
			}
			if o == "PXAT" {
				diff := time.Until(time.Unix(*v/1000, 0)).Milliseconds()
				expAt = &diff
			}
		}
	}

	return expAt, nil
}

func checkValidWriteOption(options Options, data *StringData, k string) error {
	_, hasNx := options[nx]
	_, hasXx := options[xx]

	if hasNx && hasXx {
		return fmt.Errorf("option NX AND XX specified, can't use both")
	}

	if hasNx {
		if _, ok := data.values[k]; ok {
			return fmt.Errorf("option NX was specified but key already has a value")
		}
	}

	if hasXx {
		if _, ok := data.values[k]; !ok {
			return fmt.Errorf("option XX was specified but key does not exist yet")
		}
	}

	return nil
}

func Set(data *StringData, k, v string, opts []string) (string, error) {
	options := make(Options)

	for i := 0; i < len(opts); i++ {
		opt := opts[i]
		if isExpOpt(Opt(opt)) {
			n, err := strconv.Atoi(opts[i+1])
			if err != nil {
				return "", fmt.Errorf("invalid number parameter: %v \n", opts[i+1])
			}
			n64 := int64(n)
			options[Opt(opt)] = &n64
			i++
			continue
		}
		options[Opt(opt)] = nil
	}

	error := checkValidWriteOption(options, data, k)
	if error != nil {
		return "", error
	}

	oldValue, error := getOldValue(options, data, k)
	if error != nil {
		return "", error
	}

	expIn, err := getExpAt(options, data, k)
	if err != nil {
		return "", err
	}

	value := StringValue{
		data:  v,
		expIn: expIn,
	}

	data.mutex.Lock()
	data.values[k] = value
	data.mutex.Unlock()

	if oldValue != nil {
		return oldValue.(string), nil
	}

	return "OK", nil
}
