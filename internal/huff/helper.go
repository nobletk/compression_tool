package huff

import (
	"bytes"
	"cmp"
	"fmt"
	"slices"
	"sort"
	"strings"
	"testing"
)

func printTree(node *huffmanNode, level int) {
	if node == nil {
		return
	}

	printTree(node.Right, level+1)
	for i := 0; i < level; i++ {
		fmt.Print("\t")
	}

	if node.Left == nil && node.Right == nil {
		fmt.Printf("(%c -- %v:%d,0b%s)\n", node.Char, node.Char, node.Count, node.Code)
	} else {
		fmt.Printf("(:%d)=>\n", node.Count)
	}

	printTree(node.Left, level+1)
}

func printPrefixTable(prefix map[rune]string) string {
	type pair struct {
		key   rune
		value string
	}

	var pairs []pair
	for key, value := range prefix {
		pairs = append(pairs, pair{key, value})
	}

	slices.SortFunc(pairs, func(a, b pair) int { return cmp.Compare(a.key, b.key) })

	var sortedKeys []rune
	for _, pair := range pairs {
		sortedKeys = append(sortedKeys, pair.key)
	}

	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("prefixTable(%d):\n", len(prefix)))

	var prefixSorted []string
	for _, key := range sortedKeys {
		prefixSorted = append(prefixSorted, fmt.Sprintf("(%c -- %v: %s)\n", key, key, prefix[key]))
	}

	out.WriteString(strings.Join(prefixSorted, ""))

	return out.String()
}

func printSortedMap(m map[rune]int) string {
	var out bytes.Buffer

	keys := make([]rune, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}

	slices.SortFunc(keys, func(a, b rune) int { return cmp.Compare(a, b) })

	for _, k := range keys {
		out.WriteString(fmt.Sprintf("(%c: %d)\n", k, m[k]))
	}

	return out.String()
}

func printLookup(lookup []byteLookup) string {
	slices.SortFunc(lookup, func(a, b byteLookup) int { return cmp.Compare(a.char, b.char) })

	var buff bytes.Buffer

	buff.WriteString(fmt.Sprintf("lookup(%d):\n", len(lookup)))

	var sorted []string
	for i := range lookup {
		sorted = append(sorted, fmt.Sprintf("(%c -- %v : %s)\n", lookup[i].char,
			lookup[i].char, lookup[i].code))
	}

	buff.WriteString(strings.Join(sorted, ""))

	return buff.String()
}

func byteSlicesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func sortRunes(pt map[rune]string) []rune {
	runes := make([]rune, 0, len(pt))
	for r := range pt {
		runes = append(runes, r)
	}

	sort.Slice(runes, func(a, b int) bool { return runes[a] < runes[b] })

	return runes
}

func assertEqual[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if actual != expected {
		t.Errorf("got= %v;\nwant= %v", actual, expected)
	}
}

func assertEqualBytes(t *testing.T, actual, expected []byte) {
	t.Helper()

	if len(actual) != len(expected) {
		t.Fatalf("length mismatch: got= %v, want= %v", len(actual), len(expected))
	}

	for i := range actual {
		if actual[i] != expected[i] {
			t.Errorf("index %d: got= %v, want= %v", i, actual[i], expected[i])
			return
		}
	}
}
