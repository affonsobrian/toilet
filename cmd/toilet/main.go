package main

import (
	toilet "github.com/affonsobrian/toilet/internal"
)

func main() {
	input, pathSize := toilet.ParseInput()
	data, total, errors := toilet.CalculateValues(input)
	toilet.ShowData(input, data, errors, total, pathSize)
}
