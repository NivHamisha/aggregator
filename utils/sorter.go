package utils

import "sort"

type LessFunc func(p1, p2 *Expense) bool

type CompereFunctions map[string]LessFunc

type MultiSorter struct {
	expenses []Expense
	less     []LessFunc
}

func (ms *MultiSorter) Sort(expenses []Expense) {
	ms.expenses = expenses
	sort.Sort(ms)
}

func (ms *MultiSorter) Len() int {
	return len(ms.expenses)
}

func (ms *MultiSorter) Swap(i, j int) {
	ms.expenses[i], ms.expenses[j] = ms.expenses[j], ms.expenses[i]
}

func (ms *MultiSorter) Less(i, j int) bool {
	p, q := &ms.expenses[i], &ms.expenses[j]
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]
		switch {
		case less(p, q):
			return true
		case less(q, p):
			return false
		}
	}
	return ms.less[k](p, q)
}

func OrderedBy(less []LessFunc) *MultiSorter {
	return &MultiSorter{
		less: less,
	}
}