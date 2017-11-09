#!/bin/bash -xv
rm -f filesort.bin
go build -o filesort.bin main.go
FILE=./test-file
./filesort.bin -cmd generate -filePath $FILE -numLines 10 -lineLen 5
./filesort.bin -cmd check -filePath $FILE
./filesort.bin -cmd remove -filePath $FILE
