package encryption_test

import (
	"encoding/json"
	"strings"
	"testing"

	encryption "github.com/jhalag/jsonutils/encryption"
)

type TestStruct struct {
	RegularString string
	EncString     encryption.String
}

func TestMarshal(t *testing.T) {
	encryption.SetMarshalEncryptionKey("TESTKEYHEAR")
	ts := TestStruct{
		RegularString: "foo",
	}
	ts.EncString.Value = "bar"

	// marshal
	enc, err := json.MarshalIndent(ts, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(string(enc), "bar") {
		t.Error("encrypted string was found in plaintext")
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

	if ts2.EncString.Value != "bar" {
		t.Error("Encrypted string did not match")
	}
	if ts2.EncString.WasPlaintext {
		t.Error("Encrypted string was not flagged as having been encrypted")
	}
}

func TestUnmarshalRaw(t *testing.T) {
	encryption.SetMarshalEncryptionKey("TESTKEYHEAR")
	rawstr := `
		{
			"EncString":"raw str here"
		}
	`

	// unmarshal
	ts := TestStruct{}
	err := json.Unmarshal([]byte(rawstr), &ts)
	if err != nil {
		t.Fatal(err)
	}

	if ts.EncString.Value != "raw str here" {
		t.Error("Encrypted string did not match")
	}
	if !ts.EncString.WasPlaintext {
		t.Error("Encrypted string was not flagged as having been plaintext")
	}
}
