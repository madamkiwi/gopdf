package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const bufsize int = 512
const dbsize int = 256
const config string = "config"

type Opts struct {
	startPage  int
	endPage    int
	inputFile  string
	outputFile string
}

func main() {
	fmt.Println("Go PDF")
	if len(os.Args) < 2 {
		usage()
		return
	}
	op := string(os.Args[1])
	switch {
	case op == "s":
		if len(os.Args) < 5 {
			usage()
			return
		}
	case op == "m":
		if len(os.Args) < 4 {
			usage()
			return
		}
		merge(string(os.Args[2]), string(os.Args[3]))
	}

}

func merge(inputDir string, outputFile string) {
	_, err := os.Stat(inputDir)
	if err != nil || inputDir == outputFile {
		return
	}
	dir, err := os.Open(inputDir)
	if err != nil {
		return
	}
	files, err := dir.Readdir(-1)
	mergelist := []string{}
	for _, file := range files {
		if file.Name()[0] == '.' {
			continue
		}
		mergelist = append(mergelist, fmt.Sprintf("%s/%s", inputDir, file.Name()))
	}
	if len(mergelist) == 0 {
		return
	}
	args := strings.Split(
		fmt.Sprintf("-dBATCH -dNOPAUSE -q -sDEVICE=pdfwrite -dPDFSETTINGS=/prepress -sOutputFile=%s.pdf %s",
			outputFile, strings.Join(mergelist, " ")), " ")
	_, err = exec.Command("gs", args...).Output()

	if err != nil {
		fmt.Println("merge-err %s", err)
		return
	}
}

func usage() {
	fmt.Printf("\n\tprocess\t<input file>")
	fmt.Printf("\n\tsplit\t<input file>\t<start page #>\t<end page #>\t<output file>")
	fmt.Printf("\n\tmerge\t<input dir>\t<output file>")
	fmt.Printf("\n\n")
}
