// existingkeys.go - check JSON object against struct definition
// Copyright © 2016-2018 Charles Banning.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package checkjson

import (
	"encoding/json"
	"reflect"
	"strings"
)

// ExistingJSONKeys returns a list of fields of the struct 'val' that WILL BE set
// by unmarshaling the JSON object.  It is the complement of MissingJSONKeys.
// For nested structs, field labels are the dot-notation hierachical
// path for a JSON key.  Specific struct fields can be igored
// when scanning the JSON object by declaring them using SetMembersToIgnore.
// (NOTE: JSON object keys are treated as case insensitive, i.e., there
// is no distiction between "key":"value" and "Key":"value".)
//
// For embedded structs, both the field label for the embedded struct as well
// as the dot-notation label for that struct's fields are included in the list. Thus,
//		type Person struct {
//		   Name NameInfo
//		   Sex  string
//		}
//	
//		type NameInfo struct {
//		   First, Middle, Last string
//		}
//	
//		jobj := []byte(`{"name":{"first":"Jonnie","middle":"Q","last":"Public"},"sex":"unkown"}`)
//		p := Person{}
// 	
//		fields, _ := ExistingKeys(jobj, p)
//		fmt.Println(fields)  // prints: [Name Name.First Name.Middle Name.Last Sex]
//		
// Struct fields that have JSON tag "-" are never returned. Struct fields with the tag 
// attribute "omitempty" will, by default NOT be returned unless the keys exist in the JSON object.
// If you want to know if "omitempty" struct fields are actually in the JSON object, then call
// IgnoreOmitEmptyTag(false) prior to using ExistingJSONKeys.
func ExistingJSONKeys(b []byte, val interface{}) ([]string, error) {
	s := make([]string, 0)
	m := make(map[string]interface{})
	if err := json.Unmarshal(b, &m); err != nil {
		return s, ResolveJSONError(b, err)
	}
	findMembers(m, reflect.ValueOf(val), &s, "")
	return s, nil
}

// cmem is the parent struct member for nested structs
func findMembers(mv interface{}, val reflect.Value, s *[]string, cmem string) {
	// 1. Convert any pointer value.
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val) // convert ptr to struc
	}
	// zero Value?
	if !val.IsValid() {
		return
	}
	typ := val.Type()

	// json.RawMessage is a []byte/[]uint8 and has Kind() == reflect.Slice
	if typ.Name() == "RawMessage" {
		return
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
			// encoding/json must have a JSON array value to decode
			// unlike encoding/xml which will decode a list of elements
			// to a singleton or vise-versa.
			*s = append(*s, typ.Name())
			return
		}
		// 2.1. Check members of JSON array.
		//      This forces all of them to be regular and w/o typos in key labels.
		for _, sl := range slice {
			// cmem is the member name for the slice - []<T> - value
			findMembers(sl, sval, s, cmem)
		}
		return // done with reflect.Slice value
	}

	// 3a. Ignore anything that's not a struct.
	if typ.Kind() != reflect.Struct {
		return // just ignore it - don't look for k:v pairs
	}
	// 3b. map value must represent k:v pairs
	mm, ok := mv.(map[string]interface{})
	if !ok {
		*s = append(*s, typ.Name())
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
	fields := make([]*fieldSpec, 0) // use a list so members are in sequence
	var tag string
	var oempty bool
	for i := 0; i < val.NumField(); i++ {
		tag = ""
		oempty = false
		if len(typ.Field(i).PkgPath) > 0 {
			continue // field is NOT exported
		}
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
	// var err error
	cmemdepth := 1
	if len(cmem) > 0 {
		cmemdepth = len(strings.Split(cmem, ".")) + 1 // struct hierarchy
	}
	lcmem := strings.ToLower(cmem)
	name := ""
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
			name = field.tag
			v, ok = mkeys[field.tag]
		} else {
			name = field.name
			v, ok = mkeys[lm]
		}
		// If map key is missing, then record it
		// if there's no omitempty tag or we're ignoring  omitempty tag.
		if !ok && (!field.omitempty || !omitemptyOK) {
			goto next // don't drill down further; no key in JSON object
		}
		// field exists in JSON object, so add to list
		if len(cmem) > 0 {
			*s = append(*s, cmem+`.`+field.name)
		} else {
			*s = append(*s, field.name)
		}
		if len(cmem) > 0 {
			findMembers(v, field.val, s, cmem+`.`+name)
		} else {
			findMembers(v, field.val, s, name)
		}
	next:
	}
}
