package utils

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	LastOffSetFileSuffix = ".last_off_set_file.txt"
	Kb1                  = 1024       // 1K
	Mb1                  = 1024 * Kb1 // 1M
	FirstReadByte        = 10 * Kb1   // 10 Kb
	MaxReadByte          = 256 * Mb1  // 100M
	FirstReadFileFlag    = -1
)

type MonitorCount struct {
	MonitorString string `yaml:"monitor_string" json:"monitor_string"`
	MaxNum        int    `yaml:"max_num" json:"max_num"`
	CurrNum       int    `yaml:"curr_num" json:"curr_num"`
	Describe      string `yaml:"describe" json:"describe"`
}

type MonitorApp struct {
	AppName             string         `yaml:"app_name"`
	LogFiles            []string       `yaml:"log_files"`
	PreMatchTimeFormant string         `yaml:"pre_match_time_format"`
	PreMust             []string       `yaml:"pre_must"`
	PreMustNot          []string       `yaml:"pre_must_not"`
	Must                []string       `yaml:"must"`
	MustNot             []string       `yaml:"must_not"`
	Counts              []MonitorCount `yaml:"counts"`
}

type MonitorAppConfig struct {
	Label                     string                 `yaml:"label"`
	RunModel                  string                 `yaml:"run_model"`
	MaxSendLine               int                    `yaml:"max_send_line"`
	IsDebug                   bool                   `yaml:"debug"`
	DefaultMaxNum             int                    `yaml:"default_max_num"`
	DefaultPreMatchTimeFormat string                 `yaml:"default_pre_match_time_format"`
	MonitorApps               map[string]*MonitorApp `yaml:"monitor_apps"`
}

func GetLastOffSet(logFile, lastOffSetFile string) (int64, bool) {
	var lastOffSet int64 = 0
	var fileSize int64 = 0
	var isChange = true
	if fileInfo, err := os.Stat(logFile); err == nil {
		fileSize = fileInfo.Size()
	} else {
		log.Println(err)
		isChange = false
		return 0, isChange
	}

	if IsFileExists(lastOffSetFile) {
		if f, err := os.Open(lastOffSetFile); err == nil {
			defer f.Close()
			b := make([]byte, 1024)
			if n, err := f.Read(b); err == nil {
				// 先将string(b[:n]里面的\n删除
				// 然后将其转化为int64
				if lastOffSet, err = strconv.ParseInt(strings.Replace(string(b[:n]), "\n", "", -1), 10, 64); err != nil {
					lastOffSet = 0
					log.Println(err)
				}
			}
		} else {
			log.Println(err)
			lastOffSet = 0
		}
	} else {
		lastOffSet = FirstReadFileFlag
		SetLastOffSet(fileSize, lastOffSetFile)
		isChange = false
		return lastOffSet, isChange
	}

	//log.Println(lastOffSet)
	//log.Println("file size: ", fileSize)
	if lastOffSet == 0 {
		if fileSize > FirstReadByte {
			lastOffSet = fileSize - FirstReadByte
		} else if fileSize > 0 {
			lastOffSet = fileSize
		} else {
			isChange = false
		}
	} else {
		if lastOffSet < fileSize {
			if fileSize-lastOffSet > MaxReadByte {
				lastOffSet = fileSize - MaxReadByte
			}
		} else if lastOffSet > fileSize {
			//fmt.Println(fileSize-lastOffSet, MaxReadByte)
			if fileSize > FirstReadByte {
				lastOffSet = fileSize - FirstReadByte
			} else {
				lastOffSet = fileSize
			}
		} else {
			isChange = false
		}
	}

	//fmt.Println("start position:", lastOffSet, ", is change:", isChange)
	return lastOffSet, isChange
}

func SetLastOffSet(currOffSet int64, lastOffSetFile string) {
	if f, err := os.Create(lastOffSetFile); err == nil {
		defer f.Close()
		//fmt.Println("end position:", currOffSet)
		if _, err = f.WriteString(strconv.FormatInt(currOffSet, 10)); err != nil {
			log.Println(err)
		}
	}
}

func IsMatch(str string, must, mustNot []string) bool {
	var isMatch = true
	if len(must) != 0 && must[0] != "" {
		for _, s := range must {
			isMatch = isMatch && strings.Contains(str, s)
			if !isMatch {
				return isMatch
			}
		}
	}

	if len(mustNot) != 0 && mustNot[0] != "" {
		for _, s := range mustNot {
			isMatch = isMatch && !(strings.Contains(str, s))
			if !isMatch {
				return isMatch
			}
		}
	}

	return isMatch
}

