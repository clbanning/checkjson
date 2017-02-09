package checkjson

import (
	"fmt"
)

func ExampleMissingJSONKeys() {
	// struct to which we want to decode JSON object
	type test3 struct {
		Something string
		Else      string
	}
	type test2 struct {
		Why     string
		Not     string
		Another test3
	}
	type test struct {
		Ok   bool
		Why  string
		More test2
	}

	tv := test{}
	data := []byte(`{"ok":true,"more":{"why":"again","another":{"else":"ok"}}}`)

	mems, err := MissingJSONKeys(data, tv)
	if err != nil {
		// handle error
	}
	fmt.Println("missing keys:", mems)
	// Output:
	// missing keys: [Why More.Not More.Another.Something]
}

func ExampleSetMembersToIgnore() {
	// struct to which we want to decode JSON object
	type test3 struct {
		Something string
		Else      string
	}
	type test2 struct {
		Why     string
		Not     string
		Another test3
	}
	type test struct {
		Ok   bool
		Why  string
		More test2
	}

	data := []byte(`{"ok":true,"more":{"why":"again","another":{"else":"ok"}}}`)
	SetMembersToIgnore("why", "more.not", "more.another.something")
	defer SetMembersToIgnore()

	tv := test{}
	mems, err := MissingJSONKeys(data, tv)
	if err != nil {
		// handle error
	}
	fmt.Println("missing keys:", mems)
	// Output:
	// missing keys: []
}

func ExampleSetKeysToIgnore() {
	// struct to which we want to decode JSON object
	type test2 struct {
		Maybe bool
	}
	type test struct {
		Ok  bool
		Why test2
	}

	tv := test{}
	data := []byte(`{"ok":true, "why":{"maybe":true,"maybenot":false}, "not":"I don't know"}`)
	SetKeysToIgnore("why.maybenot", "not")
	defer SetKeysToIgnore("")

	keys, err := UnknownJSONKeys(data, tv)
	if err != nil {
		// handle error
	}

	fmt.Println("unknown keys:", keys)
	// Output:
	// unknown keys: []
}
