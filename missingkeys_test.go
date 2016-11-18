package checkjson

import (
	"fmt"
	"testing"
)

func TestJSONKeys(t *testing.T) {
	fmt.Println("===================== TestJSONKeys ...")

	type test struct {
		Ok  bool
		Why string
	}
	tv := test{}
	data := []byte(`{"ok":true,"why":"it's a test"}`)
	mems, err := MissingJSONKeys(data, tv)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(mems) > 0 {
		t.Fatalf(fmt.Sprintf("len(mems) == %d >> %v", len(mems), mems))
	}

	data = []byte(`{"ok":true}`)
	mems, err = MissingJSONKeys(data, tv)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(mems) != 1 {
		t.Fatalf(fmt.Sprintf("missing mems: %d - %#v", len(mems), mems))
	}
	fmt.Println("missing keys:", mems)
}

func TestJSONSubKeys(t *testing.T) {
	fmt.Println("===================== TestJSONSubKeys ...")
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
	data := []byte(`{"ok":true,"why":"it's a test","more":{"why":"again","not":"no more","another":{"something":"a thing","else":"ok"}}}`)
	mems, err := MissingJSONKeys(data, tv)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(mems) > 0 {
		t.Fatalf(fmt.Sprintf("len(mems) == %d >> %v", len(mems), mems))
	}

	data = []byte(`{"ok":true,"more":{"why":"again","another":{"else":"ok"}}}`)
	mems, err = MissingJSONKeys(data, tv)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(mems) != 3 {
		t.Fatalf(fmt.Sprintf("missing mems: %d - %#v", len(mems), mems))
	}
	fmt.Println("missing keys:", mems)
}

func TestJSONKeysWithTags(t *testing.T) {
	fmt.Println("===================== TestJSONKeysWithTags ...")

	type test struct {
		Ok  bool
		Why string `json:"whynot"`
	}
	tv := test{}
	data := []byte(`{"ok":true,"whynot":"it's a test"}`)
	mems, err := MissingJSONKeys(data, tv)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(mems) > 0 {
		t.Fatalf(fmt.Sprintf("len(mems) == %d >> %v", len(mems), mems))
	}

	data = []byte(`{"ok":true,"why":"it's not a test"}`)
	mems, err = MissingJSONKeys(data, tv)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(mems) != 1 {
		t.Fatalf(fmt.Sprintf("missing mems: %d - %#v", len(mems), mems))
	}
	fmt.Println("missing keys:", mems)
}

