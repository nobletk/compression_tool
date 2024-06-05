package huff

import (
	"bytes"
	"errors"
	"unicode/utf8"
)

func Decompress(data []byte) ([]byte, error) {
	if len(data) < 2 {
		return nil, errors.New("data is too short")
	}
	bitLen := int(data[0])
	if bitLen < 1 || bitLen > 8 {
		return nil, errors.New("invalid bit length")
	}
	data = data[1:]

	treeBytes, err := getTreeBytes(data)
	if err != nil {
		return nil, err
	}

	encData, err := getEncData(data)
	if err != nil {
		return nil, err
	}

	treeRoot, err := rebuildEncTree(treeBytes)
	if err != nil {
		return nil, err
	}

	prefix := make(map[rune]string, 0)
	encodeTree(treeRoot, "", prefix)

	lookup := getEncLookup(prefix)

	decompressed, err := decode(encData, lookup, bitLen)
	if err != nil {
		return nil, err
	}

	return decompressed, nil
}

func getTreeBytes(b []byte) ([]byte, error) {
	treeStart := byte(31)
	treeEnd := []byte{37, 37}

	startIdx := bytes.IndexByte(b, treeStart)
	if startIdx == -1 {
		return nil, errors.New("start of tree header not found")
	}
	startIdx++

	endIdx := bytes.Index(b[startIdx:], treeEnd)
	if endIdx == -1 {
		return nil, errors.New("end of tree header not found")
	}
	endIdx += startIdx

	if startIdx >= len(b) || endIdx > len(b) || startIdx > endIdx {
		return nil, errors.New("slice bounds out of range")
	}

	return b[startIdx:endIdx], nil
}

func getEncData(b []byte) ([]byte, error) {
	treeEnd := []byte{37, 37}
	startIdx := bytes.Index(b, treeEnd)
	if startIdx == -1 {
		return nil, errors.New("end of tree delimiter not found")
	}
	startIdx += len(treeEnd)

	if startIdx >= len(b) {
		return nil, errors.New("start of encoded data not found")
	}

	return b[startIdx:], nil
}

type byteLookup struct {
	code []byte
	char rune
}

func getEncLookup(pt map[rune]string) []byteLookup {
	lookup := make([]byteLookup, 0, len(pt))

	for char, code := range pt {
		bl := byteLookup{
			code: []byte(code),
			char: char,
		}
		lookup = append(lookup, bl)
	}

	return lookup
}

func decode(enc []byte, lookup []byteLookup, validBits int) ([]byte, error) {
	decompressed := make([]byte, 0)
	currentBits := make([]byte, 0)
	totalBits := (len(enc)-1)*8 + validBits

	bitIndex := 0

	for i := 0; i < len(enc); i++ {
		for bit := 7; bit >= 0; bit-- {
			if bitIndex >= totalBits {
				break
			}

			currentBit := ((enc[i] >> bit) & 1) + '0'
			currentBits = append(currentBits, byte(currentBit))

			for _, entry := range lookup {
				if byteSlicesEqual(currentBits, entry.code) {
					decompressed = utf8.AppendRune(decompressed, entry.char)
					currentBits = make([]byte, 0)
					break
				}
			}
			bitIndex++
		}
	}

	return decompressed, nil
}

func newInternalNode() *huffmanNode {
	return &huffmanNode{
		Char:  0,
		Count: 0,
		Code:  "",
		Left:  nil,
		Right: nil,
	}
}

func newLeafNode(char rune) *huffmanNode {
	return &huffmanNode{
		Char:  char,
		Count: 0,
		Code:  "",
		Left:  nil,
		Right: nil,
	}
}

func rebuildEncTree(treeBytes []byte) (*huffmanNode, error) {
	if len(treeBytes) == 0 {
		return nil, errors.New("tree bytes are empty")
	}
	idx := 0
	return rebuildEncSubTrees(treeBytes, &idx)
}

func rebuildEncSubTrees(treeBytes []byte, idx *int) (*huffmanNode, error) {
	if *idx >= len(treeBytes) {
		return nil, errors.New("index out if range")
	}

	if treeBytes[*idx] == 0 {
		node := newInternalNode()
		*idx++
		left, err := rebuildEncSubTrees(treeBytes, idx)
		if err != nil {
			return nil, err
		}
		node.Left = left
		right, err := rebuildEncSubTrees(treeBytes, idx)
		if err != nil {
			return nil, err
		}
		node.Right = right

		return node, nil
	} else if treeBytes[*idx] == 1 {
		*idx++
		if *idx >= len(treeBytes) {
			return nil, errors.New("index out of range after leaf indicator")
		}

		char, sz := utf8.DecodeRune(treeBytes[*idx:])
		if char == utf8.RuneError {
			return nil, errors.New("invalid UTF-8 in tree bytes")
		}
		node := newLeafNode(char)
		*idx += sz

		return node, nil
	}

	return nil, errors.New("unexpected value in tree bytes")
}