func (monitorApp *MonitorApp) IsPreMatch(str string) bool {
	var isMatch = true

	if monitorApp.PreMatchTimeFormant != "" {
		preMatchTimeStr := time.Now().Format(monitorApp.PreMatchTimeFormant)
		if !(len(str) >= len(preMatchTimeStr) && str[0:len(preMatchTimeStr)] == preMatchTimeStr) {
			isMatch = false
			return isMatch
		}
	}

	isMatch = isMatch && IsMatch(str, monitorApp.PreMust, monitorApp.PreMustNot)

	return isMatch
}

func (monitorApp *MonitorApp) SendLogMess(wg *sync.WaitGroup, conf *MonitorAppConfig, logFile, noticeContext string) {
	defer wg.Done()
	// 判断是否满足must匹配和must_not匹配
	//var isMatch = IsMatch(noticeContext, monitorApp.Must, monitorApp.MustNot)
	// isMatch

	//if len(monitorApp.Must) != 0 && monitorApp.Must[0] != "" {
	//	for _, s := range monitorApp.Must {
	//		isMatch = isMatch && strings.Contains(noticeContext, s)
	//		if !isMatch {
	//			return
	//		}
	//	}
	//}
	//// 判断是否满足must_not匹配
	//if len(monitorApp.MustNot) != 0 && monitorApp.MustNot[0] != "" {
	//	for _, s := range monitorApp.MustNot {
	//		isMatch = isMatch && (!strings.Contains(noticeContext, s))
	//		if !isMatch {
	//			return
	//		}
	//	}
	//}

	// 根据isMatch是否发送报警
	if IsMatch(noticeContext, monitorApp.Must, monitorApp.MustNot) {
		// 生成报警信息
		noticeContext =
			"报警原因: " + monitorApp.AppName + "出现异常日志" +
				"\n应用名称: " + monitorApp.AppName +
				"\n日志文件: " + logFile +
				"\n日志信息: " + noticeContext +
				"\n详细日志请登陆服务器查看"
		// 如果是DEBUG，则不发送报警，打印报警信息
		if !conf.IsDebug {
			noticeMess := NoticeMess{Mess: noticeContext, Label: conf.Label, SendUser: ""}
			noticeMess.SendMess()
		} else {
			log.Println("is match:", true)
			fmt.Printf("\nnoticeContext: \n%s\n\n", noticeContext)
		}
	} else {
		// 如果是DEBUG，打印报警信息
		if conf.IsDebug {
			log.Println("is match:", true)
			fmt.Printf("\nnoticeContext: \n%s\n\n", noticeContext)
		}
	}

	//wg.Done()
}

func (monitorApp *MonitorApp) AddCounts(noticeContext string) {
	if len(monitorApp.Counts) > 0 {
		for i := range monitorApp.Counts {
			if monitorApp.Counts[i].MonitorString == "__ALL__" {
				monitorApp.Counts[i].CurrNum += 1
			} else {
				if strings.Contains(noticeContext, monitorApp.Counts[i].MonitorString) {
					monitorApp.Counts[i].CurrNum += 1
				}
			}
			//fmt.Println(monitorApp.Counts[i].MonitorString, monitorApp.Counts[i].MaxNum, monitorApp.Counts[i].CurrNum)
		}
	}
}

func (monitorApp *MonitorApp) CheckCounts(conf *MonitorAppConfig, logFile string) {
	if len(monitorApp.Counts) > 0 {
		for _, c := range monitorApp.Counts {
			if c.CurrNum > c.MaxNum {
				var noticeContext = ""
				if c.MonitorString == "__ALL__" {
					noticeContext = "报警原因: " + monitorApp.AppName + "异常日志数量大于" + strconv.Itoa(c.MaxNum) +
						"\n应用名称: " + monitorApp.AppName +
						"\n日志文件: " + logFile +
						"\n异常日志数量: " + strconv.Itoa(c.CurrNum) +
						"\n详细日志请登陆服务器查看"
				} else {
					noticeContext += "报警原因: 【" + monitorApp.AppName + "】异常日志中包含【" + c.MonitorString + "】的数量大于" + strconv.Itoa(c.MaxNum)
					if c.Describe != "" && c.Describe != c.MonitorString {
						noticeContext += "\n报警描述: " + c.Describe
					}
					//noticeContext += "报警原因: 【" + monitorApp.AppName + "】异常日志中包含【" + c.MonitorString + "】的数量大于" + strconv.Itoa(c.MaxNum) +
					//	//"\n应用名称: " + monitorApp.AppName +
					//	"\n日志文件: " + logFile +
					//	//"\n异常字段: " + c.MonitorString +
					//	"\n异常日志数量: " + strconv.Itoa(c.CurrNum) +
					//	"\n详细日志请登陆服务器查看"
					//noticeContext += "报警原因: 【" + monitorApp.AppName + "】异常日志中包含【" + c.MonitorString + "】的数量大于" + strconv.Itoa(c.MaxNum) +
					noticeContext += "\n日志文件: " + logFile +
						"\n异常日志数量: " + strconv.Itoa(c.CurrNum) +
						"\n详细日志请登陆服务器查看"
				}

				noticeMess := NoticeMess{Mess: noticeContext, Label: conf.Label, SendUser: ""}
				if conf.IsDebug {
					fmt.Printf("\nnoticeContext: \n%s\n\n", noticeContext)
				} else {
					noticeMess.SendMess()
				}
			}
		}
	}
}

