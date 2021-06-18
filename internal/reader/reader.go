package reader

import (
	"encoding/json"
	"io"
	"os"
)

// ReadJSONFile attempts to open a JSON filename and unmarshals its contents
// into output
func ReadJSONFile(filename string, output interface{}) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	b, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, output)
}
