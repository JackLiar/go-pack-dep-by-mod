package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

// ErrNoDepencencyFound No go depencency mod was found error
var ErrNoDepencencyFound error = errors.New("No go depencency module was found")

// ErrInvalidGoModFile Invalid go.mod ifle error
var ErrInvalidGoModFile error = errors.New("Invalid go.mod file")

var targetDir string = "./"
var modFile string = ""

func readGoMod(fpath string) error {
	content, err := ioutil.ReadFile(fpath)
	if err != nil {
		os.Exit(1)
	}

	contentStr := string(content)
	lBracketIndex := strings.Index(contentStr, "(")
	if lBracketIndex < 0 {
		return ErrNoDepencencyFound
	}

	rBracketIndex := strings.Index(contentStr, ")")
	if rBracketIndex < 0 {
		return ErrInvalidGoModFile
	}

	mods := strings.Split(contentStr[lBracketIndex+2:rBracketIndex-1], "\n")
	actualMods := make([]string, 0)
	goPath := os.Getenv("GOPATH")
	wd, err := os.Getwd()

	for _, modStr := range mods {
		modStr = strings.Trim(modStr, "\t")

		commentIndex := strings.Index(modStr, " // ")
		if commentIndex >= 0 {
			modStr = modStr[:commentIndex]
		}

		modStr = strings.Replace(modStr, " ", "@", 1)
		absModStr := path.Join(goPath, "pkg", "mod", strings.Replace(modStr, " ", "@", 1))
		if _, err := os.Stat(absModStr); os.IsNotExist(err) {
			log.Fatalf("%s does not exists", absModStr)
		}
		actualMods = append(actualMods, modStr)
	}

	args := make([]string, 0)
	args = append(args, "-czf", path.Join(wd, "mod.tar.gz"))
	for _, mod := range actualMods {
		args = append(args, mod)
	}

	cmd := exec.Command("tar", args...)
	cmd.Dir = path.Join(goPath, "pkg", "mod")
	_, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	return nil
}

func parseCli() {
	flag.StringVar(&targetDir, "i", "", "target go project dir")
	flag.Parse()

	if targetDir == "" {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatalln(err.Error())
		}
		fmt.Printf("No target dir is provided, using current dir:%s\n", wd)
	}

	modFile = path.Join(targetDir, "go.mod")
}

func main() {
	parseCli()
	if err := readGoMod(modFile); err != nil {
		fmt.Printf("%s, exit.\n", err.Error())
	}
}
