package data_structure

type SortedSet struct {
	Tree         *BPlusTree
	MemberScores map[string]float64
}

type SortedSetStore map[string]*SortedSet

func NewSortedSet(degree int) *SortedSet {
	return &SortedSet{
		Tree:         NewBPlusTree(degree),
		MemberScores: make(map[string]float64),
	}
}

func CreateSortedSetStore() SortedSetStore {
	var sortedSetStore SortedSetStore = make(SortedSetStore)
	return sortedSetStore
}

func (s *SortedSet) Add(score float64, member string) int {
	return s.Tree.Add(score, member)
}

func (s *SortedSet) GetScore(member string) (float64, bool) {
	return s.Tree.Score(member)
}

func (s *SortedSet) GetRank(member string) int {
	return s.Tree.GetRank(member)
}
