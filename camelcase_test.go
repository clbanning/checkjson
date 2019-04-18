package checkjson

import (
	"fmt"
	"testing"
)

func TestValidate(t *testing.T) {
	fmt.Println("===================== camelcase_test#TestValidate")
	type test struct {
		Ok     bool
		Whynot string `json:"whynot"`
	}
	tv := test{}
	data := []byte(`{"ok":false, "whyNot":"it's not a test"}`)
	err := Validate(data, tv)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestMissingJSONKeys(t *testing.T) {
	fmt.Println("===================== camelcase_test#TestMissingJSONKeys")
	type test struct {
		Ok     bool
		Whynot string `json:"whynot"`
	}
	tv := test{}
	data := []byte(`{"ok":false, "whyNot":"it's not a test"}`)
	mems, err := MissingJSONKeys(data, tv)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(mems) > 0 {
		t.Fatalf(fmt.Sprintf("len(mems) == %d >> %v", len(mems), mems))
	}
}

func TestUnknownJSONKeys(t *testing.T) {
	fmt.Println("===================== camelcase_test#TestUnknownJSONKeys")
	type test struct {
		Ok     bool
		Whynot string `json:"whynot"`
	}
	tv := test{}
	data := []byte(`{"ok":false, "whyNot":"it's not a test"}`)
	mems, err := UnknownJSONKeys(data, tv)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(mems) > 0 {
		t.Fatalf(fmt.Sprintf("len(mems) == %d >> %v", len(mems), mems))
	}
}

func TestExistingJSONKeys2(t *testing.T) {
	fmt.Println("===================== camelcase_test#TestExistingJSONKeys2")
	type test struct {
		Ok     bool
		Whynot string `json:"whynot"`
	}
	tv := test{}
	data := []byte(`{"ok":false, "whyNot":"it's not a test"}`)
	mems, err := ExistingJSONKeys(data, tv)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(mems) != 2 {
		t.Fatalf(fmt.Sprintf("len(mems) == %d >> %v", len(mems), mems))
	}
}
