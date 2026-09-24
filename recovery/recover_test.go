package recovery

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sarthak-shubham/cherry/persistence"
)

func writeTestWAL(t *testing.T, records []persistence.WALRecord) string {
	t.Helper()

	walPath := filepath.Join(t.TempDir(), "data.jsonl")

	for _, record := range records {
		err := persistence.AppendLog(walPath, record)
		if err != nil {
			t.Fatalf("failed to write test WAL: %v", err)
		}
	}

	return walPath
}

func TestRecoverMissingWAL(t *testing.T) {
	walPath := filepath.Join(t.TempDir(), "data.jsonl")

	data, err := Recover(walPath)
	if err != nil {
		t.Fatalf("Recover returned an unexpected error: %v", err)
	}

	if len(data) != 0 {
		t.Fatalf("recovered %d keys, want 0", len(data))
	}
}

func TestRecoverSetRecords(t *testing.T) {
	walPath := writeTestWAL(t, []persistence.WALRecord{
		{
			Operation: "set",
			Key:       "name",
			Value:     "Sarthak",
		},
		{
			Operation: "set",
			Key:       "city",
			Value:     "Delhi",
		},
	})

	data, err := Recover(walPath)
	if err != nil {
		t.Fatalf("Recover returned an unexpected error: %v", err)
	}

	if data["name"] != "Sarthak" {
		t.Fatalf("name = %q, want %q", data["name"], "Sarthak")
	}

	if data["city"] != "Delhi" {
		t.Fatalf("city = %q, want %q", data["city"], "Delhi")
	}
}

func TestRecoverSetOverwritesExistingValue(t *testing.T) {
	walPath := writeTestWAL(t, []persistence.WALRecord{
		{
			Operation: "set",
			Key:       "name",
			Value:     "Sarthak",
		},
		{
			Operation: "set",
			Key:       "name",
			Value:     "Rahul",
		},
	})

	data, err := Recover(walPath)
	if err != nil {
		t.Fatalf("Recover returned an unexpected error: %v", err)
	}

	if data["name"] != "Rahul" {
		t.Fatalf("name = %q, want %q", data["name"], "Rahul")
	}
}

func TestRecoverSetDeleteSetInOrder(t *testing.T) {
	walPath := writeTestWAL(t, []persistence.WALRecord{
		{
			Operation: "set",
			Key:       "name",
			Value:     "Sarthak",
		},
		{
			Operation: "delete",
			Key:       "name",
		},
		{
			Operation: "set",
			Key:       "name",
			Value:     "Rahul",
		},
	})

	data, err := Recover(walPath)
	if err != nil {
		t.Fatalf("Recover returned an unexpected error: %v", err)
	}

	if data["name"] != "Rahul" {
		t.Fatalf("name = %q, want %q", data["name"], "Rahul")
	}
}

func TestRecoverDeleteRemovesKey(t *testing.T) {
	walPath := writeTestWAL(t, []persistence.WALRecord{
		{
			Operation: "set",
			Key:       "name",
			Value:     "Sarthak",
		},
		{
			Operation: "delete",
			Key:       "name",
		},
	})

	data, err := Recover(walPath)
	if err != nil {
		t.Fatalf("Recover returned an unexpected error: %v", err)
	}

	if _, exists := data["name"]; exists {
		t.Fatal("name exists after recovery, want it to be deleted")
	}
}

func TestRecoverDeleteMissingKey(t *testing.T) {
	walPath := writeTestWAL(t, []persistence.WALRecord{
		{
			Operation: "delete",
			Key:       "name",
		},
	})

	data, err := Recover(walPath)
	if err != nil {
		t.Fatalf("Recover returned an unexpected error: %v", err)
	}

	if len(data) != 0 {
		t.Fatalf("recovered %d keys, want 0", len(data))
	}
}

func TestRecoverRejectsUnknownOperation(t *testing.T) {
	walPath := writeTestWAL(t, []persistence.WALRecord{
		{
			Operation: "update",
			Key:       "name",
			Value:     "Sarthak",
		},
	})

	_, err := Recover(walPath)
	if err == nil {
		t.Fatal("Recover returned nil, want an error")
	}
}

func TestRecoverRejectsMalformedWAL(t *testing.T) {
	walPath := filepath.Join(t.TempDir(), "data.jsonl")

	err := os.WriteFile(
		walPath,
		[]byte(`{"op":"set","key":"name","value":"Sarthak"
`),
		0644,
	)
	if err != nil {
		t.Fatalf("failed to write malformed WAL: %v", err)
	}

	_, err = Recover(walPath)
	if err == nil {
		t.Fatal("Recover returned nil, want an error")
	}
}
