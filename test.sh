#!/bin/bash -x
rm -f filesort.bin
go build -o filesort.bin main.go
FILE=./test-file
ORIG_FILE=./test-file-orig
rm -f $FILE
rm -f $ORIG_FILE
./filesort.bin -cmd generate -filePath $FILE -numLines 4 -lineLen 5
cp $FILE $ORIG_FILE
./filesort.bin -cmd sort -filePath $FILE -maxLines 2
./filesort.bin -cmd check -filePath $FILE
