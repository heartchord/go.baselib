package goblazer

import "reflect"

// GetStructFieldNum :
func GetStructFieldNum(s interface{}) uint32 {
	fieldNum := uint32(0)

	o := reflect.ValueOf(s)

	if o.Kind() != reflect.Struct {
		panic("Unexpected Data Type")
	}

	count := o.NumField()
	for i := 0; i < count; i++ {
		v := o.Field(i)
		if v.Kind() == reflect.Struct {
			fieldNum += GetStructFieldNum(v.Interface())
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

	count := o.NumField()
	for i := 0; i < count; i++ {
		v := o.Field(i)

		if v.Kind() == reflect.Struct {
			sub := GetStructFieldNames(v.Interface())
			ret = append(ret, sub...)
			continue
		}

		ret = append(ret, t.Field(i).Name)
	}

	return ret
}

// GetStructFieldTags :
func GetStructFieldTags(s interface{}) []string {
	var ret []string

	t := reflect.TypeOf(s)
	o := reflect.ValueOf(s)

	if o.Kind() != reflect.Struct {
		panic("Unexpected Data Type")
	}

	count := o.NumField()
	for i := 0; i < count; i++ {
		v := o.Field(i)

		if v.Kind() == reflect.Struct {
			sub := GetStructFieldTags(v.Interface())
			ret = append(ret, sub...)
			continue
		}

		ret = append(ret, string(t.Field(i).Tag))
	}

	return ret
}
