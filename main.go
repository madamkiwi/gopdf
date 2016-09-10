package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
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
	case op == "s" || op == "split":
		if len(os.Args) < 6 {
			usage()
			return
		}
		inputFile := fmt.Sprintf("%s.pdf", string(os.Args[2]))
		outputFile := fmt.Sprintf("%s.pdf", string(os.Args[5]))
		startPage, _ := strconv.Atoi(string(os.Args[3]))
		endPage, _ := strconv.Atoi(string(os.Args[4]))
		split(inputFile, startPage, endPage, outputFile)
	case op == "m" || op == "merge":
		if len(os.Args) < 4 {
			usage()
			return
		}
		merge(string(os.Args[2]), string(os.Args[3]))
	case op == "p" || op == "process":
	}

}

func split(inputFile string, startPage int, endPage int, outputFile string) {
	_, err := os.Stat(inputFile)
	if err != nil || inputFile == outputFile {
		fmt.Println("input-and-output-are-the-same %s", err)
		return
	}
	args := strings.Split(
		fmt.Sprintf("-dNOPAUSE -dBATCH -sOutputFile=%s -dFirstPage=%d -dLastPage=%d -sDEVICE=pdfwrite %s",
			outputFile, startPage, endPage, inputFile), " ")
	_, err = exec.Command("gs", args...).Output()
	if err != nil {
		fmt.Println("split-err %s", err)
		return
	}

}

func merge(inputDir string, outputFile string) {
	_, err := os.Stat(inputDir)
	if err != nil || inputDir == outputFile {
		fmt.Println("input-and-output-are-the-same %s", err)
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
	fmt.Printf("\n\tmerge\t<input dir>\t<output file>")
	fmt.Printf("\n\tsplit\t<input file>\t<start page #>\t<end page #>\t<output file>")
	fmt.Printf("\n\tparse\t<input file>\t<start page #>\t<file>\t<C|I>\t<start page #>\t<file2>\t<C|I>...")
	fmt.Printf("\n\n")
}
