// validate.go - check JSON object against struct definition
// Copyright © 2016 Charles Banning.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package checkjson provides functions for checking JSON object keys against
// struct member names/tags to see if they will decode using encoding/json package.
//
// There are several options: Validate returns an error on the first key:value pair
// that won't decode; UnknownJSONKeys returns a slice of all the keys that won't decode
// using dot-notation for nested JSON object keys.  A complementary function
// MissingJSONKeys provides a slice of struct members that won't be set by the JSON 
// object using dot-notation for nested struct members.  (See test cases for examples.)
/*
EXAMPLE

	import "github.com/clbanning/checkjson"
	...
	type Home struct {
		Addr string
		Port int
	}

	type config struct {
		Id       string
		Address  *Home `json:"addr"`
		Failover []*Home
		Log      bool
	}

	configData := []byte(`
		{
			"id":"example",
			"addr":{
				"addr":"127.0.0.1",
				"port":12345
			},
			"failover":[{
					"addr":"127.0.0.1",
					"port":12346
				},{
					"add":"24.33.6.1",
					"port":80
				}],
			"log":true
		}`)

	c := new(config)
	if err := checkjson.Validate(configData, c); err != nil {
		// handle error:
		// checking subkeys of JSON key: failover - [array element #2] no member for JSON key: add
	}
*/
package checkjson

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// List of JSON keys to NOT validate.
var skipkeys = []string{"config"}

// SetKeysToIgnore maintains a list of JSON keys that should
// not be validated as exported struct fields.  By default the
// JSON key "config" is not validated; it can be removed from
// the list by calling SetKeysToIgnore() with no arguments.
// The arguments are used as the list of keys to ignore and
// override the default. NOTE: keys are case insensitive - i.e.,
// "key" == "Key" == "KEY".
func SetKeysToIgnore(s ...string) {
	if len(s) == 0 {
		skipkeys = skipkeys[:0] // remove "config"
		return
	}
	skipkeys = make([]string, len(s))
	for i, v := range s {
		skipkeys[i] = strings.ToLower(v)
	}
}

// Validate scans a JSON object and returns an error when it encounters
// a key:value pair that will not decode to a member of the 'val' 
// of type struct using the "encoding/json" package. 
func Validate(b []byte, val interface{}) error {
	m := make(map[string]interface{})
	if err := json.Unmarshal(b, &m); err != nil {
		return ResolveJSONError(b, err)
	}
	if err := checkFields(m, reflect.ValueOf(val)); err != nil {
		return err
	}
	return nil
}

func checkFields(mv interface{}, val reflect.Value) error {
	// 1. Convert any pointer value.
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val) // convert ptr to struc
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
			return fmt.Errorf("JSON value not an array")
		}
		// 2.1. Check members of JSON array.
		//      This forces all of them to be regular and w/o typos in key labels.
		for n, sl := range slice {
			if err := checkFields(sl, sval); err != nil {
				return fmt.Errorf("[array element #%d] %s", n+1, err.Error())
			}
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
		return fmt.Errorf("JSON object does not have k:v pairs for member: %s",
			typ.Name())
	}

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
			return fmt.Errorf("no member for JSON key: %s", k)
		}
		if len(spec.tag) > 0 && spec.tag != k { // JSON key doesn't match Field tag
			return fmt.Errorf("key: %s -  does not match tag: %s", k, spec.tag)
		}
		if err := checkFields(m, spec.val); err != nil { // could be nested structs
			return fmt.Errorf("checking subkeys of JSON key: %s - %s", k, err.Error())
		}
	next:
	}

	return nil
}

// ResolveJSONError tries to augment json.Unmarshal errors with
// the JSON context - key:value - if possible.  (This is useful
// when errors occur when unmarshaling large JSON objects.)
func ResolveJSONError(data []byte, err error) error {
	if e, ok := err.(*json.UnmarshalTypeError); ok {
		// grab stuff ahead of the error
		var i int
		var getKey bool
		for i = int(e.Offset) - 1; i != -1; i-- {
			switch data[i] {
			case ':':
				getKey = true
			case '\n', '{', '[', ',', ' ':
				if getKey {
					goto done
				}
			}
		}
	done:
		info := strings.TrimSpace(string(data[i+1 : int(e.Offset)]))
		return fmt.Errorf("%s - at: %s", e.Error(), info)
	}
	// just report all other unmarshal errors
	return err
}

