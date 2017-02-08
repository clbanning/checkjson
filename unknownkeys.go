// validate.go - check JSON object against struct definition
// Copyright © 2016-2017 Charles Banning.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package checkjson

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
)

// UnknownJSONKeys() returns a list of the JSON object keys that will not
// be decoded to a member of 'val', which is of type struct.  For nested
// JSON objects the keys are reported using dot-notation.
func UnknownJSONKeys(b []byte, val interface{}) ([]string, error) {
	s := make([]string, 0)
	m := make(map[string]interface{})
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	if err := checkAllFields(m, reflect.ValueOf(val), &s, ""); err != nil {
		return s, err
	}
	return s, nil
}

func checkAllFields(mv interface{}, val reflect.Value, s *[]string, key string) error {
	// get the name when initially called
	if key == "" {
		key = val.Type().Name()
	}

	// 1. Convert any pointer value.
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val) // convert ptr to struct
	}
	// zero Value?
	if !val.IsValid() {
		return nil
	}
	typ := val.Type()

	// json.RawMessage is a []byte/[]uint8 and has Kind() == reflect.Slice
	if typ.Name() == "RawMessage" {
		return nil
	}

	// 2. If its a slice then 'mv' should hold a []interface{} value.
	//    Loop through the members of 'mv' and see that they are valid relative
	//    to the <T> of val []<T>.
	if typ.Kind() == reflect.Slice {
		tval := typ.Elem()
		if tval.Kind() == reflect.Ptr {
			tval = tval.Elem()
		}
		// slice may be nil, so create a Value of it's type
		// 'mv' must be of type []interface{}
		sval := reflect.New(tval)
		slice, ok := mv.([]interface{})
		if !ok {
			// return fmt.Errorf("JSON value not an array")
			*s = append(*s, key)
			return nil
		}
		// 2.1. Check members of JSON array.
		//      This forces all of them to be regular and w/o typos in key labels.
		for n, sl := range slice {
			_ = checkAllFields(sl, sval, s, key+"."+strconv.Itoa(n+1))
			// if err := checkAllFields(sl, sval, s, key+strconv.Itoa(n+1)); err != nil {
			// 	return fmt.Errorf("[array element #%d] %s", n+1, err.Error())
			// }
		}
		return nil // done with reflect.Slice value
	}

	// 3a. Ignore anything that's not a struct.
	if typ.Kind() != reflect.Struct {
		return nil // just ignore it - don't look for k:v pairs
	}
	// 3b. map value must represent k:v pairs
	mm, ok := mv.(map[string]interface{})
	if !ok {
		*s = append(*s, key)
	}
	// if !ok {
	// 	return fmt.Errorf("JSON object does not have k:v pairs for member: %s",
	// 		typ.Name)
	// }

	// 4. Build the map of struct field name:value
	//    We make every key (field) label look like an exported label - "Fieldname".
	//    If there is a JSON tag it is used instead of the field label, and saved to
	//    insure that the spec'd tag matches the JSON key exactly.
	type fieldSpec struct {
		val reflect.Value
		tag string
	}
	fieldCnt := val.NumField()
	fields := make(map[string]*fieldSpec, fieldCnt)
	for i := 0; i < fieldCnt; i++ {
		if len(typ.Field(i).PkgPath) > 0 {
			continue // field is NOT exported
		}
		tag := typ.Field(i).Tag.Get("json")
		if tag == "" {
			fields[strings.Title(strings.ToLower(typ.Field(i).Name))] = &fieldSpec{val.Field(i), ""}
		} else {
			fields[strings.Title(strings.ToLower(tag))] = &fieldSpec{val.Field(i), tag}
		}
	}

	// 5. check that map keys correspond to exported field names
	var spec *fieldSpec
	for k, m := range mm {
		lk := strings.ToLower(k)
		for _, sk := range skipkeys {
			if lk == sk {
				goto next
			}
		}
		spec, ok = fields[strings.Title(lk)]
		if !ok {
			*s = append(*s, key+"."+lk)
			return nil
			// return fmt.Errorf("no member for JSON key: %s", k)
		}
		if len(spec.tag) > 0 && spec.tag != k { // JSON key doesn't match Field tag
			// return fmt.Errorf("key: %s -  does not match tag: %s", k, spec.tag)
			*s = append(*s, key+"["+spec.tag+"]") // include tag in brackets
			return nil
		}
		_ = checkAllFields(m, spec.val, s, key+"."+lk)
		// if err := checkFields(m, spec.val); err != nil { // could be nested structs
		// 	return fmt.Errorf("checking subkeys of JSON key: %s - %s", k, err.Error())
		// }
	next:
	}

	return nil
}
