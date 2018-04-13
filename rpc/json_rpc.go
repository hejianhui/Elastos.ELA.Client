package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastos/Elastos.ELA.Client/common/config"
	"io/ioutil"
	"net/http"
	"strings"
)

type Response struct {
	Code   int         `json:"code"`
	Result interface{} `json:"result"`
}

var Url string

func GetCurrentHeight() (uint32, error) {
	result, err := CallAndUnmarshal("getcurrentheight", nil)
	if err != nil {
		return 0, err
	}
	return uint32(result.(float64)), nil
}

func GetBlockByHeight(height uint32) (*BlockInfo, error) {
	resp, err := CallAndUnmarshal("getblockbyheight", Param("height", height))
	if err != nil {
		return nil, err
	}
	block := &BlockInfo{}
	unmarshal(&resp, block)

	return block, nil
}

func Call(method string, params map[string]string) ([]byte, error) {
	if Url == "" {
		return nil, errors.New("Unknown rpc url.")
	}
	data, err := json.Marshal(map[string]interface{}{
		"method": method,
		"params": params,
	})
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(Url, "application/json", strings.NewReader(string(data)))
	if err != nil {
		fmt.Printf("POST requset: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func CallAndUnmarshal(method string, params map[string]string) (interface{}, error) {
	body, err := Call(method, params)
	if err != nil {
		return nil, err
	}

	resp := Response{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return string(body), nil
	}

	if resp.Code != 0 {
		return nil, errors.New(fmt.Sprint(resp.Result))
	}

	return resp.Result, nil
}

func unmarshal(result interface{}, target interface{}) error {
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, target)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	Url = config.Params().Host
}
