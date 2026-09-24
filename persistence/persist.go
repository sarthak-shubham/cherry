package persistence

import (
	"encoding/json"
	"os"
)

type WALRecord struct {
	Operation string `json:"op"`
	Key       string `json:"key"`
	Value     string `json:"value,omitempty"`
}

func AppendLog(path string, record WALRecord) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)

	err = encoder.Encode(record)
	if err != nil {
		return err
	}

	return file.Sync()
}
