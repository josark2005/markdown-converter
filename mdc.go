package main

import (
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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
// https://cdn.jsdelivr.net/gh/jokin1999/markdown-converter@0.0.1/
const MDC_SERVER_DEFAULT = "https://mdc.josark.com"
const MDC_SERVER_ACC = "https://cdn.jsdelivr.net/gh/jokin1999/markdown-converter"
const MDC_SERVER_DEV = "https://test.mdc.josark.com"

// Pandoc Command
var pandoc_bin string

// MDC Dir
var mdc_dir string

// MDC local configuration file
var mdc_localconf string

// MDC local configuration JSON string
var mdc_conf []byte

// Initializaion
func init() {
	pandoc_bin = "pandoc"
	userhomedir, _ := os.UserHomeDir()
	mdc_dir, _ = filepath.Abs(userhomedir + "/.mdc")
	// check pandoc binary
	_, err := exec.Command(pandoc_bin, "--version").Output()
	if err != nil {
		// try to find pandoc binary in mdc excutable directory
		execdir, _ := os.Executable()
		execdir = filepath.Dir(execdir)
		_, err := os.Stat(execdir + "/pandoc.exe")
		if err != nil {
			println("!! Pandoc binary is required.")
		} else {
			pandoc_bin = execdir + "/pandoc.exe"
		}

	}
	// check MDC DIR
	localdir, err := os.Stat(mdc_dir)
	if err != nil || !localdir.IsDir() {
		println("!! MDC DIR is unreachable. Try to create", mdc_dir)
		err := os.Mkdir(mdc_dir, 0644)
		if err != nil {
			println("!! Failed to create directory", mdc_dir)
		}
	}
	// check local configuration
	mdc_localconf, _ = filepath.Abs(mdc_dir + "/" + VERNAME + ".json")
	// try to read local configuration
	mdc_conf, err = os.ReadFile(mdc_localconf)
	if err != nil {
		println("!! Cannot read local configuration:", mdc_localconf)
	}
}

func md2html(md string) string {
	println("Running: ", pandoc_bin, "-f markdown -t html", md)
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
	println("Running: ", pandoc_bin, "-f html -t docx", html, "-o", filename)
	_, err := exec.Command(pandoc_bin, "-f", "markdown", "-t", "docx", html, "-o", filename).Output()
	if err != nil {
		fmt.Println("!! Convert failed: ", err)
	}
	return err
}

func main() {
	args := os.Args

	// at least one param should be given
	if len(args) < 2 {
		println("!! Missing arguments")
		os.Exit(0)
	}

	// command or target
	command := args[1]

	switch command {
	case "builtin":
		builtin()
	case "update":
		// Update server json
		fs_download := flag.NewFlagSet("update", flag.ExitOnError)
		// support dev/acc mode
		server := fs_download.String("server", "default", "Specify a server. 'acc' and 'dev' are available. \nThe 'acc' server may not have latest files, e.g. templates. \nDO NOT CHOOSE 'dev' SERVER IF YOU DO NOT KNOW WHAT WILL HAPPED.")
		fs_download.Parse(args[2:])
		switch *server {
		case "default":
			conf, err := mdcass.FetchConf(MDC_SERVER_DEFAULT, VERNAME, CONFVER)
			if err != nil {
				println("!! Failed to fetch configuration files")
			}
			err = mdcass.Write2File(mdc_localconf, conf)
			if err != nil {
				println("!! Failed to save configuration file:", mdc_localconf)
			}
		case "acc":
			println(">> ATTENTION: Accelerated server may not have latest files")
			_ver := strings.Split(VERSION, ".")
			_ver = _ver[:2]
			ver := strings.Join(_ver, ".")
			conf, err := mdcass.FetchConf(MDC_SERVER_ACC+"@"+ver, VERNAME, CONFVER)
			if err != nil {
				println("!! Failed to fetch configuration files")
			}
			err = mdcass.Write2File(mdc_localconf, conf)
			if err != nil {
				println("!! Failed to save configuration file:", mdc_localconf)
			}
		case "dev":
			println(">> ATTENTION: Developing server")
			println(">> ATTENTION: Developing server")
			conf, err := mdcass.FetchConf(MDC_SERVER_DEV, VERNAME, CONFVER)
			if err != nil {
				println("!! Failed to fetch configuration files")
			}
			err = mdcass.Write2File(mdc_localconf, conf)
			if err != nil {
				println("!! Failed to save configuration file:", mdc_localconf)
			}
		default:
			println("!! ERROR: Non-registered server")
		}
	default:
		// generate filepath
		mdfilepath, _ := filepath.Abs(args[2])
		mdfilename := filepath.Base(mdfilepath)

		// remove the extension name
		_filename := strings.Split(mdfilename, ".")
		_filename = _filename[:len(_filename)-1]
		mdfilename_noext := strings.Join(_filename, ".")

		var err error

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
			// help()
		}
	}
	os.Exit(0)

}

func builtin() {
	println()
	println("======================== Built-in Information ========================")
	println("MDC Version    : ", VERSION)
	println("MDC Vername    : ", VERNAME)
	println("MDC Server     : ", MDC_SERVER_DEFAULT)
	println("MDC Server ACC : ", MDC_SERVER_ACC)
	println("MDC Server DEV : ", MDC_SERVER_DEV)
	// println("Mirrors: ")
	// for k, v := range MDC_SERVERS {
	// 	println("    ", k, ":", v)
	// }
}
