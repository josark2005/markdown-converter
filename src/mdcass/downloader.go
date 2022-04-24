package mdcass

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/tidwall/gjson"
)

/**
 * @Function fetchConf
 * @param server Server URL for fetching MDC server configurations
 * @param vername Current MDC version name
 * @param confver Current MDC accepted server configration file version
 */

func FetchConf(server string, vername string, confver int) ([]byte, error) {
	server += "/" + vername + ".json"
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
	} else if resp.StatusCode == 404 {
		println("!! Error: 404 not found.")
		os.Exit(8)
		return nil, err
	} else {
		return nil, err
	}
}

// Write
func Write2File(filepath string, data []byte) error {
	err := os.WriteFile(filepath, data, 0644)
	return err
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
