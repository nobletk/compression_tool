package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"compression_tool.nobletk/internal/huff"
	"compression_tool.nobletk/internal/readwrite"
)

func main() {
	f := flags{}
	pf := f.parseFlags()

	if f.inputFlag == "" || f.outputFlag == "" {
		fmt.Println("missing argument")
		flag.Usage()
		os.Exit(1)
	}

	if pf.compFlag {
		data, err := readwrite.ReadFile(*&pf.inputFlag)
		if err != nil {
			panic(err)
		}

		compData, err := huff.Compress(data)
		if err != nil {
			panic(err)
		}

		outputPath := filepath.Join(filepath.Dir(*&pf.inputFlag), *&pf.outputFlag)

		err = readwrite.WriteFile(outputPath, compData)
		if err != nil {
			panic(err)
		}

		fmt.Printf("File compressed successfully %s\n", outputPath)
		os.Exit(0)
	}

	if pf.decompFlag {
		data, err := readwrite.ReadFile(*&pf.inputFlag)
		if err != nil {
			panic(err)
		}

		decompData, err := huff.Decompress(data)
		if err != nil {
			panic(err)
		}

		outputPath := filepath.Join(filepath.Dir(*&pf.inputFlag), *&pf.outputFlag)

		err = readwrite.WriteFile(outputPath, decompData)
		if err != nil {
			panic(err)
		}

		fmt.Printf("File decompressed successfully %s\n", outputPath)
		os.Exit(0)
	}

	fmt.Println("Incorrect argument!")
	flag.Usage()
	os.Exit(1)
}

type flags struct {
	compFlag   bool
	decompFlag bool
	inputFlag  string
	outputFlag string
}

func (f *flags) parseFlags() flags {
	flag.StringVar(&f.inputFlag, "i", "", "Input file path to compress")
	flag.StringVar(&f.outputFlag, "o", "", "Output file path to compress")
	flag.BoolVar(&f.compFlag, "c", false, "Compress the input file to the output file")
	flag.BoolVar(&f.decompFlag, "d", false, "Decompress the input file to the output file")

	flag.Parse()

	return *f
}
