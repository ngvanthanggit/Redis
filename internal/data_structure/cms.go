package data_structure

import (
	"math"

	"github.com/spaolacci/murmur3"
)

// pre-calculated value for log10(0.5)
const Log10PointFive float64 = -0.30102999566

type CMS struct {
	width   uint32
	depth   uint32
	counter [][]uint32 // count the frequencies of the items
}

type CmsStore map[string]*CMS

func CreateCmsStore() CmsStore {
	cmsStore := make(map[string]*CMS)
	return cmsStore
}

func NewCMS(w uint32, d uint32) *CMS {
	cms := &CMS{
		width: w,
		depth: d,
	}

	cms.counter = make([][]uint32, d)
	for i := range d {
		cms.counter[i] = make([]uint32, w)
	}
	return cms
}

// this function is to calculate the dimensions (w, d) of the CMS counter table
// based on the desired error rate and error probability
func CalcCMSDim(errRate float64, errProb float64) (uint32, uint32) {
	w := uint32(math.Ceil(2.0 / errRate))
	d := uint32(math.Ceil(math.Log10(errProb) / Log10PointFive))

	return w, d
}

// calculate a 32-bit hash for the given item and seed using MurmurHash
func (c *CMS) calcHash(item string, seed uint32) uint32 {
	hasher := murmur3.New32WithSeed(seed)
	hasher.Write([]byte(item))
	return hasher.Sum32()
}

func (c *CMS) IncrBy(item string, value uint32) uint32 {
	var minCount uint32 = math.MaxUint32

	for i := range c.depth {
		// calculate hash value for the item with the seed = row index, i
		hash := c.calcHash(item, i)
		// use hash to calculate the column index
		j := hash % c.width

		// safely add the value to the counter -> prevent overflow
		if math.MaxUint32-c.counter[i][j] < value {
			c.counter[i][j] = math.MaxUint32
		} else {
			c.counter[i][j] += value
		}

		// keep track of the minimum count across all rows
		if c.counter[i][j] < minCount {
			minCount = c.counter[i][j]
		}
	}
	return minCount
}

func (c *CMS) Count(item string) uint32 {
	var minCount uint32 = math.MaxUint32

	for i := range c.depth {
		hash := c.calcHash(item, i)
		j := hash % c.width

		// get the minimum count of hash(item) across all the rows
		minCount = min(minCount, c.counter[i][j])
	}
	return minCount
}
