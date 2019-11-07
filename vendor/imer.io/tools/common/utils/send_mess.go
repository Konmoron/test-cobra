package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	messUrl = "https://dc-op.yonyoucloud.com/alarm/v1/message/send"
)

var hostname, _ = os.Hostname()
var hostIP, _ = externalIP()

type NoticeMess struct {
	Mess     string `json:"notice_context"`
	Label    string `json:"label"`
	SendUser string `json:"send_user"`
}

func (c NoticeMess) SendMess() {
	sendTime := time.Now().Format("2006-01-02 15:04:05")
	c.Mess = "主机信息: " + hostname + "-" + hostIP + "\n报警时间: " + sendTime + "\n" + c.Mess
	if c.SendUser == "" {
		c.SendUser = "xiaoyou"
	}
	jsonStr, err := json.Marshal(c)
	//fmt.Println(c)
	if err != nil {
		fmt.Println(err)
	} else {
		req, _ := http.NewRequest("POST", messUrl, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		_, err = client.Do(req)
		//if err != nil {
		//	fmt.Println(err)
		//	wg.Done()
		//} else {
		//	wg.Done()
		//}
		//resp, _ := client.Do(req)
		//if err != nil {
		//	fmt.Println(err)
		//	done <- false
		//}
		//defer resp.Body.Close()
		//body, _ := ioutil.ReadAll(resp.Body)
		//fmt.Println(string(body))
		//wg.Done()
	}
}

// 发送消息
// 调用完，一定要释放done
//	<- done
func (c *NoticeMess) SendMessCh(done chan bool) {
	c.SendMess()
	done <- true
	//// 如果没有定义SendUser，则设置为小友
	//if c.SendUser == "" {
	//	c.SendUser = "xiaoyou"
	//}
	//jsonStr, err := json.Marshal(c)
	//if err != nil {
	//	fmt.Println(err)
	//	done <- false
	//} else {
	//	req, _ := http.NewRequest("POST", messUrl, bytes.NewBuffer(jsonStr))
	//	req.Header.Set("Content-Type", "application/json")
	//	client := &http.Client{}
	//	client.Do(req)
	//	//resp, _ := client.Do(req)
	//	//if err != nil {
	//	//	fmt.Println(err)
	//	//	done <- false
	//	//}
	//	//defer resp.Body.Close()
	//	//body, _ := ioutil.ReadAll(resp.Body)
	//	//fmt.Println(string(body))
	//	done <- true
	//}
}

func (c *NoticeMess) SendMessForLoop(wg *sync.WaitGroup) {
	c.SendMess()
	wg.Done()
	//// 如果没有定义SendUser，则设置为小友
	//if c.SendUser == "" {
	//	c.SendUser = "xiaoyou"
	//}
	//jsonStr, err := json.Marshal(c)
	////fmt.Println(c)
	//if err != nil {
	//	fmt.Println(err)
	//	wg.Done()
	//} else {
	//	req, _ := http.NewRequest("POST", messUrl, bytes.NewBuffer(jsonStr))
	//	req.Header.Set("Content-Type", "application/json")
	//	client := &http.Client{}
	//	_, err = client.Do(req)
	//	wg.Done()
	//	//if err != nil {
	//	//	fmt.Println(err)
	//	//	wg.Done()
	//	//} else {
	//	//	wg.Done()
	//	//}
	//	//resp, _ := client.Do(req)
	//	//if err != nil {
	//	//	fmt.Println(err)
	//	//	done <- false
	//	//}
	//	//defer resp.Body.Close()
	//	//body, _ := ioutil.ReadAll(resp.Body)
	//	//fmt.Println(string(body))
	//	//wg.Done()
	//}
}
