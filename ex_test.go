package checkjson

import (
	"fmt"
	"testing"
)

func TestEx(t *testing.T) {
	fmt.Println("===================== TestEx ...")

	data := `
	{
		"elem1":"a simple element",
		"elem2": {
			"subelem":"something more complex", 
			"notes"  :"take a look at this" },
	   "elem4":"extraneous" 
	}`

	type sub struct {
		Subelem string `json:"subelem,omitempty"`
		Another string `json:"another"`
	}
	type elem struct {
		Elem1 string `json:"elem1"`
		Elem2 sub    `json:"elem2"`
		Elem3 bool   `json:"elem3"`
	}

	e := new(elem)
	result, err := MissingJSONKeys([]byte(data), e)
	if err != nil {
		t.Fatal("MissingJSONKeys:", err)
	}
	want := `[elem2.another elem3]`
	s := fmt.Sprint(result)
	if s != want {
		t.Fatal("MissingJSONKeys", s, "!=", want)
	}

	result, err = UnknownJSONKeys([]byte(data), e)
	if err != nil {
		t.Fatal("UnknownJSONKeys:", err)
	}
	want = `[elem2.notes elem4]`
	s = fmt.Sprint(result)
	if s != want {
		t.Fatal("UnknownJSONKeys", s, "!=", want)
	}
}
