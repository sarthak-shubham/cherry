package recovery

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/sarthak-shubham/cherry/persistence"
)

func Recover(path string) (map[string]string, error) {
	data := make(map[string]string)

	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return data, nil
		}

		return nil, err
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	var lastGoodOffset int64

	for {
		line, err := reader.ReadBytes('\n')

		if len(line) == 0 && errors.Is(err, io.EOF) {
			break
		}

		if errors.Is(err, io.EOF) {
			err = file.Truncate(lastGoodOffset)

			if err != nil {
				return nil, fmt.Errorf("failed to truncate incomplete WAL tail: %w", err)
			}

			break
		}

		if err != nil {
			return nil, err
		}

		line = line[:len(line)-1]

		var record persistence.WALRecord

		err = json.Unmarshal(line, &record)
		if err != nil {
			return nil, fmt.Errorf(
				"invalid WAL record at offset %d: %w",
				lastGoodOffset,
				err,
			)
		}

		switch record.Operation {
		case "set":
			data[record.Key] = record.Value

		case "delete":
			delete(data, record.Key)

		default:
			return nil, fmt.Errorf("unknown WAL operation: %q", record.Operation)
		}

		lastGoodOffset += int64(len(line) + 1)
	}

	return data, nil
}
