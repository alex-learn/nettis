package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func LoadJsonFile(jsonFile string, verbose bool) (map[string]interface{}, error) {
		var m map[string]interface{}
        rawJson, err := loadFile(jsonFile, verbose)
        if err != nil {
                return m, err
        }
        err = json.Unmarshal(rawJson, &m)
		return m, err
}

func loadFile(jsonFile string, verbose bool) ([]byte, error) {
        file, err := ioutil.ReadFile(jsonFile)
        if err != nil {
                if os.IsNotExist(err) {
                        if verbose { //not a problem.
                                log.Printf("%s not found", jsonFile)
                        }
                } else {
                        //always log because it's unexpected
                        log.Printf("File error: %v", err)
                }
        } else {
                if verbose {
                        log.Printf("Found %s", jsonFile)
                }
        }
        return file, err
}


func ToStringSlice(v interface{}, k string) ([]string, error) {
        ret := []string{}
        switch typedV := v.(type) {
        case []interface{}:
                for _, i := range typedV {
                        ret = append(ret, i.(string))
                }
                return ret, nil
        }
        return ret, fmt.Errorf("%s should be a json array, not a %T", k, v)
}

func ToString(v interface{}, k string) (string, error) {
        switch typedV := v.(type) {
        case string:
                return typedV, nil
        }
        return "", fmt.Errorf("%s should be a json string, not a %T", k, v)
}

func ToMapStringString(v interface{}, k string) (map[string]string, error) {
		msi, err := ToMap(v, k)
		if err != nil {
				return nil, err
		}
		ret := map[string]string{}
		for k2, v2 := range msi {
			v2s, err := ToString(v2, k2)
			if err != nil {
				return nil, fmt.Errorf("%s is not a string", k2)
			}
			ret[k2] = v2s
		}
		return ret, nil
}

func ToMap(v interface{}, k string) (map[string]interface{}, error) {
        switch typedV := v.(type) {
        case map[string]interface{}:
                return typedV, nil
        }
        return nil, fmt.Errorf("%s should be a json map, not a %T", k, v)
}