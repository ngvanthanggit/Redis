package core

import "ngvanthanggit/internal/data_structure"

// SADD key mem1 mem2
func cmdSADD(args []string) []byte {
	if len(args) < 2 {
		return Encode(ErrNumberOfArgs, false)
	}

	key := args[0]
	set, exist := simpleSetStore[key]
	if !exist {
		set = data_structure.CreateSimpleSet(key)
		simpleSetStore[key] = set
	}
	count := set.Add(args[1:]...)
	return Encode(count, false)
}

// SREM key mem1 mem2
func cmdSREM(args []string) []byte {
	if len(args) < 2 {
		return Encode(ErrNumberOfArgs, false)
	}

	key := args[0]
	set, exist := simpleSetStore[key]
	if !exist {
		set = data_structure.CreateSimpleSet(key)
		simpleSetStore[key] = set
	}
	count := set.Rem(args[1:]...)
	return Encode(count, false)
}

// SISMEMBER key mem
func cmdSISMEMBER(args []string) []byte {
	if len(args) < 2 {
		return Encode(ErrNumberOfArgs, false)
	}

	key := args[0]
	set, exist := simpleSetStore[key]
	if !exist {
		return Encode(0, false)
	}
	return Encode(set.IsMember(args[1]), false)
}

// SMEMBERS key
func cmdSMEMBERS(args []string) []byte {
	if len(args) < 1 {
		return Encode(ErrNumberOfArgs, false)
	}

	key := args[0]
	set, exist := simpleSetStore[key]
	if !exist {
		return Encode(make([]string, 0), false)
	}
	return Encode(set.Members(), false)
}
