package base64

import (
	"encoding/base64"
)

// helper datatype to enable storing a field that is base64 while in transit, yet transparently (un)marshaled to/from code.
type String string

// Marshal = struct > JSON
func (e String) MarshalJSON() ([]byte, error) {

	b64str := base64.StdEncoding.EncodeToString([]byte(e))

	// wrap the base64 string in quotes so that it is valid JSON
	return []byte("\"" + b64str + "\""), nil
}

// Unmarshal = JSON > struct
func (e *String) UnmarshalJSON(b []byte) error {
	// attempt B64 decode
	b64data := make([]byte, base64.StdEncoding.DecodedLen(len(b))) // guess at maximum size the decoded data may be

	chars, err := base64.StdEncoding.Decode(
		b64data,
		b[1:len(b)-1], // trim the quotes from the data sent to b64 decode
	)
	b64data = b64data[:chars] //trim as required

	if err == nil { // was valid B64 data
		*e = String(b64data)
		return nil
	}

	// was not blank, and was not valid B64.
	return nil
}
