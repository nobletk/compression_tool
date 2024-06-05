package huff

import (
	"bytes"
	"errors"
	"unicode/utf8"
)

type FrequencyMap map[rune]int

func Compress(input []byte) ([]byte, error) {
	runesFreq, err := getRunesFrequency(input)
	if err != nil {
		return nil, err
	}

	treeRoot, prefixTable, err := buildTree(runesFreq)
	if err != nil {
		return nil, err
	}

	treeBuff, err := serializeTree(treeRoot)
	if err != nil {
		return nil, err
	}

	bitBuff, totalBits, err := encData(input, prefixTable)
	if err != nil {
		return nil, err
	}

	var outBuff bytes.Buffer
	validBits := byte(totalBits % 8)
	if validBits == 0 {
		validBits = 8
	}
	outBuff.WriteByte(validBits)
	outBuff.WriteByte(byte(0x1F)) //unit separator
	outBuff.Write(treeBuff.Bytes())
	outBuff.WriteString("%%")
	outBuff.Write(bitBuff.Bytes())

	return outBuff.Bytes(), nil
}

func getRunesFrequency(input []byte) (FrequencyMap, error) {
	freqMap := make(FrequencyMap, 0)

	for i := 0; i < len(input); {
		char, sz := utf8.DecodeRune(input[i:])
		if char == utf8.RuneError {
			return nil, errors.New("invalid UTF-8 encoding")
		}
		freqMap[char]++
		i += sz
	}

	return freqMap, nil
}

func serializeTree(node *huffmanNode) (bytes.Buffer, error) {
	var buff bytes.Buffer
	err := serializeNode(node, &buff)
	if err != nil {
		return bytes.Buffer{}, err
	}
	return buff, nil
}

func serializeNode(node *huffmanNode, buff *bytes.Buffer) error {
	if node == nil {
		return errors.New("node is nil")
	}
	if node.Left == nil && node.Right == nil {
		buff.WriteByte(1)
		if _, err := buff.WriteRune(node.Char); err != nil {
			return err
		}
	} else {
		buff.WriteByte(0)
		if err := serializeNode(node.Left, buff); err != nil {
			return err
		}
		if err := serializeNode(node.Right, buff); err != nil {
			return err
		}
	}
	return nil
}

func encData(data []byte, preTab map[rune]string) (bytes.Buffer, int, error) {
	var bitBuff bytes.Buffer
	var currentByte byte
	bitCount := 0
	totalBits := 0

	for i := 0; i < len(data); {
		char, sz := utf8.DecodeRune(data[i:])
		if char == utf8.RuneError {
			return bytes.Buffer{}, 0, errors.New("invalid UTF-8 encoding")
		}
		i += sz
		code, exists := preTab[char]
		if !exists {
			return bytes.Buffer{}, 0, errors.New("char not found in prefix table")
		}
		for _, bit := range code {
			if bit == '1' {
				currentByte |= (1 << (7 - bitCount))
			}
			bitCount++
			totalBits++
			if bitCount == 8 {
				bitBuff.WriteByte(currentByte)
				currentByte = 0
				bitCount = 0
			}
		}
	}

	if bitCount > 0 {
		bitBuff.WriteByte(currentByte)
	}

	return bitBuff, totalBits, nil
}
