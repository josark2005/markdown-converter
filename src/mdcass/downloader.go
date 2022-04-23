package mdcass

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/buger/jsonparser"
)

/**
 * @Funcion Download
 * @param server Server URL for fetching MDC server configurations
 * @param pandoc Pandoc binary filename
 * @param version Current MDC version
 * @param vername Current MDC version name
 * @param confver Current MDC accepted server configration file version
 * @param filepath The location to save the file
 */
func Download(server string, pandoc string, version string, vername string, confver int, filepath string) {
	conf, err := fetchConf(server+"/deck.json", vername, confver)
	if err != nil {
		println(err.Error())
	}
	println(conf)
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
	println(resp.StatusCode)
	if resp.StatusCode == 200 {
		// parse the json data
		body, _ := ioutil.ReadAll(resp.Body)
		_confver, _ := jsonparser.GetInt(body, "conf_ver")
		if int(_confver) != confver {
			println("!! Error CONFVER: ", confver, ", ", confver, " is expected")
			os.Exit(1)
		}
		// if the vername is acceptable
		_vername, _, _, _ := jsonparser.Get(body, "vername")
		if string(_vername) != vername {
			println("!! Error CONFVER: ", confver, ", ", confver, " is expected")
			os.Exit(1)
		}

	} else {
		return nil, err
	}
	return nil, nil

}
