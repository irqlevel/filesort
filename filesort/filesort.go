package filesort

import (
	"os"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strings"
	"bufio"
	"io"
)

func generateString(lineLen int) (string, error) {
	b := make([]byte, lineLen/2 + 1)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b)[:lineLen], nil
}

func RemoveFile(filePath string) error {
	filePath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("Can't get abs path of %s error %v", filePath, err)
	}

	return os.Remove(filePath)
}

func GenerateFile(filePath string, numLines int64, lineLen int) error {
	filePath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("Can't get abs path of %s error %v", filePath, err)
	}

	if numLines < 0 || lineLen < 0 {
		return fmt.Errorf("Invalid numLine %d or lineLen %d specified", numLines, lineLen)
	}

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer func() {
		f.Close()
		if err != nil {
			os.Remove(filePath)
		}
	} ()

	err = f.Truncate(0)
	if err != nil {
		return fmt.Errorf("Can't truncate file %s error %v", filePath, err)
	}

	for i := int64(0); i < numLines; i++ {
		var s string
		s, err = generateString(lineLen)
		if err != nil {
			return err
		}

		if _, err = f.WriteString(strings.Join([]string{s, "\n"}, "")); err != nil {
			return err
		}
	}
	return nil
}

func SortFile(filePath string) error {
	filePath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("Can't get abs path of %s error %v", filePath, err)
	}

	f, err := os.OpenFile(filePath, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

func IsFileSorted(filePath string) (bool, error) {
	filePath, err := filepath.Abs(filePath)
	if err != nil {
		return false, fmt.Errorf("Can't get abs path of %s error %v", filePath, err)
	}

	f, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {
		return false, err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	prev := ""
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return false, fmt.Errorf("Can't read file %s error %v", filePath, err)
		}

		if strings.Compare(prev, s) > 0 {
			return false, nil
		}
		prev = s
	}

	return true, nil
}
