package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/jokin1999/markdown-converter/src/mdcass"
)

// Version
const VERSION = "0.0.2"

// Version Name
const VERNAME = "deck"

// Configuraion File Version
const CONFVER = 1

// MDC Server
const MDC_SERVER_DEFAULT = "https://mdc.josark.com"
const MDC_SERVER_PANDOC = "https://pandoc.mdc.josark.com"

var pandoc_bin string
var mdc_conf string

// Initializaion
func init() {
	// check pandoc binary
	userhomedir, _ := os.UserHomeDir()
	mdcdir := userhomedir + "/.mdc"
	pandoc_bin, _ = filepath.Abs(mdcdir + "/pandoc")
	if runtime.GOOS == "windows" {
		pandoc_bin += ".exe"
	}
	_, err := os.Stat(pandoc_bin)
	if err != nil || os.IsNotExist(err) {
		println("!! Pandoc binary does NOT exist or is NOT permitted to access")
		println("!! Use mdc download to obtain a binary automatically")
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

	// show built-in information
	if command == "builtin" {
		builtin()
		os.Exit(0)
	}

	// download pandoc binary
	if command == "download" {
		mdcass.Download(MDC_SERVER_DEFAULT, MDC_SERVER_PANDOC, VERSION, VERNAME, CONFVER, pandoc_bin)
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
	builtin
		Show built-in information.
	download [Download Flags]
		Download pandoc binary from the Internet.

[Flags]
	-t <template> !!NOT SUPPORT!!
		Word output only. Generate .docx with a specific 
		template. Choose a template .docx file path or online 
		reference file name.

[Download Flags]
	-m [official mirror site name]
		Specify a mirror site name for downloading pandoc 
		binary.
		If put the mirror site empty, mdc will download from 
		a default built-in mirror site. ONLY official mirror
		site name is accept. e.g. default
	-cm [customized mirror site]
		Specify a mirror site for downloading pandoc binary.
		e.g. https://mirror.mdc.sample.com
`
	fmt.Println(helptext)
}

func builtin() {
	println()
	println("======================== Built-in Information ========================")
	println("MDC Version: ", VERSION)
	println("MDC Vername: ", VERNAME)
	println("MDC Server : ", MDC_SERVER_DEFAULT)
	println("MDC PANDOC : ", MDC_SERVER_PANDOC)
	// println("Mirrors: ")
	// for k, v := range MDC_SERVERS {
	// 	println("    ", k, ":", v)
	// }
}
