package core

import "ngvanthanggit/internal/data_structure"

var dictStore *data_structure.Dict

// it's so cool to know that this function will be automatically run
func init() {
	dictStore = data_structure.CreateDict()
}
