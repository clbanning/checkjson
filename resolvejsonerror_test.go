package checkjson

import (
	"encoding/json"
	"fmt"
	"testing"
)

func  TestResolveJSONError(t *testing.T) {
	fmt.Println("=============== TestResolveJSONError")
	// this is just to view formating w/ "key" and position info
	src := `{"this":{"is":["a", "test"],"of":"ResolveJSONError", "with":"a","quote":missing}}`
	var m interface{}
	err := json.Unmarshal([]byte(src), m)
	if err == nil {
		t.Fatal("JSON error not caught")
	}
	fmt.Println(ResolveJSONError([]byte(src), err))

	src = `{"this":{"is":["a", "test"],"of":"ResolveJSONError", "with":"a","quote:missing}}`
	err = json.Unmarshal([]byte(src), m)
	if err == nil {
		t.Fatal("JSON error not caught")
	}
	fmt.Println(ResolveJSONError([]byte(src), err))

	src = `{"this":{"is":["a", "test"],"of":"ResolveJSONError", "with":"a", quote:missing}}`
	err = json.Unmarshal([]byte(src), m)
	if err == nil {
		t.Fatal("JSON error not caught")
	}
	fmt.Println(ResolveJSONError([]byte(src), err))
}
