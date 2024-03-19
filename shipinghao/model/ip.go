package model

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

var targetUrl = "https://channels.weixin.qq.com"
var ipUrl = "http://webapi.http.zhimacangku.com/getip?neek=321a408a&num=1&type=2&pro=0&city=0&ts=1&yys=0&port=1&pack=341617&ts=0&ys=0&cs=0&lb=1&sb=&pb=4&mr=1&regions="

// {"code":0,"data":[{"ip":"110.90.15.23","port":4238,"expire_time":"2024-03-04 18:03:43"}],"msg":"0","success":true}
func IsIp(ip string) (b bool) {

	if m, _ := regexp.MatchString(`^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\:([0-5]?[0-9]{1,4}|6[0-5]{2}[0-3][0-5])$`, ip); !m {
		return false
	}
	return true
}

func getip(url string) string {
	fmt.Println("&*&*&*&*&*")
	client := &http.Client{}
	rqt, err := http.NewRequest("GET", url, nil)
	if err != nil {
		println("获取IP失败!")
		return ""
	}
	rqt.Header.Add("User-Agent", "Lingjiang")
	//处理返回结果
	response, _ := client.Do(rqt)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ""
	}
	ip := string(body)
	if c := IsIp(ip); !c {
		println("获取IP失败!")
		return ""
	}
	fmt.Println("ip:", ip)
	return "http://" + ip
}

// sock5代理
func socksproxy() {
	ip := getip(ipUrl)
	urli := url.URL{}
	urlproxy, _ := urli.Parse(ip)
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(urlproxy),
		},
	}
	rqt, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		println("接口获取IP失败!")
		return
	}

	rqt.Header.Add("User-Agent", "Lingjiang")
	//处理返回结果
	response, _ := client.Do(rqt)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	fmt.Println("socks5:", string(body))
	return

}

// http代理
func httpproxy() {
	ip := getip(ipUrl)
	fmt.Println("********************************")
	urli := url.URL{}
	urlproxy, _ := urli.Parse(ip)
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(urlproxy),
		},
	}
	fmt.Println("----------------------------------------------------------------")
	rqt, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		println("接口获取IP失败!")
		return
	}
	rqt.Header.Add("User-Agent", "Lingjiang")
	//处理返回结果
	response, _ := client.Do(rqt)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	fmt.Println("http:", string(body))
	return
}

// 本机IP
func httplocal() {
	client := &http.Client{}
	rqt, err := http.NewRequest("GET", "http://myip.top", nil)
	if err != nil {
		println("接口获取IP失败!")
		return
	}
	rqt.Header.Add("User-Agent", "Lingjiang")
	//处理返回结果
	response, _ := client.Do(rqt)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	fmt.Println("本机:", string(body))
	return
}

func GetIp() {
	httplocal()
	httpproxy()
	time.Sleep(3 * time.Second)
	socksproxy()
}
