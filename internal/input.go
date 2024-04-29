package toilet

import (
	"flag"
)

type Input struct {
	Bytes     bool
	Words     bool
	Lines     bool
	FilePaths []string
}

func ParseInput() (*Input, int) {
	input := new(Input)
	flag.BoolVar(&input.Bytes, "c", false, "print the byte counts")
	flag.BoolVar(&input.Bytes, "bytes", false, "print the byte counts")
	flag.BoolVar(&input.Words, "w", false, "print the word counts")
	flag.BoolVar(&input.Words, "chars", false, "print the word counts")
	flag.BoolVar(&input.Lines, "l", false, "print the line counts")
	flag.BoolVar(&input.Lines, "lines", false, "print the line counts")

	flag.Parse()

	if !(input.Bytes || input.Lines || input.Words) {
		input.Bytes = true
		input.Lines = true
		input.Words = true
	}

	input.FilePaths = flag.Args()
	return input, len(input.FilePaths)
}
