package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Entry struct {
	name     string
	path     string
	children int
	size     int64
	last     bool
	dir      bool
}

func str(str string) string {
	var r int
	r = strings.LastIndex(str, "\\")

	if r >= 0 {
		str = str[:r]
	} else
	{
		str = "str"  // заглушка для корневых файлов и папок
	}

	return str
}
func count(str string) int {
	var r int
	r = strings.Count(str, "\\")
	return r
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	var Entries []Entry
	var i = 0
	var children int

	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			children = count(path)
			en := []Entry{
				{
					name:     info.Name(),
					path:     path,
					children: children,
					size:     info.Size(),
					last:     true,
					dir:      info.IsDir(),
				},
			}
			Entries = append(Entries, en...)
			i++

			return nil
		})
	if err != nil {
		log.Println(err)
	}
	var (
		startSymbol, nextStartSymbol, nextPrefix, size, last string
	)
	pathP := make(map[int]string)

	// проверяем, является ли элемент последним

	for j := 1; j < i; j++ {
		for k := j + 1; k < i; k++ {
			if printFiles == true {          // проверяем, является ли элемент последним (папки и файлы)
				if str(Entries[j].path) == str(Entries[k].path) {
					Entries[j].last = false
					Entries[k].last = true
				}
			} else {                           // проверяем, является ли элемент последним (папки)
				if (Entries[j].dir == true) && (Entries[k].dir == true) && (str(Entries[j].path) == str(Entries[k].path)) {
					Entries[j].last = false
					Entries[k].last = true
				}
			}
		}

		last = nextPrefix

		if Entries[j].last == true {
			startSymbol = `└───`
			nextStartSymbol = ``

		} else {
			startSymbol = `├───`
			nextStartSymbol = `│`

		}

		if Entries[j].children <= Entries[j-1].children {
			last = pathP[Entries[j].children]
		}

		if Entries[j].dir == true {
			fmt.Fprintf(out, "%s%s%s\n", last, startSymbol, Entries[j].name)
			nextPrefix = fmt.Sprintf("%s%s\t", last, nextStartSymbol)

		} else if printFiles == true {
			if Entries[j].size > 0 {
				size = fmt.Sprintf("%db", Entries[j].size)
			} else {
				size = "empty"
			}
			fmt.Fprintf(out, "%s%s%s (%s)\n", last, startSymbol, Entries[j].name, size)
		}
		pathP[Entries[j].children] = last
	}

	return nil

}
func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
