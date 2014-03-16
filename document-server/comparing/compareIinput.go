package main

import (
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	"io/ioutil"
)

func readfile(path string) string {
	inbyte, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	text := string(inbyte[:])

	return text
}

func differances(text1, text2 string) {
	// Check diffs
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(text1, text2, false)
	// Make patch
	patch := dmp.PatchMake(text1, diffs)

	// Print patch to see the results before they're applied
	fmt.Println(dmp.PatchToText(patch))

	result := dmp.PatchToText(patch)

	// detta sker efter att vi har mottagit från klienten!
	// Apply patch
	//text3, _ := dmp.PatchApply(patch, text1)

	// Write the patched text to a new file
	err := ioutil.WriteFile("resultClient.txt", []byte(result), 0644)
	if err != nil {
		panic(err)
	}

}

func main() {
	slice1 := readfile("serverShadow.txt")
	slice2 := readfile("serverText.txt")

	differances(slice1, slice2)

}
