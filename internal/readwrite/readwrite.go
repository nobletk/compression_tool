package readwrite

import (
	"bufio"
	"errors"
	"os"
)

func ReadFile(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return data, nil
}

func WriteFile(fileName string, data []byte) error {
	file, err := os.Create(fileName)
	if err != nil {
		return errors.New(err.Error())
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.Write(data)
	if err != nil {
		return errors.New(err.Error())
	}

	err = writer.Flush()
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}
