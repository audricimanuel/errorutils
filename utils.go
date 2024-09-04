package errorutils

import (
	"fmt"
	"reflect"
)

// GetJsonTagInStruct to get the JSON tag of struct field
func GetJsonTagInStruct(fieldName string, structOfField any) string {
	res, err := getFieldJSONTagRecursive(reflect.ValueOf(structOfField).Type(), fieldName)
	if err != nil {
		return ""
	}
	return res
}

func getFieldJSONTagRecursive(t reflect.Type, fieldName string) (string, error) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Check if the field is the one we're looking for
		if field.Name == fieldName {
			return field.Tag.Get("json"), nil
		}

		// If the field is a struct, search recursively
		if field.Type.Kind() == reflect.Struct {
			if tag, err := getFieldJSONTagRecursive(field.Type, fieldName); err == nil {
				return tag, nil
			}
		}
	}
	return "", fmt.Errorf("field %s not found", fieldName)
}
