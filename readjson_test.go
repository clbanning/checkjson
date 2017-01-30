package checkjson

import (
	"fmt"
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
func TestReadJSON(t *testing.T) {
	fmt.Println("\n============= TestReadJSON ...")

	ss, err := ReadJSONFile("data.json")
	if err != nil {
		t.Errorf("ReadJsonFile ERROR: %s in %s", err.Error(), err.Error())
	}

	for i := range ss {
		if string(ss[i]) != data[i] {
			t.Errorf("rwjson ERROR: string mismatch.\nin >>%s\nout>>%s", string(ss[i]), data[i])
		}
	}
}

func TestBadJSON(t *testing.T) {
	fmt.Println("\n============= TestBadJSON ...")

	_, err := ReadJSONFile("baddata.json")
	if err == nil {
		t.Fatal("didn't catch error")
	}
	fmt.Println("err ok:", err)
}
