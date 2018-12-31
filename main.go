package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
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

type flags struct {
	strRex        string
	recursive     string
	fileFilter    string
	verbose       bool
	inverseSearch bool
}

var (
	fl            flags
	filenames     []string
	regRex        *regexp.Regexp
	rexfileFilter *regexp.Regexp
)

func init() {
	dfl := flags{
		strRex:        "",
		verbose:       false,
		inverseSearch: false,
		recursive:     "",
		fileFilter:    "",
	}

	var rexError error

	flag.StringVar(&fl.strRex, "r", dfl.strRex,
		"Regular expresion to match against the input files")

	flag.StringVar(&fl.recursive, "R", dfl.recursive,
		"Recursively find all files starting from the current folder and "+
			"apply the given search to them")

	flag.StringVar(&fl.fileFilter, "FF", dfl.fileFilter,
		"Filter to be applied to the filenames when used recursevily")

	flag.BoolVar(&fl.verbose, "v", dfl.verbose,
		"It sets verbose output (Basically showing filename and line number "+
			"for each match)")

	flag.BoolVar(&fl.inverseSearch, "i", dfl.inverseSearch,
		"It does what you might expect.. reverse the search")

	flag.Parse()

	if fl.strRex == "" {
		log.Fatalln("The regular expression flag '-r' is mandatory")
	}

	regRex, rexError = regexp.Compile(fl.strRex)

	if rexError != nil {
		log.Fatalf("Your regex '%s' cant compile. Error : %s\n", fl.strRex,
			rexError)
	}

	rexfileFilter, rexError = regexp.Compile(fl.fileFilter)

	if rexError != nil {
		log.Fatalf("Your regex '%s' cant compile. Error : %s", rexfileFilter,
			rexError)
	}

	if fl.recursive != "" {
		var err error

		if filenames, err = walkDir(fl.recursive); err != nil {
			log.Println(err)
		}
	} else {
		filenames = flag.Args()
	}
}

func main() {
	stat, err := os.Stdin.Stat()

	if err != nil {
		log.Fatalf("There is an error reading from stdin: %s", err)
	}

	var wait sync.WaitGroup

	if (stat.Mode() & os.ModeNamedPipe) != 0 {
		grepStdin(os.Stdin, regRex)
	} else {
		chResults := make(chan *result, 4)

		wait.Add(len(filenames))

		for _, fn := range filenames {
			go grep(fn, regRex, &wait, chResults)
		}

		go func(wait *sync.WaitGroup, ch chan<- *result) {
			wait.Wait()

			close(ch)
		}(&wait, chResults)

		for res := range chResults {
			if fl.verbose {
				formatRes(res, 1)
			} else {
				formatRes(res, 2)
			}
		}
	}
}

func match(reg *regexp.Regexp, line string) bool {
	return !fl.inverseSearch && reg.Match([]byte(line)) || (fl.inverseSearch && !reg.Match([]byte(line)))
}

func grepStdin(ptr io.Reader, reg *regexp.Regexp) {
	bf := bufio.NewScanner(ptr)

	for l := 1; bf.Scan(); l++ {
		if line := bf.Text(); match(reg, line) {
			formatRes(&result{
				Line:       line,
				LineNumber: l,
				Error:      nil,
			}, 3)
		}
	}
}

func grep(file string, reg *regexp.Regexp, wait *sync.WaitGroup,
	ch chan<- *result) {

	fd, err := os.Open(file)

	if err != nil {
		ch <- &result{
			Filename: file,
			Error:    err,
		}
	}

	bf := bufio.NewScanner(fd)

	for l := 1; bf.Scan(); l++ {
		if line := bf.Text(); match(reg, line) {
			ch <- &result{
				Filename:   file,
				Line:       line,
				LineNumber: l,
				Error:      nil,
			}
		}
	}

	wait.Done()
}

func formatRes(r *result, format int) {
	switch format {
	case 1:
		if r.Error == nil {
			fmt.Printf("%d - %s - %s\n", r.LineNumber, r.Filename, r.Line)
		} else {
			log.Printf("%s - %s\n", r.Filename, r.Error)
		}
		break
	case 2:
		if r.Error == nil {
			fmt.Println(r.Line)
		} else {
			log.Printf("%s - %s\n", r.Filename, r.Error)
		}
		break
	case 3:
		if r.Error == nil {
			fmt.Println(r.Line)
		} else {
			log.Println(r.Error)
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

			if fl.fileFilter != "" {
				if rexfileFilter.Match([]byte(filepath.Base(path))) {
					list = append(list, path)
				}
			} else {
				list = append(list, path)
			}

			return nil
		})

	if err != nil {
		return nil, err
	}

	return list, nil
}
