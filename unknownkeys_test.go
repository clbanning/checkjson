package checkjson

import (
	"fmt"
	"testing"
)

func TestUnknownKeys(t *testing.T) {
	fmt.Println("===================== TestUnknownKeys ...")

	data := []byte(`{"ok":true, "why":{"maybe":true,"maybenot":false}, "not":"I don't know"}`)
	check := map[string]bool{"test.why.maybenot":true, "test.not":true}
	type test2 struct {
		Maybe bool
	}
	type test struct {
		Ok  bool
		Why test2
	}

	tv := test{}
	fields, err := UnknownJSONKeys(data, tv)
	if err != nil {
		t.Fatal(err)
	}
	// fmt.Println(fields)

	for _, v := range fields {
		if _, ok := check[v]; !ok {
			t.Fatal("unknown key not in slice:", v)
		}
	}
}

func TestUnknownKeysTag(t *testing.T) {
	fmt.Println("===================== TestUnknownKeysTag ...")

	data := []byte(`{"ok":true, "why":{"maybe":true,"maybenot":false}, "not":"I don't know"}`)
	check := map[string]bool{"test.why.maybenot":true, "test.not":true}
	type test2 struct {
		Val bool `json:"maybe"`
	}
	type test struct {
		Yup  bool `json:"ok"`
		Why test2
	}

	tv := test{}
	fields, err := UnknownJSONKeys(data, tv)
	if err != nil {
		t.Fatal(err)
	}
	// fmt.Println(fields)

	for _, v := range fields {
		if _, ok := check[v]; !ok {
			t.Fatal("unknown key not in slice:", v)
		}
	}
}

func TestUnknownKeysSkip(t *testing.T) {
	fmt.Println("===================== TestUnknownKeysSkip ...")

	data := []byte(`{"ok":true, "why":{"maybe":true,"maybenot":false}, "not":"I don't know"}`)
	SetKeysToIgnore("test.why.maybenot", "test.not")
	defer SetKeysToIgnore("config")
	type test2 struct {
		Maybe bool
	}
	type test struct {
		Ok  bool
		Why test2
	}

	tv := test{}
	fields, err := UnknownJSONKeys(data, tv)
	if err != nil {
		t.Fatal(err)
	}
	// fmt.Println(fields)
	if len(fields) != 0 {
		t.Fatal("fields:", fields)
	}
}

