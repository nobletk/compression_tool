package test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestCompareFiles(t *testing.T) {
	file1 := "./testdata/test.txt"
	file2 := "./testdata/decomp.txt"

	_, _, err := compareFiles(file1, file2)
	if err != nil {
		t.Fatalf("Error comparing files:%v", err)
	}

}

func compareFiles(file1, file2 string) (bool, int, error) {
	f1, err := os.Open(file1)
	if err != nil {
		return false, -1, errors.New(fmt.Sprintf("failed to open file1: %v", err))
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		return false, -1, errors.New(fmt.Sprintf("failed to open file2: %v", err))
	}
	defer f2.Close()

	if err := skipBOM(f1); err != nil {
		return false, -1, errors.New(fmt.Sprintf("failed to skip BOM in file1: %v", err))
	}
	if err := skipBOM(f2); err != nil {
		return false, -1, errors.New(fmt.Sprintf("failed to skip BOM in file2: %v", err))
	}

	const buffSz = 4096
	buff1 := make([]byte, buffSz)
	buff2 := make([]byte, buffSz)

	pos := 0
	for {
		r1, err1 := f1.Read(buff1)
		r2, err2 := f2.Read(buff2)

		if r1 != r2 || !bytes.Equal(buff1[:r1], buff2[:r2]) {
			for i := 0; i < r1 && i < r2; i++ {
				if buff1[i] != buff2[i] {
					return false, pos + i,
						errors.New(fmt.Sprintf("files differ at (original: %v) (decomp: %v)-- pos %d",
							buff1[i], buff2[i], pos+i))
				}
			}
			if r1 != r2 {
				return false, pos + r1, nil
			}
		}

		if err1 == io.EOF && err2 == io.EOF {
			break
		} else if err1 != nil {
			return false, pos, errors.New(fmt.Sprintf("error reading file1: %v", pos))
		} else if err2 != nil {
			return false, pos, errors.New(fmt.Sprintf("error reading file2: %v", pos))
		}

		pos += r1
	}

	return true, 0, nil
}

func skipBOM(f *os.File) error {
	utf8BOM := []byte{239, 187, 191}
	buff := make([]byte, len(utf8BOM))

	_, err := f.Read(buff)
	if err != nil && err != io.EOF {
		return err
	}
	if bytes.Equal(buff, utf8BOM) {
		return nil
	}

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	return nil
}
