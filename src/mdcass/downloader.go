package mdcass

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/tidwall/gjson"
)

/**
 * @Funcion Download
 * @param server Server URL for fetching MDC server configurations
 * @param pandoc Pandoc binary filename
 * @param version Current MDC version
 * @param vername Current MDC version name
 * @param confver Current MDC accepted server configration file version
 * @param dfilepath The location to save the file
 */
func Download(server string, pandoc string, version string, vername string, confver int, dfilepath string) {
	conf, err := fetchConf(server+"/deck.json", vername, confver)
	if err != nil {
		println(err.Error())
	}
	// read download url
	var osname string
	switch runtime.GOOS {
	case "windows":
		osname = "win"
	case "linux":
		osname = "linux"
	case "darwin":
		osname = "darwin"
	}
	filename := gjson.Get(string(conf), "pandoc"+"."+vername+"."+osname)
	url := pandoc + "/" + filename.String()
	fetchBig(url, filepath.Dir(dfilepath)+"/"+filename.String(), true)
}

/**
 * @Function fetchConf
 * @param server Server URL for fetching MDC server configurations
 * @param vername Current MDC version name
 * @param confver Current MDC accepted server configration file version
 */

func fetchConf(server string, vername string, confver int) ([]byte, error) {
	println("-> Fetching from ", server)
	resp, err := http.Get(server)
	if err != nil {
		println("!! Error occurd when fetching configurations")
		return nil, err
	}
	defer resp.Body.Close()
	println("-> Fetch succeeded")
	if resp.StatusCode == 200 {
		// parse the json data
		body, _ := ioutil.ReadAll(resp.Body)
		_confver := gjson.Get(string(body), "conf_ver")
		if int(_confver.Int()) != confver {
			println("!! Error CONFVER: ", confver, ", ", confver, " is expected")
			os.Exit(1)
		}
		// if the vername is acceptable
		_vername := gjson.Get(string(body), "vername")
		if _vername.String() != vername {
			println("!! Error CONFVER: ", confver, ", ", confver, " is expected")
			os.Exit(1)
		}
		return body, nil
	} else {
		return nil, err
	}

}

// Download counter
type progress struct {
	total   int64
	current int64
}

// Write function implemention of download counter
func (pg *progress) Write(b []byte) (int, error) {
	n := len(b)
	pg.current += int64(n)
	fmt.Print("\r    ", pg.current/1024, " / ", pg.total/1024, " KB")
	return n, nil
}

func fetchBig(url string, target string, zipped bool) {
	tmpfile, _ := filepath.Abs(target + ".tmp")
	println("-> Downloading from", url)
	println("   * Creating temporary file", tmpfile)
	f_tmpfile, err := os.Create(tmpfile)
	if err != nil {
		println("!! Failed to create temporary file")
		os.Exit(7)
	}
	resp, err := http.Get(url)
	if err != nil {
		println("!! Failed to download file:", err.Error())
		os.Exit(8)
	}
	defer resp.Body.Close()
	// 初始化计数器
	f_writer := &progress{total: resp.ContentLength}
	_, err = io.Copy(f_tmpfile, io.TeeReader(resp.Body, f_writer))
	if err != nil {
		println("!! Failed to download file")
		os.Exit(8)
	}
}
