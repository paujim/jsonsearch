package jsonsearch

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

var test1 []byte
var test2 []byte

func init() {
	loadFile := func(fileName string) []byte {
		jsonFile, _ := os.Open(fileName)
		byteValue, _ := ioutil.ReadAll(jsonFile)
		return byteValue
	}

	test1 = loadFile("test1.json")
	test2 = loadFile("test2.json")
}

func TestNewJsonObjectWithInvalidParameters(t *testing.T) {
	var buff bytes.Buffer
	logger := log.New(
		&buff,
		"DEBUG:",
		log.Ldate|log.Ltime|log.Lshortfile,
	)
	actual := NewJsonObject(logger, nil, "JsonName", "")
	if actual != nil {
		t.Errorf("expected: a nil response")
	}

	if !strings.Contains(buff.String(), "Creating JsonName") {
		t.Errorf("expected: the log shoudl contain [Creating JsonName]")
	}
}

func TestJsonObject(t *testing.T) {
	// Test valid object
	obj := NewJsonObject(nil, test1, "test", "id")
	if obj == nil {
		t.Errorf("expected: a valid json object")
	}

	// Test index
	if obj.IsIdField("SOMETHIN") && !obj.IsIdField("id") {
		t.Errorf("expected: idField is [id]")
	}
	_, ok := obj.GetById("1")
	if !ok {
		t.Errorf("expected: [id = 1] should exist")
	}
	_, ok = obj.GetById("31")
	if ok {
		t.Errorf("expected: [id = 31] should NOT exist")
	}

	// Test Properties
	keys := obj.Keys()

	if len(keys) != 5 {
		t.Errorf("expected: 5 keys")
	}

	// Test Searching
	areEqual := func(s1 string, s2 string) bool {
		return s1 == s2
	}
	actualResult := obj.SearchAll("details", "MegaCorp", areEqual)
	if len(actualResult) != 1 {
		t.Errorf("expected: 1 values")
	}
	actualResult = obj.SearchAll("tags", "Fuentes", areEqual)
	if len(actualResult) != 2 {
		t.Errorf("expected: 2 values")
	}

	actualResult = obj.SearchAll("details", "non found", areEqual)
	if len(actualResult) != 0 {
		t.Errorf("expected: 0 values")
	}

	actualResult = obj.SearchAll("non_existing", "non found", areEqual)
	if len(actualResult) != 0 {
		t.Errorf("expected: 0 values")
	}
}
