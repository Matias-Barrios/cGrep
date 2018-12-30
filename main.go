package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

type result struct {
	Filename   string
	Line       string
	LineNumber int
	Error      error
}

var strRex string
var filenames []string
var regRex *regexp.Regexp
var wg sync.WaitGroup
var allResults []result
var verbose = false
var recursive string
var recursiveFileList []string
var fileFilter string
var rexfileFilter *regexp.Regexp
var inverseSearch bool

func init() {
	var rexError error
	flag.StringVar(&strRex, "r", "", "Regular expresion to match against the input files")
	flag.BoolVar(&verbose, "v", false, "It sets verbose output (Basically showing filename and line number for each match)")
	flag.BoolVar(&inverseSearch, "i", false, "It does what you might expect.. reverse the search")
	flag.StringVar(&recursive, "R", "", "Recursively find all files starting from the current folder and apply the given search to them")
	flag.StringVar(&fileFilter, "FF", "", "Filter to be applied to the filenames when used recursevily")
	flag.Parse()
	if strRex == "" {
		fmt.Fprintf(os.Stderr, "The '-r' (regular expression flag is mandatory)\n")
		os.Exit(1)
	}
	regRex, rexError = regexp.Compile(strRex)
	if rexError != nil {
		fmt.Fprintf(os.Stderr, "Your regex '%s' cant compile. Error : %s", strRex, rexError.Error())
		os.Exit(2)
	}
	rexfileFilter, rexError = regexp.Compile(fileFilter)
	if rexError != nil {
		fmt.Fprintf(os.Stderr, "Your regex '%s' cant compile. Error : %s", rexfileFilter, rexError.Error())
		os.Exit(3)
	}
	if recursive != "" {
		var err error
		filenames, err = walkDir(recursive)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err.Error())
		}

	} else {

		filenames = flag.Args()
	}

}

func main() {

	stat, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "There is an error reading from stdin : %s", err)
		os.Exit(3)
	}
	if (stat.Mode() & os.ModeNamedPipe) != 0 {
		grepStdin(os.Stdin, regRex)
	} else {
		chResults := make(chan *result, 4)
		wg.Add(len(filenames))
		for _, fn := range filenames {
			go grep(fn, regRex, &wg, chResults)
		}
		go func(wait *sync.WaitGroup, ch chan<- *result) {
			wg.Wait()
			close(ch)
		}(&wg, chResults)

		for res := range chResults {
			if verbose {
				formatRes(res, 1)
			} else {
				formatRes(res, 2)
			}

		}
	}
}

func grepStdin(ptr io.Reader, reg *regexp.Regexp) {
	bf := bufio.NewScanner(ptr)
	var lineno = 1

	for bf.Scan() {
		// There is no XOR in Golang, so you ahve to do this :
		if line := bf.Text(); (reg.Match([]byte(line)) && !inverseSearch) || (!reg.Match([]byte(line)) && inverseSearch) {

			formatRes(&result{
				Line:       line,
				LineNumber: lineno,
				Error:      nil,
			}, 3)
		}
		lineno++
	}
}

func grep(file string, reg *regexp.Regexp, wait *sync.WaitGroup, ch chan<- *result) {
	fd, err := os.Open(file)
	if err != nil {
		ch <- &result{
			Filename: file,
			Error:    err,
		}
	}
	bf := bufio.NewScanner(fd)
	var lineno = 1

	for bf.Scan() {
		// There is no XOR in Golang, so you ahve to do this :
		if line := bf.Text(); (reg.Match([]byte(line)) && !inverseSearch) || (!reg.Match([]byte(line)) && inverseSearch) {

			ch <- &result{
				Filename:   file,
				Line:       line,
				LineNumber: lineno,
				Error:      nil,
			}
		}
		lineno++
	}
	wg.Done()
}

func formatRes(r *result, format int) {
	if format == 1 {
		if r.Error == nil {
			fmt.Printf("%d - %s - %s\n", r.LineNumber, r.Filename, r.Line)
		} else {
			fmt.Fprintf(os.Stderr, "%s - %s \n", r.Filename, r.Error)
		}
	}
	if format == 2 {
		if r.Error == nil {
			fmt.Printf("%s\n", r.Line)
		} else {
			fmt.Fprintf(os.Stderr, "%s - %s \n", r.Filename, r.Error)
		}
	}
	if format == 3 {
		if r.Error == nil {
			fmt.Printf("%s\n", r.Line)
		} else {
			fmt.Fprintf(os.Stderr, "%s\n", r.Error)
		}
	}
}

func walkDir(path string) ([]string, error) {
	list := make([]string, 0, 50)
	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if fileFilter != "" {
				if rexfileFilter.Match([]byte(filepath.Base(path))) {
					list = append(list, path)
				}
			} else {
				list = append(list, path)
			}
			return nil // Unreachable code
		})
	if err != nil {
		return nil, err
	}
	return list, nil
}
