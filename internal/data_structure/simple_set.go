package data_structure

type SimpleSet struct {
	key  string
	dict map[string]struct{}
}

type SetStore map[string]*SimpleSet

func CreateSimpleSet(key string) *SimpleSet {
	return &SimpleSet{
		key:  key,
		dict: make(map[string]struct{}),
	}
}

func CreateSimpleSetStore() SetStore {
	var setStore SetStore = make(SetStore)
	return setStore
}

func (s *SimpleSet) Add(members ...string) int {
	var addedMem int = 0
	for _, mem := range members {
		if _, exist := s.dict[mem]; !exist {
			s.dict[mem] = struct{}{}
			addedMem++
		}
	}
	return addedMem
}

func (s *SimpleSet) Rem(members ...string) int {
	var removedMem int = 0
	for _, mem := range members {
		if _, exist := s.dict[mem]; exist {
			delete(s.dict, mem)
			removedMem++
		}
	}
	return removedMem
}

func (s *SimpleSet) IsMember(mem string) int {
	_, exist := s.dict[mem]
	if exist {
		return 1
	}
	return 0
}

func (s *SimpleSet) Members() []string {
	var res []string = make([]string, 0)
	for k, _ := range s.dict {
		res = append(res, k)
	}
	return res
}
