package filesort

import (
	"fmt"
	"log"
	"os"
	"testing"
)

//test file sorting in memory
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

type sortParams struct {
	numLines int64
	lineLen  int
	maxLines int
}

func testFileSortWithParams(params sortParams) error {
	filePath := "./test-file"
	err := GenerateFile(filePath, params.numLines, params.lineLen)
	if err != nil {
		return err
	}
	defer RemoveFile(filePath)

	ok, err := IsFileSorted(filePath)
	if err != nil {
		return err
	}

	if ok {
		return fmt.Errorf("file already sorted")
	}

	err = SortFile(filePath, params.maxLines)
	if err != nil {
		return err
	}

	ok, err = IsFileSorted(filePath)
	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("file not sorted")
	}

	return nil
}

//test file sorting with different parameters numLines, lineLen, maxLines(memory limit)
func TestFileSort(t *testing.T) {
	log.SetOutput(os.Stdout)

	params := []sortParams{{10, 5, 2},
		{7, 3, 2},
		{105097, 73, 541}}

	for _, param := range params {
		err := testFileSortWithParams(param)
		if err != nil {
			t.Fatalf("sort with params %v failed error %v", param, err)
			return
		}
	}
}
