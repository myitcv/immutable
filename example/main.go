// Copyright (c) 2016 Paul Jolly <paul@myitcv.org.uk>, all rights reserved.
// Use of this document is governed by a license found in the LICENSE document.

//go:generate immutableGen -licenseFile license_header.txt -G "echo \"hello world\""

package main

func main() {
}

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
}
