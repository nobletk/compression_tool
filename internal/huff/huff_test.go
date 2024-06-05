package huff

import (
	"testing"
)

func TestValidBuildTree(t *testing.T) {
	tests := []struct {
		name                string
		input               map[rune]int
		expectedFrequencies []int
		expectedPrefix      string
	}{
		{
			name: "case 1",
			input: map[rune]int{
				'c': 32,
				'd': 42,
				'e': 120,
				'k': 7,
				'l': 42,
				'm': 24,
				'u': 37,
				'z': 2,
			},
			expectedFrequencies: []int{306, 120, 186, 79, 37, 42, 107, 42, 65, 32, 33, 9, 2, 7, 24},
			expectedPrefix:      "prefixTable(8):\n(c -- 99: 1110)\n(d -- 100: 101)\n(e -- 101: 0)\n(k -- 107: 111101)\n(l -- 108: 110)\n(m -- 109: 11111)\n(u -- 117: 100)\n(z -- 122: 111100)\n",
		},
		{
			name: "case 2",
			input: map[rune]int{
				' ': 3,
				'a': 4,
				'b': 3,
				'c': 2,
				'd': 1,
			},
			expectedFrequencies: []int{13, 6, 3, 1, 2, 3, 7, 3, 4},
			expectedPrefix:      "prefixTable(5):\n(  -- 32: 01)\n(a -- 97: 11)\n(b -- 98: 10)\n(c -- 99: 001)\n(d -- 100: 000)\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := tt.input
			root, prefixTable, err := buildTree(input)
			if root == nil {
				t.Fatal("root node is nil")
			}
			if err != nil {
				t.Fatal(err.Error())
			}

			actualFreq := getTreeFrequency(root)
			if len(actualFreq) != len(tt.expectedFrequencies) {
				t.Fatalf("len(expectedFrequencies)=%d, len(actualFreq)=%d",
					len(tt.expectedFrequencies), len(actualFreq))
			}

			for i, node := range actualFreq {
				if node.Count != tt.expectedFrequencies[i] {
					t.Errorf("node.Frequency not equal expected at index %d. got=%d, want=%d",
						i, node.Count, tt.expectedFrequencies[i])
				}
			}

			// printTree(root, 0) //printing for debuging

			actualPrefix := printPrefixTable(prefixTable)

			assertEqual(t, actualPrefix, tt.expectedPrefix)
		})
	}
}

func TestInvalidBuildTree(t *testing.T) {
	tests := []struct {
		name        string
		freqMap     map[rune]int
		expectedErr string
	}{
		{"Empty frequency map", map[rune]int{}, "frequency map is empty"},
		{"Zero frequency", map[rune]int{'a': 0}, "frequency must be greater than zero"},
		{"Single node", map[rune]int{'a': 1}, "priority queue must contain at least two nodes"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := buildTree(tt.freqMap)
			if err == nil || err.Error() != tt.expectedErr {
				t.Errorf("buildTree() error got= %v, want=%v", err, tt.expectedErr)
			}
		})
	}
}

func TestEncodeTree(t *testing.T) {
	root := &huffmanNode{
		Char:  0,
		Count: 306,
		Left:  &huffmanNode{Char: 'e', Count: 120},
		Right: &huffmanNode{
			Char:  0,
			Count: 120,
			Left: &huffmanNode{
				Char:  0,
				Count: 79,
				Left:  &huffmanNode{Char: 'u', Count: 37},
				Right: &huffmanNode{Char: 'd', Count: 42},
			},
			Right: &huffmanNode{
				Char:  0,
				Count: 107,
				Left:  &huffmanNode{Char: 'l', Count: 42},
				Right: &huffmanNode{
					Char:  0,
					Count: 65,
					Left:  &huffmanNode{Char: 'c', Count: 32},
					Right: &huffmanNode{
						Char:  0,
						Count: 33,
						Left: &huffmanNode{
							Char:  0,
							Count: 9,
							Left:  &huffmanNode{Char: 'z', Count: 2},
							Right: &huffmanNode{Char: 'k', Count: 7},
						},
						Right: &huffmanNode{Char: 'm', Count: 24},
					},
				},
			},
		},
	}

	expectedPrefix := map[rune]string{
		'e': "0",
		'u': "100",
		'd': "101",
		'l': "110",
		'c': "1110",
		'z': "111100",
		'k': "111101",
		'm': "11111",
	}

	prefixTable := make(prefixTable, 0)
	encodeTree(root, "", prefixTable)

	if len(prefixTable) != len(expectedPrefix) {
		t.Fatalf("len(prefixTable)=%d, len(expectedPrefix)=%d",
			len(prefixTable), len(expectedPrefix))
	}

	actualPrefix := printPrefixTable(prefixTable)
	expectedSortedPrefix := printPrefixTable(expectedPrefix)

	assertEqual(t, actualPrefix, expectedSortedPrefix)
}
