package toilet

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"unicode"
	"unicode/utf8"
)

func formatOutput(input *Input, output *Output) string {
	outputString := ""
	if input.Lines {
		outputString = fmt.Sprintf("%s %d", outputString, output.LineCount)
	}
	if input.Words {
		outputString = fmt.Sprintf("%s %d", outputString, output.CharCount)
	}
	if input.Bytes {
		outputString = fmt.Sprintf("%s %d", outputString, output.ByteCount)
	}
	outputString = fmt.Sprintf("%s %s", outputString, output.Filepath)
	return outputString
}

func ShowData(input *Input, data *[]Output, errors *[]error, total *Output, pathSize int) {
	// Show values
	for _, fd := range *data {
		fmt.Println(formatOutput(input, &fd))
	}
	// Show failed values
	for _, err := range *errors {
		fmt.Println(err.Error())
	}
	// Show total
	if pathSize > 1 {
		fmt.Println(formatOutput(input, total))
	}
}

func getFileByteSize(fileInfo *fs.FileInfo) int64 {
	return (*fileInfo).Size()
}

func isSpaceOrBreakLine(b byte) bool {
	bb, _ := utf8.DecodeRune([]byte{b})
	return unicode.IsSpace(bb)
}

func wordCount(s []byte, lByte *byte) int {
	c := 0
	for _, cByte := range s {
		if isSpaceOrBreakLine(cByte) && !isSpaceOrBreakLine(*lByte) {
			c += 1
		}
		*lByte = cByte
	}
	return c
}

func getFileLineSize(r io.Reader) (int64, int64, error) {
	buf := make([]byte, 32*1024)
	var count int64 = 0
	var countWords int64 = 0
	breaklineByte := []byte{'\n'}
	lByte := byte(' ')
	for {
		n, err := r.Read(buf)
		count += int64(bytes.Count(buf[:n], breaklineByte))
		countWords += int64(wordCount(buf[:n], &lByte))
		if err == io.EOF {
			break
		}
		if err != nil {
			return count, countWords, err
		}
	}
	if !isSpaceOrBreakLine(lByte) {
		countWords += 1
	}
	return count, countWords, nil
}

func CalculateValues(input *Input) (*[]Output, *Output, *[]error) {
	var datas []Output
	var total Output = Output{Filepath: "total"}
	var errors []error

	for _, filepath := range input.FilePaths {
		file, err := os.Open(filepath)
		if err != nil {
			errors = append(errors, fmt.Errorf("wc: %s: No such file or directory", filepath))
			continue
		}
		defer file.Close()

		fd := Output{Filepath: filepath}
		if input.Bytes {
			fStat, _ := file.Stat()
			fd.ByteCount = getFileByteSize(&fStat)
			total.ByteCount += fd.ByteCount
		}

		if input.Lines || input.Words {
			fd.LineCount, fd.CharCount, _ = getFileLineSize(file)
			total.LineCount += fd.LineCount
			total.CharCount += fd.CharCount
		}

		datas = append(datas, fd)
	}

	return &datas, &total, &errors
}
