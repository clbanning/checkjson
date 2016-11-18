package checkjson

import (
	"fmt"
)

func ExampleMissingJSONKeys() {
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

	tv := test{}
	mems, err := MissingJSONKeys(data, tv)
	if err != nil {
		// handle error
	}
	fmt.Println("missing keys:", mems)
	// Output:
	// missing keys: []
}

