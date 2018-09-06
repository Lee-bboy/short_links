package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type Controller struct {
}

func (c *Controller) parseParams(r *http.Request) (map[string]interface{}, error) {
	//解析请求参数
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.New("请求不正确")
	}

	var data map[string]interface{}
	err = json.Unmarshal(raw, &data)
	if err != nil {
		return nil, errors.New("请求不正确")
	}

	return data, nil
}

func (c *Controller) errorResponse(w http.ResponseWriter, code int, message string) {
	result, _ := json.Marshal(map[string]interface{}{"errcode": code, "errmsg": message})
	w.Write(result)
}

func (c *Controller) jsonResponse(w http.ResponseWriter, data interface{}) {
	resp, _ := json.Marshal(data)
	w.Write(resp)
}
