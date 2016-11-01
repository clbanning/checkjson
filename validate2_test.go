package checkjson

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

type struct1 struct {
	This string
	Is   struct2
	Not  []struct2
	But  map[string]interface{}
	anon int
}

type struct2 struct {
	A    string
	Json string
}

func TestJdata(t *testing.T) {
	fmt.Println("===================== TestJdata ...")
	s := new(struct1)

	json1 := []byte(`{
			"this":"is",
			"is":{"a":"simple","json":"object"},
			"not":[
				{"a":"simple"},
				{"JSON":"object"},
				{"a":"not so","json":"simple"}
			],
			"but": {"it":"is a", "little":"goofy"}}`)
	if err := Validate(json1, s); err != nil {
		t.Fatal(err)
	}

	// field doesn't exist for JSON key 'else'
	json2 := []byte(`{
			"this":"is",
			"not":[
				{"a":"simple"},
				{"JSON":"object"}
			],
			"but": {"it":"is a", "little":"goofy"},
			"else": false}`)
	// json2 is valid encoding of struct1
	if err := Validate(json2, s); err == nil {
		t.Fatal("no error returned")
	} else {
		fmt.Println("err ok:", err)
	}

	// field doesn't exist for JSON key 'else' in member struct
	json3 := []byte(`{
			"this":"is",
			"not":[
				{"a":"simple"},
				{"json":"object"},
				{"else": false}
			],
			"but": {"it":"is a", "little":"goofy"}}`)
	// json2 is valid encoding of struct1
	if err := Validate(json3, s); err == nil {
		t.Fatal("no error returned")
	} else {
		fmt.Println("err ok:", err)
	}

	json4 := []byte(`{
			"this":"is",
			"is":{"a":"simple","json":"object", "else":false},
			"not":[
				{"a":"simple"},
				{"JSON":"object"},
				{"a":"not so","json":"simple"}
			],
			"but": {"it":"is a", "little":"goofy"}}`)
	if err := Validate(json4, s); err == nil {
		t.Fatal("no error returned")
	} else {
		fmt.Println("err ok:", err)
	}

	json5 := []byte(`{
			"this":"is",
			"is":{"a":"simple","json":"object"},
			"not": {"a":"simple"},
			"but": {"it":"is a", "little":"goofy"}}`)
	if err := Validate(json5, s); err == nil {
		t.Fatal("no error returned")
	} else {
		fmt.Println("err ok:", err)
	}
}

type VmonConns struct {
	Display bool
	Active  bool
	Conns   map[string]bool
}
type VmonSpec struct {
	topology         string
	Id               string
	Defined          bool
	View             string
	Passive          bool
	Status           *VmonConns
	OldestRecallTime string
	EncodingMask     string
	KeyFilter        json.RawMessage
	BumpFilter       json.RawMessage
	SortBy           json.RawMessage
	BumpOtherVmons   json.RawMessage
	BumpItems        bool
	BumpDependents   bool
}
type TopoSpec struct {
	Id        string
	ConfigId  string
	Startup   bool
	StartTime time.Time
	StopTime  time.Time
	Vmons     []*VmonSpec
	vmons     map[string]*VmonSpec
}

func TestKdata(t *testing.T) {
	fmt.Println("===================== TestKdata ...")
	kdata1 := []byte(`
{
	"config"  : "topo",
	"id"      : "expos",
	"startup" : false,
	"vmons"   : [
		{
			"id"             : "expeditor_all",
			"view"           : "item",
			"keyfilter"      : {"child":false,"modifier":false},
			"bumpfilter"     : {"orderbumped":false},
			"bumpdependents" : true,
			"bumpothervmons" : ["grill_dt","fryer_dt","grill","fryer","expeditor","expeditor_dt"],
			"sortby"         : ["orderatime","priority","line"],
			"encodingmask"   : "QSRItem_Hierarchical"
		},
		{
			"id"             : "expeditor_order",
			"view"           : "order",
			"keyfilter"      : {"dest":["EAT IN","CARRY OUT", "DRIVE THRU"]},
			"bumpfilter"     : {"vmonbumped":false},
			"bumpitems"      : true,
			"sortby"         : ["atime"],
			"encodingmask"   : "QSROrder_Hierarchical"
		}
	]
}`)
	tspec := new(TopoSpec)
	if err := Validate(kdata1, tspec); err != nil {
		t.Fatal(err)
	}

	kdata2 := []byte(`
{
	"config"  : "topo",
	"id"      : "expos",
	"startup" : false,
	"vmons"   : [
		{
			"id"             : "expeditor_all",
			"view"           : "item",
			"keyfilter"      : {"child":false,"modifier":false},
			"bumpfilter"     : {"orderbumped":false},
			"bumpdependents" : true,
			"bumpothervmons" : ["grill_dt","fryer_dt","grill","fryer","expeditor","expeditor_dt"],
			"sortby"         : ["orderatime","priority","line"],
			"encodingmask"   : "QSRItem_Hierarchical"
		},
		{
			"not_id"         : "expeditor_order",
			"view"           : "order",
			"keyfilter"      : {"dest":["EAT IN","CARRY OUT", "DRIVE THRU"]},
			"bumpfilter"     : {"vmonbumped":false},
			"bumpitems"      : true,
			"sortby"         : ["atime"],
			"encodingmask"   : "QSROrder_Hierarchical"
		}
	]
}`)
	if err := Validate(kdata2, tspec); err == nil {
		t.Fatal("no error reported")
	} else {
	fmt.Println("err ok:", err)
	}
}
