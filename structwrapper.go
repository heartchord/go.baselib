package goblazer

import "reflect"

// GetStructFieldNum :
func GetStructFieldNum(s interface{}) uint32 {
	fieldNum := uint32(0)

	o := reflect.ValueOf(s)

	if o.Kind() != reflect.Struct {
		panic("Unexpected Data Type")
	}

	for i := 0; i < o.NumField(); i++ {
		if o.Field(i).Kind() == reflect.Struct {
			fieldNum += GetStructFieldNum(o.Field(i).Interface())
			continue
		}

		fieldNum++
	}

	return fieldNum
}

// GetStructFieldNames :
func GetStructFieldNames(s interface{}) []string {
	var ret []string

	t := reflect.TypeOf(s)
	o := reflect.ValueOf(s)

	if o.Kind() != reflect.Struct {
		panic("Unexpected Data Type")
	}

	for i := 0; i < o.NumField(); i++ {
		if o.Field(i).Kind() == reflect.Struct {
			sub := GetStructFieldNames(o.Field(i).Interface())
			ret = append(ret, sub...)
			continue
		}

		ret = append(ret, t.Field(i).Name)
	}

	return ret
}
