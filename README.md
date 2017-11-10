### filesort Golang package to sort files

#### Project structure:
main.go - command line utility
filesort/filesort.go - core code to sort files with memory limit
filesort/filesort_test.go - tests

#### Usage:
go build -o filesort.bin main.go - to create command line utility which
can generate/sort/check/remove files
./filesort.bin -cmd generate -filePath ./file-to-sort -numLines 4 -lineLen 5
./filesort.bin -cmd sort -filePath ./file-to-sort -maxLines 2
./filesort.bin -cmd check - filePath ./file-to-sort

#### Run tests:
cd file_sort && go test

#### File sorting algorithm:
Split input file into several parts(temporary run files) and sort them in memory
Merge temporary run files into one last run file
Rename last run file into output file(which is input usually)
