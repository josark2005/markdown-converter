package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// Version
const VERSION = "0.0.1"

// PANDOC Release URL
const PANDOC_RELEASE_URL = ""

// PANDOC Accelerated Release URL
const PANDOC_RELEASE_ACC_URL = ""

var pandoc_bin string

// Initializaion
func init() {
	// check pandoc binary
	userhomedir, _ := os.UserHomeDir()
	mdcdir := userhomedir + "/.mdc"
	pandoc_bin, _ = filepath.Abs(mdcdir + "/pandoc.exe")
	_, err := os.Stat(pandoc_bin)
	if err != nil || os.IsNotExist(err) {
		println("!! Pandoc binary does NOT exist or is NOT permitted to access")
		os.Exit(2)
	}
}

func md2html(md string) string {
	println("Running: ", pandoc_bin, " -f markdown -t html ", md)
	output, err := exec.Command(pandoc_bin, "-f", "markdown", "-t", "html", md).Output()
	if err != nil {
		fmt.Println("!! Convert failed: ", err)
	}
	re := regexp.MustCompile(`(<img src=".+" alt=".*" />)[\s]{1,2}<figcaption aria-hidden="true">.+</figcaption>`)
	repl := `$1`
	repl_output := string(re.ReplaceAll(output, []byte(repl)))
	return repl_output
}

func md2html_w(mdfilepath string, filename string, perm fs.FileMode) error {
	// convert markdown to html
	res := md2html(mdfilepath)

	// write result to file
	err := ioutil.WriteFile(filename, []byte(res), perm)
	if err != nil {
		println("!! Error: ", err.Error())
		return err
	}
	return nil
}

func html2docx_w(html string, filename string) error {
	println("Running: ", pandoc_bin, " -f html -t docx ", html, "-o ", filename)
	_, err := exec.Command(pandoc_bin, "-f", "markdown", "-t", "docx", html, "-o", filename).Output()
	if err != nil {
		fmt.Println("!! Convert failed: ", err)
	}
	return err
}

func main() {
	args := os.Args

	// at least one param should be given
	if len(args) < 2 || args == nil {
		help()
		os.Exit(0)
	}

	// command or target
	command := args[1]

	// if command is help, ignore other params
	if command == "help" {
		help()
		os.Exit(0)
	}

	// generate filepath
	mdfilepath, _ := filepath.Abs(args[2])
	mdfilename := filepath.Base(mdfilepath)

	// remove the extension name
	_filename := strings.Split(mdfilename, ".")
	_filename = _filename[:len(_filename)-1]
	mdfilename_noext := strings.Join(_filename, ".")

	var err error

	// show pandoc version
	output, err := exec.Command(pandoc_bin, "--version").Output()
	if err != nil {
		println("!! Error occured when fetching the pandoc version")
		os.Exit(3)
	} else {
		println()
		println("============================ Pandoc  Info ============================")
		pandoc_version := strings.Split(string(output), "\n")
		println(pandoc_version[0])
		println(pandoc_version[1])
		println(pandoc_version[2])
		println(pandoc_version[3])
		println()
	}

	// if markdown file is exist
	_, err = os.Stat(mdfilepath)
	if err != nil {
		println("!! markdown file DOES NOT exist")
		os.Exit(4)
	}

	// command
	switch command {
	case "html":
		// write result to file
		md2html_w(mdfilepath, mdfilename_noext+".html", 0644)
	case "docx", "doc", "word":
		var err error
		// convert markdown to html
		err = md2html_w(mdfilepath, mdfilename_noext+".html", 0644)
		if err != nil {
			println("!! Error: ", err.Error())
			os.Exit(6)
		}
		// convert html to docx
		err = html2docx_w(mdfilename_noext+".html", mdfilename_noext+".docx")
		if err != nil {
			println("!! Error: ", err.Error())
			os.Exit(7)
		}
	default:
		help()
	}
}

func help() {
	println("MDC version: ", VERSION)
	helptext := `Manual for Markdown Converter

mdc <COMMAND> [filepath] [-t <template>] [-o <filename>]

[COMMAND]
	pdf
		Convert markdown file to PDF.
	docx | doc | word
		Convert markdown file to word (.docx).
	help
		Show help document.

[OTHERS]
	-t <template> !!NOT SUPPORT!!
		Word output only. Generate .docx with a specific template. 
		Choose a template .docx file path or online reference file 
		name.
`
	fmt.Printf(helptext)
}
