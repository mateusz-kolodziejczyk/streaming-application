package main

import "os"

// clearDirectory Deletes all files within a directory. It does this by removing the directory and recreating it afterwards.
func clearDirectory(directory string) {
	err := os.RemoveAll(directory)
	if err != nil {
		println(err)
	}
	err = os.Mkdir(directory, os.ModePerm)
	if err != nil {
		println(err)
	}
}
