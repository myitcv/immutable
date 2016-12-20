package main

//go:generate immutableGen -licenseFile license.txt -G "echo \"hello world\""

// a comment about myMap
type _Imm_myMap map[string]int

// a comment about Slice
type _Imm_Slice []*string

// a comment about myStruct
type _Imm_myStruct struct {

	// my field comment
	//somethingspecial
	/*

		Heelo

	*/
	Name, surname string `tag:"value"`
	age           int    `tag:"age"`

	fieldWithoutTag bool
}

func main() {
}
