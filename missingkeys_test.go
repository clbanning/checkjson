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

func TestJSONMissingSkipMems(t *testing.T) {
	fmt.Println("===================== TestJSONMissingSkipMems ...")

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
	SetMembersToIgnore("why", "more.not", "more.another.something")
	defer SetMembersToIgnore()

	mems, err := MissingJSONKeys(data, tv)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(mems) != 0 {
		t.Fatalf(fmt.Sprintf("missing mems: %d - %#v", len(mems), mems))
	}
	fmt.Println("missing keys:", mems)
	SetMembersToIgnore() // reset to default
}

func TestJSONKeysWithIgnoreTags(t *testing.T) {
	fmt.Println("===================== TestJSONKeysWithIgnoreTag ...")

	type test struct {
		Ok  bool
		Why string
		Whynot string `json:"-"`
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
}

func TestJSONKeysWithOmitemptyTags(t *testing.T) {
	fmt.Println("===================== TestJSONKeysWithmitemptyIgnoreTag ...")

	type nested struct {
		Irrel bool
	}
	type test struct {
		Ok  bool
		Why string
		Whynot string `json:",omitempty"`
		Ignored nested `checkjson:"norecurse"`
	}
	tv := test{}
	data := []byte(`{"ok":true,"why":"it's a test","ignored":{"wildcard":"any"}}`)

	IgnoreOmitemptyTag(true)
	mems, err := MissingJSONKeys(data, tv)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(mems) > 0 {
		t.Fatalf(fmt.Sprintf("len(mems) == %d >> %v", len(mems), mems))
	}

	IgnoreOmitemptyTag()
	defer IgnoreOmitemptyTag()
	mems, err = MissingJSONKeys(data, tv)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(mems) != 1 {
		t.Fatalf(fmt.Sprintf("len(mems) == %d >> %v", len(mems), mems))
	}
	if mems[0] != "Whynot" {
		t.Fatalf("MissingJSONKeys did't get: Whynot")
	}
}

