package filesort

import (
	"testing"
	"log"
	"os"
)

func TestFileSortInMemory(t *testing.T) {
	log.SetOutput(os.Stdout)

	filePath := "./test-file"
	err := GenerateFile(filePath, 10, 5)
	if err != nil {
		t.Error(err)
		return
	}
	defer RemoveFile(filePath)

	ok, err := IsFileSorted(filePath)
	if err != nil {
		t.Error(err)
		return
	}
	if ok {
		t.Fatal("file already sorted")
		return
	}

	err = SortFileInMemory(filePath)
	if err != nil {
		t.Error(err)
		return
	}

	ok, err = IsFileSorted(filePath)
	if err != nil {
		t.Error(err)
		return
	}

	if !ok {
		t.Fatal("file is not sorted")
	}
}

func TestFileSort(t *testing.T) {
	log.SetOutput(os.Stdout)

	filePath := "./test-file"
	err := GenerateFile(filePath, 999971, 5)
	if err != nil {
		t.Error(err)
		return
	}
	defer RemoveFile(filePath)

	ok, err := IsFileSorted(filePath)
	if err != nil {
		t.Error(err)
		return
	}
	if ok {
		t.Fatal("file already sorted")
		return
	}

	err = SortFile(filePath)
	if err != nil {
		t.Error(err)
		return
	}

	ok, err = IsFileSorted(filePath)
	if err != nil {
		t.Error(err)
		return
	}

	if !ok {
		t.Fatal("file is not sorted")
	}
}
