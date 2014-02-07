package main

import (
	"fmt"
	tp "github.com/vizidrix/goio/tarpack"
	"io/ioutil"
	"os"
)

func main() {
	var err error
	argIndex := 1
	targetDir := "."
	outFile := "tarpack.tar"
	filters := make([]tp.FilePredicate, 1, 10)
	// Filter out git files by default
	filters[0] = tp.OnDirs(tp.Not(tp.NameContainsAny(".git")))
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "toml":
			filters = append(filters, tp.OnFiles(tp.PathMatchAny("*.toml")))
			outFile = "toml.tar"
			argIndex++
			break
		case "web":
			filters = append(filters, tp.OnFiles(tp.PathMatchAny("*.js", "*.html", "*.png", "*.css", "*.ttf")))
			outFile = "web.tar"
			argIndex++
			break
		default:
			break
		}
		if len(os.Args) > 2 {
			targetDir = os.Args[argIndex]
		}
		if len(os.Args) > 3 {
			outFile = os.Args[argIndex+1]
		}
	}
	var buffer []byte
	if buffer, err = tp.TarDir(targetDir, "", filters...); err != nil {
		fmt.Printf("Error zipping TarPack files:\n\t%s\n\n", err)
	}
	ioutil.WriteFile(outFile, buffer, 0666)
}
