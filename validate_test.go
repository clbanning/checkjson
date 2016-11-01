package checkjson

import (
	"fmt"
	"testing"
)

func TestStruct(t *testing.T) {
	fmt.Println("===================== TestStruct ...")
	a := []string{}
	err := Validate([]byte{}, a)
	if err == nil {
		t.Fatal("[]string cannot be a struct")
	}
	fmt.Println("err: ok -", err.Error())

	err = Validate([]byte{}, &a)
	if err == nil {
		t.Fatal("[]string cannot be a struct")
	}
	fmt.Println("err: ok -", err.Error())

	type test struct {
		Ok  bool
		Why string
	}
	tv := test{}
	data := []byte(`{"ok":true,"why":"it's a test"}`)
	if err = Validate(data, tv); err != nil {
		t.Fatalf(err.Error())
	}
	if err = Validate(data, &tv); err != nil {
		t.Fatalf(err.Error())
	}
}

func TestStructTag(t *testing.T) {
	fmt.Println("===================== TestStructTag ...")
	var err error

	type test struct {
		Ok  bool
		Why string `json:"whynot"`
	}
	tv := test{}
	data := []byte(`{"ok":true,"whynot":"it's a tag test"}`)
	if err = Validate(data, tv); err != nil {
		t.Fatalf(err.Error())
	}

	data = []byte(`{"ok":true,"why":"it's a tag test"}`)
	if err = Validate(data, tv); err == nil {
		t.Fatalf("didn't catch key error with whynot tag for member")
	}
	fmt.Println("err ok:", err.Error())
}

func TestIgnoreKeys(t *testing.T) {
	fmt.Println("===================== TestIgnoreKeys ...")
	var err error

	type test struct {
		Ok  bool
		Why string
	}
	tv := test{}
	data := []byte(`{"ok":true,"why":"it's a tag test","config":"test"}`)
	if err = Validate(data, tv); err != nil {
		t.Fatalf(err.Error())
	}

	SetKeysToIgnore("config", "cfg")
	data = []byte(`{"ok":true,"why":"it's a tag test","config":"test","cfg":true}`)
	if err = Validate(data, tv); err != nil {
		t.Fatalf(err.Error())
	}

	SetKeysToIgnore()
	data = []byte(`{"ok":true,"why":"it's a tag test","config":"test"}`)
	if err = Validate(data, tv); err == nil {
		t.Fatalf("didn't catch 'conig' key error")
	}
	fmt.Println("err ok:", err.Error())
}
