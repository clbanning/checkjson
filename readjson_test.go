package checkjson

import (
	"fmt"
	"os"
	"testing"
)

var data = []string{
	`{"test":"data"}`,
	`{"test2":"2","again":"3","and more":"5"}`,
	`{"test3":[{"sub1":"s1"},{"sub2":"s2","sub3":"s3","3":[{"here":"again"}]}],"some":"more","sentence":"Now is the time for all good men ..."}`,
	`{"encrypt":"Encrypt","decrypt":"Decrypt","files":[{"encryptfile":"EncryptFile","encryptjsonfile":"EncryptJsonFile#","decryptfile":"DecryptFile\n","decryptjsonfile":"DecryptJsonFile\""}]}`,
	`{"myname":"is","Inigo":"Montoya"}`,
	`{"key":"value","key2":{"key3":"value2","key4":"value4"}}`,
}

// read in a file: should see if it will unmarshal properly, then write it
// reread it and compare with original -
// read/write on Buffer are implicit, since used by JsonFile functions
func TestReadJSONFile(t *testing.T) {
	fmt.Println("============= TestReadJSONFile ...")

	ss, err := ReadJSONFile("data.json")
	if err != nil {
		t.Errorf("ReadJsonFile err: %s", err.Error())
	}

	for i := range ss {
		if string(ss[i]) != data[i] {
			t.Errorf("rwjson ERROR: string mismatch.\nin >>%s\nout>>%s", string(ss[i]), data[i])
		}
	}
}

func TestReadJSONReader(t *testing.T) {
	fmt.Println("============= TestReadJSONReader ...")

	fh, err := os.Open("data.json")
	if err != nil {
		t.Errorf("err opening data.json: %s", err.Error())
	}
	defer fh.Close()

	for i := 0; i < len(data); i++ {
		j, err := ReadJSONReader(fh)
		if err != nil {
			t.Errorf("ReadJSONFromReader err: %s", err.Error())
		}
		if string(j) != data[i] {
			t.Errorf("rwjson ERROR: string mismatch.\nin >>%s\nout>>%s", string(j), data[i])
		}
	}
}

func TestBadJSON(t *testing.T) {
	fmt.Println("============= TestBadJSON ...")

	_, err := ReadJSONFile("baddata.json")
	if err == nil {
		t.Fatal("didn't catch error")
	}
	fmt.Println("err ok:", err)
}
