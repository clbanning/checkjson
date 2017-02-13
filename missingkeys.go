// validate.go - check JSON object against struct definition
// Copyright © 2016-2017 Charles Banning.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// ==== from: https://forum.golangbridge.org/t/how-to-detect-missing-bool-field-in-json-post-data/3861

package checkjson

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type skipmems struct {
	val   string
	depth int
}

// Should we ignore "omitempty" struct tags. By default accept tag.
var omitemptyOK = true

// IgnoreOmitemptyTag determines whether a `json:",omitempty"` tag is recognized or
// not with respect to the JSON object.  By default MissingJSONKeys will not include
// struct members that are tagged with "omitempty" in the list of missing JSON keys.
// If the function is toggled or passed the optional argument 'false' then missing
// JSON keys may include those for struct members with the 'omitempty' JSON tag.
//
// Calling IgnoreOmitemptyTag with no arguments toggles the handling on/off.  If
// the alternative argument is passed, then the argument value determines the
// "omitempty" handling behavior.
func IgnoreOmitemptyTag(ok ...bool) {
	if len(ok) == 0 {
		omitemptyOK = !omitemptyOK
		return
	}
	omitemptyOK = ok[0]
}

// Slice of dot-notation struct members that can be missing in JSON object.
var skipmembers = []skipmems{}

// SetMembersToIgnore creates a list of exported struct member names that should not be checked
// for as keys in the JSON object.  For hierarchical struct members provide the full path for
// the member name using dot-notation. Calling SetMembersToIgnore with no arguments -
// SetMembersToIgnore() - clears the list.
func SetMembersToIgnore(s ...string) {
	if len(s) == 0 {
		skipmembers = skipmembers[:0]
		return
	}
	skipmembers = make([]skipmems, len(s))
	for i, v := range s {
		skipmembers[i] = skipmems{strings.ToLower(v), len(strings.Split(v, "."))}
	}
}

// MissingJSONKeys returns a list of members of val of struct type  that will NOT be set
// by unmarshaling the JSON object; rather, they will assume their default
// values. For nested structs, member labels are the dot-notation hierachical
// path for the missing JSON key.  Specific struct members can be igored
// when scanning the JSON object by declaring them using SetMembersToIgnore().
// (NOTE: JSON object keys are treated as case insensitive, i.e., there
// is no distiction between "key":"value" and "Key":"value".)
//
// By default struct members that have JSON tags "-" and "omitempty" are ignored.
// IgnoreOmitemptyTag(false) can be called to override the handling of "omitempty"
// tags.
func MissingJSONKeys(b []byte, val interface{}) ([]string, error) {
	s := make([]string, 0)
	m := make(map[string]interface{})
	if err := json.Unmarshal(b, &m); err != nil {
		return s, ResolveJSONError(b, err)
	}
	if err := checkMembers(m, reflect.ValueOf(val), &s, ""); err != nil {
		return s, err
	}
	return s, nil
}

// cmem is the parent struct member for nested structs
func checkMembers(mv interface{}, val reflect.Value, s *[]string, cmem string) error {
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
			// cmem is the member name for the slice - []<T> - value
			if err := checkMembers(sl, sval, s, cmem); err != nil {
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
	// 3c. Coerce keys to lower case.
	mkeys := make(map[string]interface{}, len(mm))
	for k, v := range mm {
		mkeys[strings.ToLower(k)] = v
	}

	// 4. Build the list of struct field name:value
	//    We make every key (field) label look like an exported label - "Fieldname".
	//    If there is a JSON tag it is used instead of the field label, and saved to
	//    insure that the spec'd tag matches the JSON key exactly.
	type fieldSpec struct {
		name      string
		val       reflect.Value
		tag       string
		omitempty bool
	}
	fieldCnt := val.NumField()
	fields := make([]*fieldSpec, 0) // use a list so members are in sequence
	for i := 0; i < fieldCnt; i++ {
		if len(typ.Field(i).PkgPath) > 0 {
			continue // field is NOT exported
		}
		var tag string
		var oempty bool
		t := typ.Field(i).Tag.Get("json")
		tags := strings.Split(t, ",")
		tag = tags[0]
		// handle ignore member JSON tag, "-"
		if tag == "-" {
			continue
		}
		// scan rest of tags for "omitempty"
		for _, v := range tags[1:] {
			if v == "omitempty" {
				oempty = true
				break
			}
		}
		if tag == "" {
			fields = append(fields, &fieldSpec{typ.Field(i).Name, val.Field(i), "", oempty})
		} else {
			fields = append(fields, &fieldSpec{typ.Field(i).Name, val.Field(i), tag, oempty})
		}
	}

	// 5. check that field names/tags have corresponding map key
	// var ok bool
	var v interface{}
	var err error
	cmemdepth := 1
	if len(cmem) > 0 {
		cmemdepth = len(strings.Split(cmem, ".")) + 1 // struct hierarchy
	}
	lcmem := strings.ToLower(cmem)
	for _, field := range fields {
		lm := strings.ToLower(field.name)
		for _, sm := range skipmembers {
			// skip any skipmembers values that aren't at same depth
			if cmemdepth != sm.depth {
				continue
			}
			if len(cmem) > 0 {
				if lcmem+`.`+lm == sm.val {
					goto next
				}
			} else if lm == sm.val {
				goto next
			}
		}
		if len(field.tag) > 0 {
			v, ok = mkeys[field.tag]
		} else {
			v, ok = mkeys[lm]
		}
		// If map key is missing, then record it
		// if there's no omitempty tag or we're ignoring  omitempty tag.
		if !ok && (!field.omitempty || !omitemptyOK) {
			if len(cmem) > 0 {
				*s = append(*s, cmem+`.`+field.name)
			} else {
				*s = append(*s, field.name)
			}
			goto next // don't drill down further; no key in JSON object
		}
		if len(cmem) > 0 {
			err = checkMembers(v, field.val, s, cmem+`.`+field.name)
		} else {
			err = checkMembers(v, field.val, s, field.name)
		}
		if err != nil { // could be nested structs
			return fmt.Errorf("checking submembers of member: %s - %s", cmem, err.Error())
		}
	next:
	}

	return nil
}
