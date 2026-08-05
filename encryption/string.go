package encryption

import (
	"encoding/base64"
	"errors"
	"strconv"
)

var marshalEncryptionKey string = ""

func SetMarshalEncryptionKey(k string) {
	marshalEncryptionKey = k
}

var ErrEncryptionKeyNotSet = errors.New("Marshal Encryption Key not set")

// helper struct to enable storing a field securely in JSON.
//
//   - supports reading of encrypted fields previously saved
//   - can also load plaintext for future re-saving encrypted
//   - encrypted data is stored Base64 encoded - storing a valid B64 string unencrypted will cause unmarshal to be incorrect.
type String struct {
	WasPlaintext bool
	Value        string
}

// Marshal = struct > JSON
func (e String) MarshalJSON() ([]byte, error) {
	if marshalEncryptionKey == "" {
		return nil, ErrEncryptionKeyNotSet
	}
	if e.Value == "" {
		return []byte(`""`), nil // empty string
	}
	encval, err := Encrypt(marshalEncryptionKey, []byte(e.Value), nil)
	if err != nil {
		return nil, err
	}

	b64str := base64.StdEncoding.EncodeToString(encval)

	// wrap the base64 string in quotes so that it is valid JSON
	return []byte("\"" + b64str + "\""), nil
}

// Unmarshal = JSON > struct
func (e *String) UnmarshalJSON(b []byte) error {
	if marshalEncryptionKey == "" {
		return ErrEncryptionKeyNotSet
	}
	if len(b) == 0 || string(b) == `""` {
		e.Value = ""
		e.WasPlaintext = false
		return nil
	}

	// attempt B64 decode
	b64data := make([]byte, base64.StdEncoding.DecodedLen(len(b))) // guess at maximum size the decoded data may be

	chars, err := base64.StdEncoding.Decode(
		b64data,
		b[1:len(b)-1], // trim the quotes from the data sent to b64 decode
	)
	b64data = b64data[:chars] //trim as required

	if err == nil { // was valid B64 data
		unEncStr, err := Decrypt(marshalEncryptionKey, b64data, nil)
		if err != nil {
			return err // decryption error
		}
		e.Value = string(unEncStr)
		e.WasPlaintext = false
		return nil
	}

	// was not blank, and was not valid B64.

	e.Value, _ = strconv.Unquote(string(b))
	e.WasPlaintext = true
	return nil
}
