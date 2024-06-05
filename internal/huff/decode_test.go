package huff

import (
	"testing"
	"unicode/utf8"
)

func TestDecompress(t *testing.T) {
	input := []byte{5, 31, 0, 0, 0, 1, 100, 1, 99, 1, 32, 0, 1, 98, 1, 97, 37, 37, 255, 106, 73, 64}

	decompressed, err := Decompress(input)
	if err != nil {
		t.Fatalf(err.Error())
	}

	actualDecompText := string(decompressed)
	expectedDecompText := `aaaa bbb cc d`

	assertEqual(t, actualDecompText, expectedDecompText)
}

func TestGetEncLookup(t *testing.T) {
	input := map[rune]string{
		' ': "01",
		'a': "11",
		'b': "10",
		'c': "001",
		'…': "000",
	}

	lookup := getEncLookup(input)
	actualSorted := printLookup(lookup)

	expectedLookup := []byteLookup{
		{code: []byte{48, 49}, char: ' '},     // "01" -> ' '
		{code: []byte{49, 48}, char: 'b'},     // "10" -> 'b'
		{code: []byte{49, 49}, char: 'a'},     // "11" -> 'a'
		{code: []byte{48, 48, 48}, char: '…'}, // "000" -> '…'
		{code: []byte{48, 48, 49}, char: 'c'}, // "001" -> 'c'
	}
	expectedSorted := printLookup(expectedLookup)

	assertEqual(t, actualSorted, expectedSorted)
}

func TestDecode(t *testing.T) {
	tests := []struct {
		enc         []byte
		bits        int
		lookup      []byteLookup
		expectedDec []byte
	}{
		{
			enc:  []byte{255, 106, 73, 64},
			bits: 5,
			lookup: []byteLookup{
				{code: []byte{48, 48, 48}, char: '…'}, // "000" -> '…'
				{code: []byte{48, 48, 49}, char: 'c'}, // "001" -> 'c'
				{code: []byte{48, 49}, char: ' '},     // "01" -> ' '
				{code: []byte{49, 48}, char: 'b'},     // "10" -> 'b'
				{code: []byte{49, 49}, char: 'a'},     // "11" -> 'a'
			},

			expectedDec: utf8.AppendRune([]byte{
				'a', 'a', 'a', 'a', ' ',
				'b', 'b', 'b', ' ',
				'c', 'c', ' ',
			}, '…'),
		},
	}

	for _, tt := range tests {
		t.Helper()

		actualDecode, err := decode(tt.enc, tt.lookup, tt.bits)
		if err != nil {
			t.Fatal(err.Error())
		}
		assertEqualBytes(t, actualDecode, tt.expectedDec)
	}
}

func TestGetTreeBytes(t *testing.T) {
	input := []byte{5, 31,
		0, 0, 0, 1, 100, 1, 99, 1, 32, 0, 1, 98, 1, 97, 37, 37,
		255, 106, 73, 64}

	actualTreeBytes, err := getTreeBytes(input)
	if err != nil {
		t.Fatal(err.Error())
	}

	expectedTreeBytes := []byte{0, 0, 0, 1, 100, 1, 99, 1, 32, 0, 1, 98, 1, 97}

	assertEqualBytes(t, actualTreeBytes, expectedTreeBytes)
}

func TestRebuildEncTree(t *testing.T) {
	input := []byte{0, 0, 0, 0, 1, 'd', 1, 'c', 1, ' ',
		1, 'b', 0, 1, 'a', 0, 1, 195, 170, 0, 1, 'f', 1, 'g',
	}

	root, err := rebuildEncTree(input)
	if err != nil {
		t.Fatal(err.Error())
	}
	prefixTable := make(prefixTable, 0)
	encodeTree(root, "", prefixTable)
	// printTree(root, 0)

	expectedPrefix := map[rune]string{
		'd': "0000",
		'c': "0001",
		' ': "001",
		'b': "01",
		'a': "10",
		'ê': "110",
		'f': "1110",
		'g': "1111",
	}

	if len(prefixTable) != len(expectedPrefix) {
		t.Fatalf("len(prefixTable)=%d; len(expectedPrefix)=%d", len(prefixTable),
			len(expectedPrefix))
	}

	actualPrefix := printPrefixTable(prefixTable)
	expectedSortedPrefix := printPrefixTable(expectedPrefix)
	assertEqual(t, actualPrefix, expectedSortedPrefix)
}

func TestGetEncData(t *testing.T) {
	input := []byte{5, 31, 0, 0, 0, 1, 100, 1, 99, 1, 32, 0, 1, 98, 1, 97, 37, 37, 255, 106, 73, 64}

	actualEncData, err := getEncData(input)
	if err != nil {
		t.Fatal(err.Error())
	}

	expectedEncData := []byte{255, 106, 73, 64}

	assertEqualBytes(t, actualEncData, expectedEncData)
}

func TestDecompressInvalidInput(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		expectedErr string
	}{
		{"Short input", []byte{0}, "data is too short"},
		{"Invalid bit length", []byte{0, 31, 37, 37, 0}, "invalid bit length"},
		{"No tree start", []byte{8, 1, 2, 3, 4}, "start of tree header not found"},
		{"No tree end", []byte{8, 31, 1, 2, 3, 4}, "end of tree header not found"},
		{"No encoded data start", []byte{8, 31, 1, 2, 37, 37}, "start of encoded data not found"},
		{"Empty tree bytes", []byte{8, 31, 37, 37, 37, 37}, "tree bytes are empty"},
		{"Invalid UTF-8 in tree", []byte{8, 31, 1, 0xff, 37, 37, 37, 37}, "invalid UTF-8 in tree bytes"},
		{"Empty encoded data", []byte{8, 31, 1, 'a', 37, 37}, "start of encoded data not found"},
	}

	for _, tt := range tests {
		t.Helper()
		t.Run(tt.name, func(t *testing.T) {
			_, err := Decompress(tt.input)
			if err == nil || err.Error() != tt.expectedErr {
				t.Errorf("Decompress() error got=%v, want=%v", err, tt.expectedErr)
			}
		})
	}
}
