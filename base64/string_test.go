package base64_test

import (
	"encoding/json"
	"strings"
	"testing"

	base64 "github.com/jhalag/jsonutils/base64"
)

type TestStruct struct {
	RegularString string
	B64String     base64.String
}

func TestMarshal(t *testing.T) {
	ts := TestStruct{
		RegularString: "foo",
		B64String:     "bar",
	}

	// marshal
	enc, err := json.MarshalIndent(ts, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(string(enc), "bar") {
		t.Error("b64 string was found in plaintext")
	}

	// unmarshal
	ts2 := TestStruct{}
	err = json.Unmarshal(enc, &ts2)
	if err != nil {
		t.Fatal(err)
	}

	if ts2.RegularString != "foo" {
		t.Error("sibling string was altered")
	}

	if ts2.B64String != "bar" {
		t.Error("decoded string did not match")
	}
}
