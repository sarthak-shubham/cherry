package recovery

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/sarthak-shubham/cherry/persistence"
)

func Recover(path string) (map[string]string, error) {
	data := make(map[string]string)

	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return data, nil
		}

		return nil, err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	for {
		var record persistence.WALRecord

		err := decoder.Decode(&record)
		if errors.Is(err, io.EOF) {
			break
		}

		switch record.Operation {
		case "set":
			data[record.Key] = record.Value

		case "delete":
			delete(data, record.Key)

		default:
			return nil, fmt.Errorf("unknown WAL operation: %q", record.Operation)
		}
	}

	return data, nil
}
