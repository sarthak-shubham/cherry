package store

import (
	"errors"
	"fmt"
	"sync"
	"testing"
)

const (
	concurrentWorkers   = 12
	operationsPerWorker = 100
	concurrentKeyCount  = 10
)

func TestConcurrentStoreOperations(t *testing.T) {
	kvStore := newTestStore(t)

	for i := 0; i < concurrentKeyCount; i++ {
		key := fmt.Sprintf("key-%d", i)

		err := kvStore.Set(key, "initial")
		if err != nil {
			t.Fatalf("failed to seed store: %v", err)
		}
	}

	var waitGroup sync.WaitGroup
	start := make(chan struct{})
	errCh := make(chan error, concurrentWorkers)

	waitGroup.Add(concurrentWorkers)

	for workerID := 0; workerID < concurrentWorkers; workerID++ {
		go runStoreWorker(
			&waitGroup,
			start,
			errCh,
			kvStore,
			workerID,
		)
	}

	close(start)
	waitGroup.Wait()
	close(errCh)

	for err := range errCh {
		t.Error(err)
	}

	for i := 0; i < concurrentKeyCount; i++ {
		key := fmt.Sprintf("key-%d", i)

		err := kvStore.Set(key, "final")
		if err != nil {
			t.Fatalf("final Set failed for %s: %v", key, err)
		}

		value, err := kvStore.Get(key)
		if err != nil {
			t.Fatalf("final Get failed for %s: %v", key, err)
		}

		if value != "final" {
			t.Fatalf("value for %s = %q, want %q", key, value, "final")
		}
	}
}

func runStoreWorker(
	waitGroup *sync.WaitGroup,
	start <-chan struct{},
	errCh chan<- error,
	kvStore *Store,
	workerID int,
) {
	defer waitGroup.Done()

	<-start

	for operation := 0; operation < operationsPerWorker; operation++ {
		keyIndex := (workerID + operation) % concurrentKeyCount
		key := fmt.Sprintf("key-%d", keyIndex)

		switch operation % 3 {
		case 0:
			value := fmt.Sprintf("worker-%d-operation-%d", workerID, operation)

			err := kvStore.Set(key, value)
			if err != nil {
				errCh <- fmt.Errorf(
					"worker %d Set failed on operation %d: %w",
					workerID,
					operation,
					err,
				)
				return
			}

		case 1:
			_, err := kvStore.Get(key)
			if err != nil && !errors.Is(err, ErrKeyNotFound) {
				errCh <- fmt.Errorf(
					"worker %d Get failed on operation %d: %w",
					workerID,
					operation,
					err,
				)
				return
			}

		case 2:
			err := kvStore.Delete(key)
			if err != nil && !errors.Is(err, ErrKeyNotFound) {
				errCh <- fmt.Errorf(
					"worker %d Delete failed on operation %d: %w",
					workerID,
					operation,
					err,
				)
				return
			}
		}
	}
}
