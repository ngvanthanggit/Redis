package core

import "ngvanthanggit/internal/data_structure"

var dictStore *data_structure.Dict
var simpleSetStore data_structure.SetStore
var sortedSetStore data_structure.SortedSetStore
var cmsStore data_structure.CmsStore

// it's so cool to know that this function will be automatically run
func init() {
	dictStore = data_structure.CreateDict()
	simpleSetStore = data_structure.CreateSimpleSetStore()
	sortedSetStore = data_structure.CreateSortedSetStore()
	cmsStore = data_structure.CreateCmsStore()
}
