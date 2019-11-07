package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
)

func ExecCommand(hostIP string, command string, group *sync.WaitGroup) {
	//fmt.Println(url1, command)
	//use http.Post
	//resp, err := http.Post(url1,
	//	"application/x-www-form-urlencoded", strings.NewReader("command=ls"))

	url1 := fmt.Sprintf("http://%s:19880/runner/getresponse", hostIP)
	// use http.PostForm
	resp, err := http.PostForm(url1,
		url.Values{"command": {command}})
	CheckErr(err)
	defer group.Done()
	defer RecoverName()

	//defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		CheckErr(err)
		fmt.Println(hostIP + "\n" + string(body))
	}
}
