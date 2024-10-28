package memory

func MSet(d *StringData, vals []string) (string, error) {
	i := 0

	for i < len(vals) {
		k := vals[i]
		i++
		v := vals[i]

		d.values[k] = StringValue{
			data: v,
		}

		i++
	}

	return "OK", nil
}
