package huff

import (
	"container/heap"
	"errors"
)

type huffmanNode struct {
	Char  rune
	Count int
	Code  string
	Left  *huffmanNode
	Right *huffmanNode
}

type priorityQueue []*huffmanNode

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	if pq[i].Count == pq[j].Count {
		return pq[i].Char < pq[j].Char
	}
	return pq[i].Count < pq[j].Count
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *priorityQueue) Push(x any) {
	item := x.(*huffmanNode)
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]

	return item
}

type prefixTable map[rune]string

func buildTree(freqMap map[rune]int) (*huffmanNode, prefixTable, error) {
	if len(freqMap) == 0 {
		return nil, nil, errors.New("frequency map is empty")
	}

	pq := make(priorityQueue, len(freqMap))
	i := 0
	for char, freq := range freqMap {
		if freq <= 0 {
			return nil, nil, errors.New("frequency must be greater than zero")
		}
		pq[i] = &huffmanNode{Char: char, Count: freq}
		i++
	}
	heap.Init(&pq)

	if len(pq) < 2 {
		return nil, nil, errors.New("priority queue must contain at least two nodes")
	}

	count := 1
	for len(pq) > 1 {
		left := heap.Pop(&pq).(*huffmanNode)
		right := heap.Pop(&pq).(*huffmanNode)

		parent := &huffmanNode{
			Count: left.Count + right.Count,
			Left:  left,
			Right: right,
		}

		heap.Push(&pq, parent)
		count += 2
	}

	head := heap.Pop(&pq).(*huffmanNode)
	if head == nil {
		return nil, nil, errors.New("failed to build tree")
	}

	prefix := make(prefixTable, len(freqMap))
	encodeTree(head, "", prefix)

	return head, prefix, nil
}

func encodeTree(node *huffmanNode, code string, preTab prefixTable) {
	if node == nil {
		return
	}

	if node.Left == nil && node.Right == nil {
		node.Code = code
		preTab[node.Char] = code
	}

	encodeTree(node.Left, code+"0", preTab)
	encodeTree(node.Right, code+"1", preTab)
}

func getTreeFrequency(root *huffmanNode) []*huffmanNode {
	var nodes []*huffmanNode

	traverseTree(root, &nodes)

	return nodes
}

func traverseTree(node *huffmanNode, nodes *[]*huffmanNode) {
	if node == nil {
		return
	}

	*nodes = append(*nodes, node)

	traverseTree(node.Left, nodes)

	traverseTree(node.Right, nodes)
}
