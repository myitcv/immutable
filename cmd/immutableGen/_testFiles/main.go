package main

//go:generate immutableGen -licenseFile license.txt -G true

// a comment about myMap
type Imm_myMap map[string]int

// a comment about Slice
type Imm_Slice []*string

// a comment about myStruct
type Imm_myStruct struct {

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
