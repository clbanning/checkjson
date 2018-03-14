package checkjson

import (
	"fmt"
	"testing"
)

func TestExistingJSONKeys(t *testing.T) {
	fmt.Println("===================== TestExistingJSONKeys ...")

	type test struct {
		Ok  bool
		Why string
	}
	tv := test{}
	data := []byte(`{"ok":true,"why":"it's a test"}`)
	mems, err := ExistingJSONKeys(data, tv)
	if err != nil {
		t.Fatalf(err.Error())
	}
	// returns: [Ok Why]
	if len(mems) != 2 {
		t.Fatalf(fmt.Sprintf("len(mems) == %d >> %v", len(mems), mems))
	}

	data = []byte(`{"ok":true}`)
	mems, err = ExistingJSONKeys(data, tv)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(mems) != 1 {
		t.Fatalf(fmt.Sprintf("existing mems: %d - %#v", len(mems), mems))
	}
	if mems[0] != "Ok" {
		t.Fatalf(fmt.Sprintf("existing keys: %v", mems))
	}
}
func TestExistingJSONSubKeys(t *testing.T) {
	fmt.Println("===================== TestExistingJSONSubKeys ...")
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
	mems, err := ExistingJSONKeys(data, tv)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(mems) != 8 {
		t.Fatalf(fmt.Sprintf("len(mems) == %d >> %v", len(mems), mems))
	}

	data = []byte(`{"ok":true,"more":{"why":"again","another":{"else":"ok"}}}`)
	mems, err = ExistingJSONKeys(data, tv)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(mems) != 5 {
		t.Fatalf(fmt.Sprintf("existing mems: %d - %#v", len(mems), mems))
	}
	fmt.Println("existing keys:", mems)
}

func TestExistingJSONKeysWithTags(t *testing.T) {
	fmt.Println("===================== TestExistingJSONKeysWithTags ...")

	type test struct {
		Ok  bool
		Why string `json:"whynot"`
	}
	tv := test{}
	data := []byte(`{"ok":true,"whynot":"it's a test"}`)
	mems, err := ExistingJSONKeys(data, tv)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(mems) != 2 {
		t.Fatalf(fmt.Sprintf("len(mems) == %d >> %v", len(mems), mems))
	}
	fmt.Println("existing keys:", mems)

	data = []byte(`{"ok":true,"why":"it's not a test"}`)
	mems, err = ExistingJSONKeys(data, tv)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(mems) != 1 || mems[0] != "Ok" {
		t.Fatalf(fmt.Sprintf("existing mems: %d - %#v", len(mems), mems))
	}
	fmt.Println("existing keys:", mems)
}

func TestExistingJSONKeysSkipMems(t *testing.T) {
	fmt.Println("===================== TestExistingJSONKeysSkipMems ...")

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
	SetMembersToIgnore("ok", "more.why")
	defer SetMembersToIgnore()

	mems, err := ExistingJSONKeys(data, tv)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(mems) != 3 {
		t.Fatalf(fmt.Sprintf("existing mems: %d - %#v", len(mems), mems))
	}
	fmt.Println("existing keys:", mems)
}

func TestExistingJSONKeysWithIgnoreTag(t *testing.T) {
	fmt.Println("===================== TestExistingJSONKeysWithIgnoreTag ...")

	type test struct {
		Ok     bool
		Why    string
		Whynot string `json:"-"`
	}
	tv := test{}
	data := []byte(`{"ok":true,"why":"it's a test"}`)
	mems, err := ExistingJSONKeys(data, tv)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(mems) != 2 {
		t.Fatalf(fmt.Sprintf("len(mems) == %d >> %v", len(mems), mems))
	}
	fmt.Println("existing keys:", mems)
}

func TestExistingJSONKeysWithOmitemptyAttr(t *testing.T) {
	fmt.Println("===================== TestExistingJSONKeysWithmitemptyIgnoreAttr ...")

	type test struct {
		Ok     bool
		Why    string
		Whynot string `json:",omitempty"`
	}
	tv := test{}
	data := []byte(`{"ok":true,"why":"it's a test"}`)

	IgnoreOmitemptyTag(true) // make sure it's set
	mems, err := ExistingJSONKeys(data, tv)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(mems) != 3 {
		t.Fatalf(fmt.Sprintf("len(mems) == %d >> %v", len(mems), mems))
	}
	if mems[2] != "Whynot" {
		t.Fatalf("ExistingJSONKeys did't get: Whynot")
	}

	IgnoreOmitemptyTag(false)      // ignore attribute
	defer IgnoreOmitemptyTag(true) // reset on return
	mems, err = ExistingJSONKeys(data, tv)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(mems) != 2 {
		t.Fatalf(fmt.Sprintf("len(mems) == %d >> %v", len(mems), mems))
	}
}
