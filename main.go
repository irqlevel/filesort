package main

import (
	"./filesort"
	"flag"
	"log"
	"os"
)

//command line utility to generate/sort/check/remove files
func do_job() int {
	var filePath string
	var cmd string
	var numLines int64
	var lineLen int
	var memory bool
	var maxLines int

	flag.StringVar(&filePath, "filePath", "", "file to sort")
	flag.StringVar(&cmd, "cmd", "", "command to execute")
	flag.Int64Var(&numLines, "numLines", -1, "number of file lines")
	flag.IntVar(&lineLen, "lineLen", -1, "file line length")
	flag.BoolVar(&memory, "memory", false, "sort file in memory")
	flag.IntVar(&maxLines, "maxLines", -1, "max file lines to use in sort")

	flag.Parse()

	log.SetOutput(os.Stdout)
	log.Printf("cmd %s filePath %s memory %t numLines %d lineLen %d maxLines\n",
		cmd, filePath, memory, numLines, lineLen, maxLines)

	switch cmd {
	case "sort":
		var err error
		if memory {
			err = filesort.SortFileInMemory(filePath)
		} else {
			err = filesort.SortFile(filePath, maxLines)
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
	os.Exit(do_job())
}
