package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func PrintToJson(v interface{}) {
	jsonIndent, err := json.MarshalIndent(&v, "", "    ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(jsonIndent))
}

func PrintRespBody2Json(resp *http.Response) {
	var temp = make(map[string]interface{})
	err := json.NewDecoder(resp.Body).Decode(&temp)
	if err == nil {
		PrintToJson(temp)
	} else {
		fmt.Println(err)
	}
}
