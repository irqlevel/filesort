package filesort

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

//generate hex encoded string like 'abc0312' with len LineLen
func generateString(lineLen int) (string, error) {
	b := make([]byte, lineLen/2+1)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b)[:lineLen], nil
}

//write one line into output file
func writeLine(f *os.File, l string) error {
	_, err := f.WriteString(strings.Join([]string{l, "\n"}, ""))
	return err
}

//write multiple lines into output file
func writeLines(f *os.File, lines []string) error {
	for _, s := range lines {
		err := writeLine(f, s)
		if err != nil {
			return err
		}
	}
	return nil
}

//read line from input file
func readLine(r *bufio.Reader) (string, bool, error) {
	s, err := r.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return "", true, nil
		}
		return "", false, err
	}
	return strings.TrimSuffix(s, "\n"), false, nil
}

//read multiple lines from input file
func readLines(f *bufio.Reader, maxLines int) ([]string, error) {
	lines := make([]string, 0)
	for i := 0; i < maxLines; i++ {
		s, eof, err := readLine(f)
		if err != nil {
			return nil, err
		}
		if eof {
			break
		}
		lines = append(lines, s)
	}

	return lines, nil
}

//delete file from disk
func RemoveFile(filePath string) error {
	filePath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("Can't get abs path of %s error %v", filePath, err)
	}

	return os.Remove(filePath)
}

//generate file with numLines lines with width lineLen
func GenerateFile(filePath string, numLines int64, lineLen int) error {
	filePath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("Can't get abs path of %s error %v", filePath, err)
	}

	if numLines < 0 || lineLen < 0 {
		return fmt.Errorf("Invalid numLine %d or lineLen %d specified", numLines, lineLen)
	}

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer func() {
		f.Close()
		if err != nil {
			os.Remove(filePath)
		}
	}()

	for i := int64(0); i < numLines; i++ {
		var s string
		s, err = generateString(lineLen)
		if err != nil {
			return err
		}

		err = writeLine(f, s)
		if err != nil {
			return err
		}
	}
	return nil
}

//get temporary run file path by it's id
func getRunFilePath(i int64) string {
	return "./run_" + strconv.FormatInt(i, 10)
}

//delete temporary run file by it's id
func removeRun(i int64) {
	os.Remove(getRunFilePath(i))
}

//delete temporary run files with ids [start, end)
func removeRuns(start int64, end int64) {
	for i := start; i < end; i++ {
		removeRun(i)
	}
}

//merge sorted lines from r[0] r[1] input files into output file
func mergeLines(r [2]*bufio.Reader, out *os.File, fn [3]int64) error {
	var eof [2]bool
	var line [2]string
	var lineValid [2]bool
	var err error

	//check input files eofs, read lines if not eofs and merge them
	//into out file
	for !eof[0] || !eof[1] {
		for i := 0; i < 2; i++ {
			if !lineValid[i] && !eof[i] {
				line[i], eof[i], err = readLine(r[i])
				if err != nil {
					return err
				}
				if !eof[i] {
					lineValid[i] = true
				}
			}
		}

		var s string

		if lineValid[0] && !lineValid[1] {
			s = line[0]
			lineValid[0] = false
		} else if lineValid[1] && !lineValid[0] {
			s = line[1]
			lineValid[1] = false
		} else if lineValid[0] && lineValid[1] {
			if strings.Compare(line[0], line[1]) < 0 {
				s = line[0]
				lineValid[0] = false
			} else {
				s = line[1]
				lineValid[1] = false
			}
		} else {
			continue
		}

		err = writeLine(out, s)
		if err != nil {
			return err
		}
	}

	return nil
}

//merge two temporary run files into temporary out run file
func mergeTwoRuns(first int64, second int64, out int64) error {

	var err error
	fr1, err := os.OpenFile(getRunFilePath(first), os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer fr1.Close()

	fr2, err := os.OpenFile(getRunFilePath(second), os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer fr2.Close()

	fout, err := os.OpenFile(getRunFilePath(out), os.O_WRONLY|os.O_CREATE, 666)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			removeRun(out)
		}
	}()
	defer fout.Close()

	return mergeLines([2]*bufio.Reader{bufio.NewReader(fr1), bufio.NewReader(fr2)}, fout,
		[3]int64{first, second, out})
}

//merge temporary run files with id[0, numRuns) into output file(outFilePath)
func mergeRuns(numRuns int64, outFilePath string) error {
	start := int64(0)
	end := numRuns

	defer func() {
		removeRuns(start, end)
	}()

	//merge run files while not have one run file
	for (end - start) > 1 {
		err := mergeTwoRuns(start, start+1, end)
		if err != nil {
			return err
		}
		removeRuns(start, start+2)
		start += 2
		end += 1
	}

	//rename last run file into output file
	err := os.Rename(getRunFilePath(start), outFilePath)
	if err != nil {
		return err
	}
	start++
	return nil
}

//generate neccessary temprorary run files based on each run file
//contains not more than maxLines lins
func generateRuns(filePath string, maxLines int) (int64, error) {
	if maxLines < 2 {
		return 0, fmt.Errorf("Invalid maxLines %d", maxLines)
	}

	f, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	fr := bufio.NewReader(f)
	i := int64(0)

	defer func() {
		if err != nil {
			removeRuns(0, i)
		}
	}()

	for ; i < math.MaxInt64; i++ {
		lines, err := readLines(fr, maxLines)
		if err != nil {
			return 0, err
		}

		if len(lines) == 0 {
			break
		}

		sort.Strings(lines)

		fr, err := os.OpenFile(getRunFilePath(i), os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return 0, err
		}

		err = writeLines(fr, lines)
		fr.Close()
		if err != nil {
			return 0, err
		}
	}

	return i, nil
}

//sort file with memory limit (maxLines)
func SortFile(filePath string, maxLines int) error {
	filePath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("Can't get abs path of %s error %v", filePath, err)
	}

	numRuns, err := generateRuns(filePath, maxLines)
	if err != nil {
		return err
	}

	if numRuns == 0 {
		return nil
	}

	err = mergeRuns(numRuns, filePath)
	if err != nil {
		return err
	}
	return nil
}

//sort file by reading full of it in memory
func SortFileInMemory(filePath string) error {
	filePath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("Can't get abs path of %s error %v", filePath, err)
	}

	f, err := os.OpenFile(filePath, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	maxLines := (1 << 31) - 1
	lines, err := readLines(bufio.NewReader(f), maxLines)
	if err != nil {
		return err
	}

	sort.Strings(lines)

	err = f.Truncate(0)
	if err != nil {
		return fmt.Errorf("Can't truncate file %s error %v", filePath, err)
	}

	for _, line := range lines {
		err = writeLine(f, line)
		if err != nil {
			return err
		}
	}

	return nil
}

//check if file contains sorted lines
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

	fr := bufio.NewReader(f)
	prev := ""
	var pos int64
	for {
		s, eof, err := readLine(fr)
		if err != nil {
			return false, err
		}
		if eof {
			break
		}

		if pos > 0 {
			if strings.Compare(prev, s) > 0 {
				return false, nil
			}
		}
		prev = s
		pos++
	}

	return true, nil
}
