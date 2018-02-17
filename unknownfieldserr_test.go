// +build go1.10

package checkjson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

// Per go1.10 (*Decoder)DisallowUnknownFields seems to mimic
// checkjson.UnknownKeys() in its net effect.  What does it
// really do?
func TestDisallowUnknFields(t *testing.T) {
	fmt.Println("===================== TestDisallowUnknFields ...")

	type test2 struct {
		Maybe bool
	}
	type test struct {
		Ok  bool
		Why test2
	}

	tst := test{}
	data := []byte(`{"ok":true, "why":{"maybe":true,"maybenot":false}, "not":"I don't know"}`)
	r := bytes.NewReader(data)
	d := json.NewDecoder(r)
	d.DisallowUnknownFields()
	if err := d.Decode(&tst); err != nil {
		fmt.Println("err ok:", err)
	} else {
		t.Fatal("no err on decode with DisallowUnknownFields")
	}

	tst = test{}
	data = []byte(`{"ok":true, "why":{"maybe":true}}`)
	r = bytes.NewReader(data)
	d = json.NewDecoder(r)
	d.DisallowUnknownFields()
	if err := d.Decode(&tst); err != nil {
		t.Fatal(err)
	}
}
