package main

import (
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	"io/ioutil"
)

func main() {

	// Read whole file as string
	inbyte, err := ioutil.ReadFile("text1.txt")
	if err != nil {
		panic(err)
	}
	text1 := string(inbyte[:])

	inbyte, err = ioutil.ReadFile("text2.txt")
	if err != nil {
		panic(err)
	}
	text2 := string(inbyte[:])

	// Check diffs
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(text1, text2, false)

	// Make patch
	patch := dmp.PatchMake(text1, diffs)

	// Print patch to see the results before they're applied
	fmt.Println(dmp.PatchToText(patch))

	// Apply patch
	text3, _ := dmp.PatchApply(patch, text1)

	// Write the patched text to a new file
	err = ioutil.WriteFile("text3.txt", []byte(text3), 0644)
	if err != nil {
		panic(err)
	}

}
