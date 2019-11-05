package main

import (
	"fmt"
	"yonyou.com/iuap/tools/common/utils"

	"github.com/spf13/viper"
)

type MonitorResponseTime struct {
	MaxNum          int     `yaml:"max_num" mapstructure:"max_num"`
	CurrNum         int     `yaml:"curr_num" mapstructure:"curr_num"`
	Label           string  `yaml:"label" mapstructure:"label"`
	MaxResponseTime float64 `yaml:"max_response_time" mapstructure:"max_response_time"`
	MatchString     string  `yaml:"match_string" mapstructure:"match_string"`
}

type MonitorResponseStatus struct {
	MaxNum      int    `yaml:"max_num" mapstructure:"max_num"`
	CurrNum     int    `yaml:"curr_num" mapstructure:"curr_num"`
	Label       string `yaml:"label" mapstructure:"label"`
	Status      string `yaml:"status" mapstructure:"status"`
	MatchString string `yaml:"match_string" mapstructure:"match_string"`
}

type MonitorCount struct {
	MaxNum      int    `yaml:"max_num" mapstructure:"max_num"`
	CurrNum     int    `yaml:"curr_num" mapstructure:"curr_num"`
	Label       string `yaml:"label" mapstructure:"label"`
	MatchString string `yaml:"match_string" mapstructure:"match_string"`
}

type MonitorUrl struct {
	/*
		主机信息: iZ25hjygnzaZ-10.3.7.130
		报警时间: 2019-08-16 11:26:01
		报警原因: 【DescribeTitle】出现异常日志
		报警描述: 【31】条请求【DescribeName】的【响应时间】超过【2】秒
		监控字段: PreMatchMust
		日志信息(仅显示1条): 10.3.5.52|-|[16/Aug/2019:11:25:02 +0800]|"POST /apptenant/rest/tenantrest/userIdsInTenant?ts=1565925900550&username=iform&appId=873925c8d9f81e6204f0a9f85e80887a39ab21db HTTP/1.1"|200|900|apcenter.yonyoucloud.com|"-"|"Apache-HttpClient/4.4.1 (Java/1.8.0_151)"|10.3.37.244, 100.120.36.204|-|"text/plain;charset=UTF-8"|2.037|2.034|3079|841|"10.3.15.8:80"|"http"


		开发者中心值班表: https://wiki.yonyoucloud.com/pages/viewpage.action?pageId=8553406
	*/
	DescribeTitle            string                   `yaml:"describe_title" mapstructure:"describe_title"`
	DescribeName             string                   `yaml:"describe_name" mapstructure:"describe_name"`
	PreMatchMust             []string                 `yaml:"pre_match_must" mapstructure:"pre_match_must"`
	PreMatchMustNot          []string                 `yaml:"pre_match_must_not" mapstructure:"pre_match_must_not"`
	NotMonitorResponseTime   bool                     `yaml:"not_monitor_response_time" mapstructure:"not_monitor_response_time"`
	ResponseTime             []*MonitorResponseTime   `yaml:"response_time" mapstructure:"response_time"`
	NotMonitorResponseStatus bool                     `yaml:"not_monitor_response_status" mapstructure:"not_monitor_response_status"`
	ResponseStatus           []*MonitorResponseStatus `yaml:"response_status" mapstructure:"response_status"`
	NotMonitorCount          bool                     `yaml:"not_monitor_count" mapstructure:"not_monitor_count"`
	Count                    []*MonitorCount          `yaml:"count" mapstructure:"count"`
}

type MonitorUrlConfig struct {
	DefaultDescribeTitle         string         `yaml:"default_describe_title" mapstructure:"default_describe_title"`
	DefaultLabel                 string         `yaml:"default_label" mapstructure:"default_label"`
	RunModel                     string         `yaml:"run_model" mapstructure:"run_model"`
	ConcurrentReadLineNumber     int            `yaml:"concurrent_read_line_number" mapstructure:"concurrent_read_line_number"`
	IsDebug                      bool           `yaml:"debug" mapstructure:"debug"`
	DefaultMaxResponseTime       float64        `yaml:"default_max_response_time" mapstructure:"default_max_response_time"`
	DefaultMonitorResponseStatus string         `yaml:"default_monitor_response_status" mapstructure:"default_monitor_response_status"`
	DefaultMaxNum                int            `yaml:"default_max_num" mapstructure:"default_max_num"`
	LogFileConfs                 []*LogFileConf `yaml:"log_file_confs" mapstructure:"log_file_confs"`
	//LogFiles                     []string      `yaml:"log_files" mapstructure:"log_files"`
	//MonitorUrlS                  []*MonitorUrl  `yaml:"monitor_urls" mapstructure:"monitor_urls"`
}

type LogFileConf struct {
	LogFiles    []string      `yaml:"log_files" mapstructure:"log_files"`
	MonitorUrlS []*MonitorUrl `yaml:"monitor_urls" mapstructure:"monitor_urls"`
}


// Create private data struct to hold config options.
type config struct {
	Hostname string `yaml:"hostname"`
	Port     string `yaml:"port"`
}

// Create a new config instance.
var (
	conf *MonitorUrlConfig
)

// Read the config file from the current directory and marshal
// into the conf config struct.
func getConf() *MonitorUrlConfig {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()

	if err != nil {
		fmt.Printf("%v", err)
	}
	viper.Set("concurrent_read_line_number", 10)
	conf := &MonitorUrlConfig{}
	err = viper.Unmarshal(conf)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}

	return conf
}

// Initialization routine.
func init() {
	// Retrieve config options.
	conf = getConf()
}

// Main program.
func main() {

	utils.PrintToJson(conf)
}