func (monitorApp *MonitorApp) CheckLogFile(conf *MonitorAppConfig, workDir, logFile string) {
	if f, err := os.Open(logFile); err == nil {
		defer f.Close()

		// References:
		//	1.https://stackoverflow.com/questions/34654514/how-to-read-a-file-starting-from-a-specific-line-number-using-scanner
		//	2.https://juejin.im/post/5bd7b1b7e51d4547f763fe79
		//	3.https://stackoverflow.com/questions/17863821/how-to-read-last-lines-from-a-big-file-with-go-every-10-secs
		r := bufio.NewReader(f)
		var line []byte

		//preMatch := time.Now().Format(monitorApp.PreMatchTimeFormant)

		// 生成OffSet文件名
		//
		// 将/tmp/error.log转换为tmp_error.log
		// 	strings.Replace(logFile, "/", "", 1): 去掉第一个/
		// 	strings.Replace(strings.Replace(logFile, "/", "", 1), "/", "_", -1): 替换其余的/为_
		logFileName := strings.Replace(strings.Replace(logFile, "/", "", 1), "/", "_", -1)
		// 生成lastOffSetFile
		lastOffSetFile := workDir + "/." + logFileName + LastOffSetFileSuffix

		// 获取OffSet和文件是否改变
		lastOffSet, isChange := GetLastOffSet(logFile, lastOffSetFile)
		fPos := lastOffSet

		if isChange {
			if _, err = f.Seek(fPos, 0); err == nil {
				var wg sync.WaitGroup
				var noticeContext = ""
				var currLineMatch = 0
				var errorLogFlag = false
				for i := 1; ; i++ {
					if line, err = r.ReadBytes('\n'); err == nil {
						lineLen := len(line)
						if lineLen != 0 {
							fPos += int64(lineLen)
							// 将byte转换为string
							lineStr := string(line[:lineLen])

							if monitorApp.IsPreMatch(lineStr) {
								if currLineMatch > 0 {
									wg.Add(1)
									go monitorApp.SendLogMess(&wg, conf, logFile, noticeContext)
									monitorApp.AddCounts(noticeContext)
									currLineMatch = 0
								}
								currLineMatch += 1
								errorLogFlag = true
								noticeContext = lineStr
								continue
							}
							//log.Println(errorLogFlag, lineStr)
							if errorLogFlag {
								currLineMatch += 1
								noticeContext += lineStr
							}

							if currLineMatch >= conf.MaxSendLine {
								wg.Add(1)
								go monitorApp.SendLogMess(&wg, conf, logFile, noticeContext)
								monitorApp.AddCounts(noticeContext)
								currLineMatch = 0
								errorLogFlag = false
								noticeContext = ""
							}
						}
					} else {
						// 如果最后两行的日志都需报警，那么，最后一行的报警需要在此发送
						if currLineMatch > 0 {
							wg.Add(1)
							go monitorApp.SendLogMess(&wg, conf, logFile, noticeContext)
							monitorApp.AddCounts(noticeContext)
							currLineMatch = 0
							errorLogFlag = false
							noticeContext = ""
						}
						break
					}
				}

				if err != io.EOF {
					log.Fatal(err)
				}
				wg.Wait()
				//defer SetLastOffSet(fPos, lastOffSetFile)
				monitorApp.CheckCounts(conf, logFile)
				log.Println("log file:", logFile, ", start position:", lastOffSet, ", end position:", fPos, ", log file is change:", isChange)
				defer SetLastOffSet(fPos, lastOffSetFile)
			} else {
				log.Fatal(err)
			}
		} else {
			if lastOffSet == FirstReadFileFlag {
				log.Println(logFile, "first read.")
			} else {
				log.Println(logFile, "not change.")
			}
		}
	} else {
		log.Fatal(err)
	}
}
