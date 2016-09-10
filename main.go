package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

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
		inputFile := fmt.Sprintf("%s.pdf", os.Args[2])
		outputFile := fmt.Sprintf("%s.pdf", os.Args[5])
		startPage, _ := strconv.Atoi(os.Args[3])
		endPage, _ := strconv.Atoi(os.Args[4])
		split(inputFile, startPage, endPage, outputFile)
	case op == "m" || op == "merge":
		if len(os.Args) < 4 {
			usage()
			return
		}
		merge(os.Args[2], os.Args[3])
	case op == "p" || op == "process":
		if len(os.Args) < 5 {
			usage()
			return
		}
		os.Mkdir("./Complete", 0777)
		os.Mkdir("./Incomplete", 0777)
		process(os.Args...)
	}

}

func process(args ...string) {
	var (
		outputFile string
		startPage, endPage  int
		dirpath    string
	)

	inputFile := fmt.Sprintf("%s.pdf", args[2])

	fmt.Println("processing input file %s", inputFile)
	idx := 0
	for _, _ = range args {
		if idx >= len(args) {
			return
		}
		if idx < 3 {
			idx++
			continue
		}
		startPage, _ = strconv.Atoi(args[idx])
		idx++
		outputFile = fmt.Sprintf("%s.pdf", args[idx])
		idx++
		switch {
		case args[idx] == "c" || args[idx] == "C":
			dirpath = "./Complete"
		case args[idx] == "i" || args[idx] == "I":
			dirpath = "./Incomplete"
		default:
			fmt.Println("expecting-c-or-i-but-got %s", args[idx])
			return
		}
		idx++
		if idx >= len(args) {
			endPage = totalPages(inputFile)
		} else {
			endPage, _ = strconv.Atoi(args[idx])
			endPage = endPage - 1
		}
		fmt.Printf("%s (p.%d-p.%d) -> %s\n", inputFile, startPage, endPage, fmt.Sprintf("%s/%s", dirpath, outputFile))
		err := split(inputFile, startPage, endPage, fmt.Sprintf("%s/%s", dirpath, outputFile))
		if err != nil {
			return
		}
	}

}

func totalPages(file string) int{
	output, _ := exec.Command("gs", "-q", "-dNODISPLAY", "-c", fmt.Sprintf("(%s) (r) file runpdfbegin pdfpagecount = quit", file)).Output()
	page, _:= strconv.Atoi(strings.Replace(string(output), "\n", "", -1))

	return page
}

func split(inputFile string, startPage int, endPage int, outputFile string) error {
	_, err := os.Stat(inputFile)
	if err != nil || inputFile == outputFile {
		fmt.Println("input-and-output-are-the-same %s", err)
		return err
	}
	args := strings.Split(
		fmt.Sprintf("-dNOPAUSE -dBATCH -sOutputFile=%s -dFirstPage=%d -dLastPage=%d -sDEVICE=pdfwrite %s",
			outputFile, startPage, endPage, inputFile), " ")
	_, err = exec.Command("gs", args...).Output()
	if err != nil {
		fmt.Println("split-err %s", err)
		return err
	}
	return nil
}

func merge(inputDir string, outputFile string) {
	_, err := os.Stat(inputDir)
	if err != nil || inputDir == fmt.Sprintf("%s.pdf", outputFile) {
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
