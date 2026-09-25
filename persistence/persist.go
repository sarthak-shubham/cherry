package persistence

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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

	info, err := file.Stat()
	if err != nil {
		return err
	}

	startSize := info.Size()

	data, err := json.Marshal(record)
	if err != nil {
		return err
	}

	data = append(data, '\n')

	written, err := file.Write(data)

	if err != nil || written != len(data) {
		truncateErr := file.Truncate(startSize)

		if truncateErr != nil {
			if err != nil {
				return fmt.Errorf("WAL write failed: %v; WAL rollback failed: %w", err, truncateErr)
			}

			return fmt.Errorf("WAL write failed: %w; WAL rollback failed: %v", io.ErrShortWrite, truncateErr)
		}

		if err != nil {
			return err
		}

		return io.ErrShortWrite
	}

	err = file.Sync()

	if err != nil {
		log.Printf("warning: WAL sync failed: %v", err)
	}

	return nil
}
