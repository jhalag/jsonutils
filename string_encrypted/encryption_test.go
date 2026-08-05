package string_encrypted

import (
	"testing"
)

func TestEncryptionRoundtrip(t *testing.T) {
	// matching key, no extra data
	{
		enc, err := Encrypt("foo test1", []byte("test text"), nil)
		if err != nil {
			t.Error(err)
		} else {
			unEnc, err := Decrypt("foo test1", enc, nil)
			if err != nil {
				t.Error(err)
			} else if string(unEnc) != "test text" {
				t.Error("decrypted text did not match")
			}
		}
	}

	// matching key, matching extra data
	{
		enc, err := Encrypt("foo test1", []byte("test text"), []byte("extra data"))
		if err != nil {
			t.Error(err)
		} else {
			unEnc, err := Decrypt("foo test1", enc, []byte("extra data"))
			if err != nil {
				t.Error(err)
			} else if string(unEnc) != "test text" {
				t.Error("decrypted text did not match")
			}
		}
	}

	// matching key, mismatched extradata
	{
		enc, err := Encrypt("foo test1", []byte("test text"), []byte("extra data"))
		if err != nil {
			t.Error(err)
		} else {
			_, err := Decrypt("foo test1", enc, []byte("ALTERED data"))
			if err == nil {
				t.Error("ExtraData was altered, yet decrypted without error")
			}
		}
	}

	// mismatched key, no extradata
	{
		enc, err := Encrypt("foo test1", []byte("test text"), nil)
		if err != nil {
			t.Error(err)
		} else {
			_, err := Decrypt("foo test2", enc, nil)
			if err == nil {
				t.Error("key was wrong, yet decrypted without error")
			}
		}
	}

}
