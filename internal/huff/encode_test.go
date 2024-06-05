package huff

import (
	"testing"
)

func TestGetRunesFrequency(t *testing.T) {
	input := []byte{'a', 'a', 'a', 'a', ' ', 'b', 'b', 'b', ' ', 'c', 'c', ' ', '-'}

	frequencyMap, err := getRunesFrequency(input)
	if err != nil {
		t.Fatalf(err.Error())
	}

	actualFreqMap := printSortedMap(frequencyMap)
	expectedFreqMap := "( : 3)\n(-: 1)\n(a: 4)\n(b: 3)\n(c: 2)\n"

	assertEqual(t, actualFreqMap, expectedFreqMap)
}

func TestEncodeData(t *testing.T) {
	tests := []struct {
		input             []byte
		prefixTable       map[rune]string
		expectedBuff      []byte
		expectedTotalBits int
	}{
		{
			input: []byte{'a', 'a', 'a', 'a', ' ', 'b', 'b', 'b', ' ', 'c', 'c', ' ', '-'},
			prefixTable: map[rune]string{
				' ': "01",
				'a': "11",
				'b': "10",
				'c': "001",
				'-': "000",
			},
			expectedBuff:      []byte{255, 106, 73, 64},
			expectedTotalBits: 29,
		},
	}

	for _, tt := range tests {
		t.Helper()

		buff, actualTotalBits, err := encData(tt.input, tt.prefixTable)
		if err != nil {
			t.Fatalf(err.Error())
		}
		actualBuff := buff.Bytes()

		assertEqualBytes(t, actualBuff, tt.expectedBuff)
		assertEqual(t, actualTotalBits, tt.expectedTotalBits)
	}
}

func TestCompress(t *testing.T) {
	input := []byte{'a', 'a', 'a', 'a', ' ', 'b', 'b', 'b', ' ', 'c', 'c', ' ', 'd'}

	actualEncodedText, err := Compress(input)
	if err != nil {
		t.Fatalf(err.Error())
	}

	expectedEncodedText := []byte{5, 31, 0, 0, 0, 1, 100, 1, 99, 1, 32, 0, 1, 98, 1,
		97, 37, 37, 255, 106, 73, 64}

	assertEqualBytes(t, actualEncodedText, expectedEncodedText)
}

func TestCompressInvalidUTF8(t *testing.T) {
	invalidUTF8 := []byte{255, 254, 253}

	_, err := Compress(invalidUTF8)
	if err == nil {
		t.Fatal("expected error for invalid UTF-8 input, got nil")
	}

	expectedErr := "invalid UTF-8 encoding"
	if err.Error() != expectedErr {
		t.Fatalf("expected error: %s, got: %s", expectedErr, err.Error())
	}
}

func TestSerializTreeNilNode(t *testing.T) {
	_, err := serializeTree(nil)
	if err == nil {
		t.Fatal("expected error for nil node, got nil")
	}

	expectedErr := "node is nil"
	if err.Error() != expectedErr {
		t.Fatalf("expected error: %s, got: %s", expectedErr, err.Error())
	}
}

func TestEndDataMissingCharInPrefixTable(t *testing.T) {
	validUTF8 := []byte{'a', 'b', 'c'}
	preTab := map[rune]string{
		'a': "00",
		'b': "01",
	}

	_, _, err := encData(validUTF8, preTab)
	if err == nil {
		t.Fatal("expected error for missing char in prefix table, got nil")
	}

	expectedErr := "char not found in prefix table"
	if err.Error() != expectedErr {
		t.Fatalf("expected error: %s, got: %s", expectedErr, err.Error())
	}

}
