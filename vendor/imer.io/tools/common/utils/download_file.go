package utils

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func DownloadFile(filepath, downloadUrl string) error {
	// 参考：
	//  https://cloud.tencent.com/developer/ask/36815
	//  https://blog.csdn.net/fyxichen/article/details/46915285

	timeout := time.Duration(10 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	// Get the data
	resp, err := client.Get(downloadUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New("获取downloadUrl失败, 异常状态码: " + strconv.Itoa(resp.StatusCode) + ", downloadUrl: " + downloadUrl)
	}

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
