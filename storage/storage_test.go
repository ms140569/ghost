package storage

import (
	"testing"
)

func TestStorageFactoryMem(t *testing.T) {

	provider, err := New("mem")

	if err != nil {
		t.Fatalf("Could not create memory storage provider.")
	}

	provider.Initialize()

}

func TestStorageFactoryFile(t *testing.T) {

	provider, err := New("file")

	if err != nil {
		t.Fatalf("Could not create file storage provider.")
	}

	provider.Initialize()

}
