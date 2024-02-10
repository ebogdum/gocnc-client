package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// FileMode is a custom type that wraps os.FileMode
type FileMode os.FileMode

// UnmarshalJSON customizes the unmarshalling of a file mode from JSON
func (fm *FileMode) UnmarshalJSON(data []byte) error {
	// Attempt to unmarshal as a string
	var fileModeStr string
	if err := json.Unmarshal(data, &fileModeStr); err == nil {
		// If successful, parse the string assuming it's in octal
		mode, err := strconv.ParseUint(fileModeStr, 8, 32)
		if err != nil {
			return err
		}
		*fm = FileMode(mode)
		return nil
	}

	// Attempt to unmarshal as an integer
	var fileModeInt uint64
	if err := json.Unmarshal(data, &fileModeInt); err == nil {
		*fm = FileMode(fileModeInt)
		return nil
	}

	// If neither string nor integer, return an error
	return fmt.Errorf("file mode should be a string or an integer")
}
