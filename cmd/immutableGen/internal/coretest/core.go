package coretest

//go:generate immutableGen -licenseFile license.txt -G "echo \"hello world\""

// a comment about MyMap
type _Imm_MyMap map[string]int

// a comment about Slice
type _Imm_MySlice []string

// a comment about myStruct
type _Imm_MyStruct struct {

	// my field comment
	//somethingspecial
	/*

		Heelo

	*/
	Name, surname string `tag:"value"`
	age           int    `tag:"age"`

	*string

	fieldWithoutTag bool
}

func main() {
}
