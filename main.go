package main

import (
	"os"
	"./filesort"
	"log"
	"flag"
)

func do() int {
	var filePath string
	var cmd string
	var numLines int64
	var lineLen int
	var memory bool

	flag.StringVar(&filePath, "filePath", "", "file to sort")
	flag.StringVar(&cmd, "cmd", "", "command to execute")
	flag.Int64Var(&numLines, "numLines", -1, "number of file lines")
	flag.IntVar(&lineLen, "lineLen", -1, "file line length")
	flag.BoolVar(&memory, "memory", false, "sort file in memory")

	flag.Parse()

	log.SetOutput(os.Stdout)
	switch cmd {
	case "sort":
		var err error
		if memory {
			err = filesort.SortFileInMemory(filePath)
		} else {
			err = filesort.SortFile(filePath)
		}

		if err != nil {
			log.Printf("sort file %s error %v", filePath, err)
			return -1
		}
	case "generate":
		err := filesort.GenerateFile(filePath, numLines, lineLen)
		if err != nil {
			log.Printf("generate file %s error %v", filePath, err)
			return -1
		}
	case "remove":
		err := filesort.RemoveFile(filePath)
		if err != nil {
			log.Printf("remove file %s error %v", filePath, err)
			return -1
		}
	case "check":
		sorted, err := filesort.IsFileSorted(filePath)
		if err != nil {
			log.Printf("file %s check sorted error %v", filePath, err)
			return -1
		}
		if !sorted {
			log.Printf("file %s is not sorted", filePath)
			return -1
		}
	default:
		log.Printf("unknown command %s", cmd)
		return -1
	}
	return 0
}

func main() {
	os.Exit(do())
}
