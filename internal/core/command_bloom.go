package core

import (
	"fmt"
	"ngvanthanggit/internal/constant"
	"ngvanthanggit/internal/data_structure"
	"strconv"
)

// BF.RESERVE key error_rate capacity
func cmdBFRESERVE(args []string) []byte {
	if len(args) != 3 {
		return Encode(ErrNumberOfArgs, false)
	}

	key := args[0]
	errRate, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return Encode(fmt.Errorf("error rate must be a floating point number %s", args[1]), false)
	}
	capacity, err := strconv.ParseUint(args[2], 10, 64)
	if err != nil {
		return Encode(fmt.Errorf("capacity must be an integer number %s", args[2]), false)
	}

	_, exist := bloomStore[key]
	if exist {
		return Encode(fmt.Errorf("Bloom filter with key '%s' already exist", key), false)
	}
	bloom := data_structure.NewBloomFilter(capacity, errRate)
	bloomStore[key] = bloom

	return constant.RespOk
}

// BF.MADD key item [item...]
func cmdBFMADD(args []string) []byte {
	if len(args) < 2 {
		return Encode(ErrNumberOfArgs, false)
	}

	key := args[0]
	bloom, exist := bloomStore[key]
	if !exist {
		bloom = data_structure.NewBloomFilter(constant.BfDefaultInitCapacity, constant.BfDefaultErrRate)
		bloomStore[key] = bloom
	}

	var res []string = make([]string, 0)
	for i := 1; i < len(args); i++ {
		item := args[i]
		isNew := bloom.Add(item)
		if isNew {
			res = append(res, fmt.Sprintf("%d", 1))
		} else {
			res = append(res, fmt.Sprintf("%d", 0))
		}
	}

	return Encode(res, false)
}

// BF.EXISTS key item
func cmdBFEXISTS(args []string) []byte {
	if len(args) != 2 {
		return Encode(ErrNumberOfArgs, false)
	}

	key := args[0]
	bloom, exist := bloomStore[key]
	if !exist {
		return Encode(fmt.Errorf("Bloom filter with key %s does not exist", key), false)
	}

	if !bloom.Exist(args[1]) {
		return Encode(0, false)
	}
	return Encode(1, false)
}
