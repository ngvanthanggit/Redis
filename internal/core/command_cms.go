package core

import (
	"errors"
	"fmt"
	"math"
	"ngvanthanggit/internal/constant"
	"ngvanthanggit/internal/data_structure"
	"strconv"
)

// CMS.INITBYDIM key width depth
func cmdCMSINITBYDIM(args []string) []byte {
	if len(args) != 3 {
		return Encode(ErrNumberOfArgs, false)
	}

	key := args[0]
	width, err := strconv.ParseUint(args[1], 10, 32)
	if err != nil {
		return Encode(fmt.Errorf("width must be an integer number %s", args[1]), false)
	}
	depth, err := strconv.ParseUint(args[2], 10, 32)
	if err != nil {
		return Encode(fmt.Errorf("depth must be an integer number %s", args[1]), false)
	}

	_, exist := cmsStore[key]
	if exist {
		return Encode(errors.New("CMS: key already exists"), false)
	}
	cms := data_structure.NewCMS(uint32(width), uint32(depth))
	cmsStore[key] = cms

	return constant.RespOk
}

// CMS.INITBYPROB key error probability
func cmdCMSINITBYPROB(args []string) []byte {
	if len(args) != 3 {
		return Encode(ErrNumberOfArgs, false)
	}

	key := args[0]
	errRate, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return Encode(fmt.Errorf("error must be a float number %s", args[1]), false)
	}
	errProb, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return Encode(fmt.Errorf("probability must be a float number %s", args[2]), false)
	}

	_, exist := cmsStore[key]
	if exist {
		return Encode(errors.New("CMS: key already exists"), false)
	}
	width, depth := data_structure.CalcCMSDim(errRate, errProb)
	cms := data_structure.NewCMS(width, depth)
	cmsStore[key] = cms

	return constant.RespOk
}

// CMS.INCRBY key item increment [item increment ...]
func cmdCMSINCRBY(args []string) []byte {
	if len(args) < 3 || len(args)%2 == 0 {
		return Encode(ErrNumberOfArgs, false)
	}

	key := args[0]

	cms, exist := cmsStore[key]
	if !exist {
		return Encode(errors.New("CMS: key does not exist"), false)
	}

	// store the frequencies of the items after the command
	// this array should store values in string format - for cases where the overflow happens
	var frequencies []string = make([]string, 0)
	for i := 1; i < len(args); i += 2 {
		// item = args[i] & increment = args[i+1]
		item := args[i]
		increment, err := strconv.ParseUint(args[i+1], 10, 32)
		if err != nil {
			return Encode(fmt.Errorf("increment must be a non negative integer number %s", args[i+1]), false)
		}

		counter := cms.IncrBy(item, uint32(increment))
		if counter == math.MaxUint32 {
			frequencies = append(frequencies, "CMS: INCRBY overflow")
		}
		frequencies = append(frequencies, fmt.Sprintf("%d", counter))
	}
	return Encode(frequencies, false)
}

// CMS.QUERY key item [item ...]
func cmdCMSQUERY(args []string) []byte {
	if len(args) < 2 {
		return Encode(ErrNumberOfArgs, false)
	}

	key := args[0]
	cms, exist := cmsStore[key]
	if !exist {
		return Encode(errors.New("CMS: key does not exist"), false)
	}

	var frequencies []string = make([]string, 0)
	for i := 1; i < len(args); i++ {
		item := args[i]
		counter := cms.Count(item)

		frequencies = append(frequencies, fmt.Sprintf("%d", counter))
	}
	return Encode(frequencies, false)
}
