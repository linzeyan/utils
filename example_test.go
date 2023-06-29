package utils_test

import (
	"fmt"

	"github.com/linzeyan/utils"
)

func ExampleHasNullByte() {
	type exampleType struct {
		data []byte
	}

	var x = []exampleType{
		{[]byte{0, 1, 2, 3}},
		{[]byte{44, 55, 00, 77, 88}},
		{[]byte{111, 222}},
	}
	for i := range x {
		fmt.Println(utils.HasNullByte(x[i].data))
	}
	// output:
	// true
	// true
	// false
}
