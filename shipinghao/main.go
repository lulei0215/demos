package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"shipinghao/model"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
)

// var uin = "0000000000"

type Data struct {
	Timestamp       string      `json:"timestamp"`
	LogFinderUin    string      `json:"_log_finder_uin"`
	LogFinderId     string      `json:"_log_finder_id"`
	RawKeyBuff      interface{} `json:"rawKeyBuff"`
	PluginSessionId interface{} `json:"pluginSessionId"`
	Scene           int         `json:"scene"`
	ReqScene        int         `json:"reqScene"`
}

// 1.获取二维码
func getQRCode1() (string, string) {
	fmt.Println(123)
	timestamp := strconv.Itoa(int(time.Now().Unix() * 1000))
	url := "https://channels.weixin.qq.com/cgi-bin/mmfinderassistant-bin/auth/auth_login_code"
	headers := map[string]string{
		"X-Wechat-Uin": uin,
	}
	data := Data{
		Timestamp:       timestamp,
		LogFinderUin:    "",
		LogFinderId:     "",
		RawKeyBuff:      nil,
		PluginSessionId: nil,
		Scene:           7,
		ReqScene:        7,
	}

	// 将data转换为JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}

	// 创建请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalln(err)
	}

	// 添加请求头
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	// 解析响应
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	// 获取token和二维码链接
	token := result["data"].(map[string]interface{})["token"].(string)
	img := "https://channels.weixin.qq.com/mobile/confirm_login.html?token=" + token

	err = qrcode.WriteFile(img, qrcode.Medium, 256, "images/vx_login_code.png")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Token: ", token)
	fmt.Println("QRCode Link: ", img)

	// Call create_session with the obtained token
	// return createSession(token)
	return img, token
}

// 1.获取二维码
func getQRCode() (string, string) {
	fmt.Println(123)
	timestamp := strconv.Itoa(int(time.Now().Unix() * 1000))
	url := "https://channels.weixin.qq.com/cgi-bin/mmfinderassistant-bin/auth/auth_login_code"
	headers := map[string]string{
		"X-Wechat-Uin": uin,
	}
	data := Data{
		Timestamp:       timestamp,
		LogFinderUin:    "",
		LogFinderId:     "",
		RawKeyBuff:      nil,
		PluginSessionId: nil,
		Scene:           7,
		ReqScene:        7,
	}

	// 将data转换为JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(123)
	// 创建请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalln(err)
	}

	// 添加请求头
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	// 解析响应
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	// 获取token和二维码链接
	token := result["data"].(map[string]interface{})["token"].(string)
	img := "https://channels.weixin.qq.com/mobile/confirm_login.html?token=" + token

	err = qrcode.WriteFile(img, qrcode.Medium, 256, "vx_login_code.png")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Token: ", token)
	fmt.Println("QRCode Link: ", img)

	// Call create_session with the obtained token
	// createSession(token)

	// img := "https://channels.weixin.qq.com/mobile/confirm_login.html?token=" + token
	// fmt.Println("--------------------------------")
	// fmt.Println(img)
	return img, token
}

var Cookie string

type ResponseData struct {
	Status     int `json:"status"`
	AcctStatus int `json:"acctStatus"`
}

//	type ResponseData struct {
//		ErrorCode int    `json:"errorCode"`
//		Msg       string `json:"msg"`
//		BaseResp  struct {
//			Errcode int `json:"errcode"`
//		} `json:"baseResp"`
//		ExportId string `json:"exportId"`
//	}

type AuthDataResponse struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	Data    struct {
		FinderUser struct {
			FinderUsername string `json:"finderUsername"`
		} `json:"finderUser"`
	} `json:"data"`
}

type Response struct {
	ErrCode int          `json:"errCode"`
	ErrMsg  string       `json:"errMsg"`
	Data    ResponseData `json:"data"`
}

// # 2.获取token
func createSession(token string) (string, string) {
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
	baseUrl := "https://channels.weixin.qq.com/cgi-bin/mmfinderassistant-bin/auth/auth_login_status"

	headers := map[string]string{
		"X-Wechat-Uin": "your_uin", // replace with your uin
	}
	data := map[string]interface{}{
		"token":           token,
		"timestamp":       timestamp,
		"_log_finder_uin": "",
		"_log_finder_id":  "",
		"rawKeyBuff":      nil,
		"pluginSessionId": nil,
		"scene":           7,
		"reqScene":        7,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return "", ""
	}

	params := url.Values{}
	params.Add("token", token)
	params.Add("timestamp", timestamp)
	params.Add("_log_finder_uin", "")
	params.Add("_log_finder_id", "")
	params.Add("scene", "7")
	params.Add("reqScene", "7")
	url2 := fmt.Sprintf("%s?%s", baseUrl, params.Encode())

	req, err := http.NewRequest("POST", url2, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
		return "", ""
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", ""
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var res Response
	json.Unmarshal(body, &res)

	// go func() {
	switch {
	case res.Data.Status == 0 && res.Data.AcctStatus == 0:
		fmt.Println("未登录")
		// time.Sleep(2 * time.Second)
		// createSession(token)
	case res.Data.Status == 5 && res.Data.AcctStatus == 1:
		fmt.Println("已经扫码未确认")
		// time.Sleep(2 * time.Second)
		// createSession(token)
	case res.Data.Status == 1 && res.Data.AcctStatus == 1:
		fmt.Println("登录成功")
		Cookie = resp.Header.Get("Set-Cookie")
		if Cookie != "" {
			getAuthData()
		} else {
			fmt.Println("Cookie获取失败")
		}
		fmt.Println(res.Data.Status)
	case res.Data.Status == 5 && res.Data.AcctStatus == 2:
		fmt.Println("没有可登录的视频号")
	case res.Data.Status == 4:
		fmt.Println("二维码已经过期")
	default:
		fmt.Println("网络错误")
	}
	// }()

	img := "https://channels.weixin.qq.com/mobile/confirm_login.html?token=" + token
	fmt.Println("--------------------------------")
	fmt.Println(img)
	return img, ""
}

// # 3.获取authData
func getAuthData() (string, error) {
	fmt.Println("get auth data")
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
	url := "https://channels.weixin.qq.com/cgi-bin/mmfinderassistant-bin/auth/auth_data"

	// Create JSON data
	data := map[string]interface{}{
		"timestamp":       timestamp,
		"_log_finder_uin": "",
		"_log_finder_id":  "",
		"rawKeyBuff":      nil,
		"pluginSessionId": nil,
		"scene":           7,
		"reqScene":        7,
	}
	fmt.Println("jsonData")
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return "数据异常", err
	}
	fmt.Println("Create request")
	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
		return "数据异常", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Wechat-Uin", "your_uin")
	req.Header.Set("Cookie", Cookie)

	fmt.Println("client")
	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "数据异常", err
	}
	defer resp.Body.Close()

	// Decode the response
	body, _ := ioutil.ReadAll(resp.Body)
	var res AuthDataResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return "数据异常", err
	}
	fmt.Println(res)
	fmt.Println("finderUsername")
	// v2_060000231003b20faec8c6eb8b19cad4c807ec3cb0774ca0dd7b6df1e4aa4943f964c1951ff7
	finderUsername := res.Data.FinderUser.FinderUsername
	if finderUsername != "" {
		get_x_wechat_uin(finderUsername)
		go getLoginCookie(finderUsername)
		go getVideoList(finderUsername)
		// upload_video()
		// locaction, err := get_location(finderUsername)
		// if err != nil {
		// 	fmt.Println("get_location获取失败")
		// }
		// fmt.Println(locaction)
		fmt.Println("authData获取ok")
		// traceKey := get_trace_key(finderUsername)
		// fmt.Println(traceKey)
		// time.Sleep(2 * time.Second)
		// publish_video(traceKey, locaction.Latitude, locaction.Longitude, locaction.City, finderUsername)
	} else {
		fmt.Println("authData获取失败")
	}
	return finderUsername, nil
}

type UinResponse struct {
	Data struct {
		Uin int `json:"uin"`
	} `json:"data"`
}

// 4.获取登录cookie 接收私信消息的 token
// type Data struct {
// 	Timestamp int64 `json:"timestamp"`
// 	LogFinderUin string `json:"_log_finder_uin"`
// 	LogFinderId string `json:"_log_finder_id"`
// 	RawKeyBuff interface{} `json:"rawKeyBuff"`
// 	PluginSessionId interface{} `json:"pluginSessionId"`
// 	Scene int `json:"scene"`
// 	ReqScene int `json:"reqScene"`
// }

//	type Cookie struct {
//		Data struct {
//			Cookie string `json:"cookie"`
//		} `json:"data"`
//	}
type CookieData struct {
	Timestamp       int64       `json:"timestamp"`
	LogFinderUin    string      `json:"_log_finder_uin"`
	LogFinderId     string      `json:"_log_finder_id"`
	RawKeyBuff      interface{} `json:"rawKeyBuff"`
	PluginSessionId interface{} `json:"pluginSessionId"`
	Scene           int         `json:"scene"`
	ReqScene        int         `json:"reqScene"`
}

type CookieCookie struct {
	Data struct {
		Cookie string `json:"cookie"`
	} `json:"data"`
}
type Student struct {
	Name string
	Age  int
}

func getLoginCookie(finderUsername string) {
	timestamp := int(time.Now().UnixNano() / 1000000)
	fmt.Println("get_login_cookie")

	url := "https://channels.weixin.qq.com/cgi-bin/mmfinderassistant-bin/private-msg/get-login-cookie"
	headers := map[string]string{
		"X-Wechat-Uin": uin,    // 假设uin已定义
		"Cookie":       Cookie, // 需要在实际代码中替换为有效的Cookie值
	}

	data := map[string]interface{}{
		"timestamp":       timestamp,
		"_log_finder_uin": "",
		"_log_finder_id":  finderUsername,
		"rawKeyBuff":      nil,
		"pluginSessionId": nil,
		"scene":           7,
		"reqScene":        7,
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Failed to marshal data:", err)
		return
	}

	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		request.Header.Set(k, v)
	}
	request.Header.Set("Cookie", Cookie) // 在实际代码中替换为有效的Cookie值

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Failed to send request:", err)
		return
	}
	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return
	}

	var res map[string]interface{}
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		fmt.Println("Failed to parse JSON response:", err)
		return
	}

	// s1 := Student{"小红", 12}
	// s2 := Student{"小兰", 10}
	// s3 := Student{"小黄", 11}

	// _ := model.MongoDB()

	// document := bson.M{"name": "John", "age": 30, "city": "New York"}

	// result, err := model.InsertDocument(client1, "hehe", "chat", document)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("mongodb: ********************************")
	// fmt.Println(result)
	// fmt.Println("mongodb: ********************************")

	// err = client1.Disconnect(context.Background())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Disconnected from MongoDB!")

	cookie, ok := res["data"].(map[string]interface{})["cookie"]
	if ok && cookie != nil {
		fmt.Println("--------------cookie------------------")
		fmt.Println(cookie)
		fmt.Println(res)
		fmt.Println("---------------cookie-----------------")
		getNewMsg(finderUsername, cookie.(string)) // 接收私信消息，假设getNewMsg函数已定义
	} else {
		fmt.Println("登录cookie获取失败")
	}
}

// 5.接收私信消息

func getNewMsg(finderUsername string, cookie string) {

	timestamp := int(time.Now().UnixNano() / 1000000)
	fmt.Println("接收私信消息 get_new_msg")

	url := "https://channels.weixin.qq.com/cgi-bin/mmfinderassistant-bin/private-msg/get-new-msg"
	headers := map[string]string{
		"X-Wechat-Uin": uin, // 假设uin已定义
		// "Cookie": Cookie, // 在实际代码中替换为有效的Cookie值
	}

	data := map[string]interface{}{
		"cookie":          cookie,
		"timestamp":       timestamp,
		"_log_finder_uin": "",
		"_log_finder_id":  finderUsername,
		"rawKeyBuff":      nil,
		"pluginSessionId": nil,
		"scene":           7,
		"reqScene":        7,
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Failed to marshal data:", err)
		return
	}

	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		request.Header.Set(k, v)
	}
	request.Header.Set("Cookie", Cookie) // 在实际代码中替换为有效的Cookie值

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Failed to send request:", err)
		return
	}
	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return
	}

	var res map[string]interface{}
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		fmt.Println("Failed to parse JSON response:", err)
		return
	}
	fmt.Println("*************************msglist***************************************")
	fmt.Println(res)

	fmt.Println("***************************msglist*************************************")
	msg, ok := res["data"].(map[string]interface{})["msg"].([]interface{})
	if !ok {
		fmt.Println("Failed to extract 'msg' from response")
		return
	}

	for _, i := range msg {
		message, ok := i.(map[string]interface{})
		if !ok {
			fmt.Println("Invalid message format")
			continue
		}
		if message["rawContent"] == "你好" {
			time.Sleep(2 * time.Second)
			// 调用send_private_msg函数发送回复（假设已定义）
			sendPrivateMsg(
				finderUsername,
				message["sessionId"].(string),
				message["toUsername"].(string),
				message["fromUsername"].(string),
			)
		} else {
			fmt.Println("其他消息")
		}
	}

	// 每隔5秒请求一次
	time.Sleep(5 * time.Second)
	fmt.Println("每5秒自动请求555-----------------")

	// 使用新cookie重新获取消息
	newCookie, ok := res["data"].(map[string]interface{})["cookie"].(string)
	if ok {
		getNewMsg(finderUsername, newCookie)
	} else {
		fmt.Println("无法从响应中获取新的cookie")
	}
}

// 6.回复私信消息
func sendPrivateMsg(finderUsername string, sessionId string, toUsername string, fromUsername string) {
	myUUID := uuid.New().String()
	timestamp := fmt.Sprintf("%d", int(time.Now().UnixNano()/1000000))

	url := "https://channels.weixin.qq.com/cgi-bin/mmfinderassistant-bin/private-msg/send-private-msg"
	headers := map[string]string{
		"X-Wechat-Uin": uin, // 假设uin已定义
		"Content-Type": "application/json",
	}

	data := map[string]interface{}{
		"msgPack": map[string]interface{}{
			"sessionId":    sessionId,
			"fromUsername": toUsername,
			"toUsername":   fromUsername,
			"msgType":      1,
			"textMsg":      map[string]interface{}{"content": "你好呀！"},
			"cliMsgId":     myUUID,
		},
		"timestamp":       timestamp,
		"_log_finder_uin": "",
		"_log_finder_id":  finderUsername,
		"rawKeyBuff":      nil,
		"pluginSessionId": nil,
		"scene":           7,
		"reqScene":        7,
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Failed to marshal data:", err)
		return
	}

	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(reqBody)))
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}

	for k, v := range headers {
		request.Header.Set(k, v)
	}
	request.Header.Set("Cookie", Cookie) // 在实际代码中替换为有效的Cookie值

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Failed to send request:", err)
		return
	}
	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return
	}

	var res map[string]interface{}
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		fmt.Println("Failed to parse JSON response:", err)
		return
	}

	fmt.Println("回复消息", res)
}

var uin string

// 8.获取X-Wechat-Uin
func get_x_wechat_uin(finderUsername string) {
	fmt.Println("get_x_wechat_uin")
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
	url := "https://channels.weixin.qq.com/cgi-bin/mmfinderassistant-bin/helper/helper_upload_params"

	// Create JSON data
	data := map[string]interface{}{
		"timestamp":       timestamp,
		"_log_finder_uin": "",
		"_log_finder_id":  finderUsername,
		"rawKeyBuff":      nil,
		"pluginSessionId": nil,
		"scene":           7,
		"reqScene":        7,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Wechat-Uin", "0000000000")
	req.Header.Set("Cookie", Cookie)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	// Decode the response
	body, _ := ioutil.ReadAll(resp.Body)
	var res UinResponse
	err = json.Unmarshal(body, &res)

	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	uin = strconv.Itoa(res.Data.Uin)
	fmt.Println("****************uni")
	fmt.Println(uin)
	fmt.Println(res)
	fmt.Println("****************uni")
}

// Helper function to convert struct to map
func structToMap(data Data) map[string]interface{} {
	dataJSON, _ := json.Marshal(data)
	dataMap := make(map[string]interface{})
	json.Unmarshal(dataJSON, &dataMap)
	return dataMap
}

type LocationResponse struct {
	Data struct {
		Address LocationAddress `json:"address"`
	} `json:"data"`
}
type LocationAddress struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	City      string  `json:"city"`
}

type TraceKeyResponse struct {
	Data struct {
		TraceKey string `json:"traceKey"`
	} `json:"data"`
}

func upload_video() {
	fmt.Println("upload_video")
	url := "https://finderassistancea.video.qq.com/applyuploaddfs"
	headers := map[string]string{
		"X-Arguments":   "apptype=251&filetype=20304&weixinnum=981816192&filekey=finder_video_img.jpeg&filesize=369410&taskid=0153ddfd-305f-4dca-aea6-748af8bafa47&scene=0",
		"Content-Type":  "application/json",
		"Authorization": "303e0201010437303502010102010102043a855380020101020102020404030201020320141d020412124475020414efe97902046622167602045e79fee50400",
	}
	data := map[string]interface{}{
		"BlockSum":        1,
		"BlockPartLength": []int{369410},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create request
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", headers["Content-Type"])
	req.Header.Set("X-Arguments", headers["X-Arguments"])
	req.Header.Set("Authorization", headers["Authorization"])

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	// Decode the response
	body, _ := ioutil.ReadAll(resp.Body)
	var res map[string]interface{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	fmt.Println("upload_video")
	fmt.Println(res)
	// Print DownloadURL
	fmt.Println(res["DownloadURL"])
}

func get_location(finderUsername string) (LocationAddress, error) {
	fmt.Println("get_location location")
	var laddress LocationAddress
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
	url := "https://channels.weixin.qq.com/cgi-bin/mmfinderassistant-bin/helper/helper_search_location"
	headers := map[string]string{
		"X-Wechat-Uin": uin,
	}
	fmt.Println(1)
	data := map[string]interface{}{
		"query":           "",
		"cookies":         "",
		"longitude":       0,
		"latitude":        0,
		"timestamp":       timestamp,
		"_log_finder_uin": "",
		"_log_finder_id":  finderUsername,
		"rawKeyBuff":      nil,
		"pluginSessionId": nil,
		"scene":           7,
		"reqScene":        7,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return laddress, err
	}
	fmt.Println(2)
	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
		return laddress, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Wechat-Uin", headers["X-Wechat-Uin"])
	fmt.Println(3)
	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return laddress, err
	}
	defer resp.Body.Close()
	fmt.Println(4)
	// Decode the response
	body, _ := ioutil.ReadAll(resp.Body)
	var res LocationResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return laddress, err
	}
	fmt.Println(5)
	fmt.Println(res.Data)
	fmt.Println(6)
	fmt.Println(res.Data.Address.City, res.Data.Address.Latitude)
	laddress.City = res.Data.Address.City
	laddress.Latitude = res.Data.Address.Latitude
	laddress.Longitude = res.Data.Address.Longitude
	fmt.Println(6)
	return laddress, nil

}

func get_trace_key(finderUsername string) string {
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
	url := "https://channels.weixin.qq.com/cgi-bin/mmfinderassistant-bin/post/get-finder-post-trace-key"
	headers := map[string]string{
		"X-Wechat-Uin": uin,
	}
	data := map[string]interface{}{
		"objectId":        nil,
		"timestamp":       timestamp,
		"_log_finder_uin": "",
		"_log_finder_id":  finderUsername,
		"rawKeyBuff":      nil,
		"pluginSessionId": nil,
		"scene":           7,
		"reqScene":        7,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Wechat-Uin", headers["X-Wechat-Uin"])

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()

	// Decode the response
	body, _ := ioutil.ReadAll(resp.Body)
	var res TraceKeyResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return ""
	}

	return res.Data.TraceKey
}

func publish_video(traceKey string, latitude string, longitude string, city string, finderUsername string) {
	myUUID := uuid.New().String()
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
	m_timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	url := "https://channels.weixin.qq.com/cgi-bin/mmfinderassistant-bin/post/post_create"
	headers := map[string]string{
		"X-Wechat-Uin": uin,
		"Content-Type": "application/json",
	}
	data := map[string]interface{}{
		"objectType":      0,
		"longitude":       0,
		"latitude":        0,
		"feedLongitude":   0,
		"feedLatitude":    0,
		"originalFlag":    0,
		"topics":          []interface{}{},
		"isFullPost":      1,
		"handleFlag":      2,
		"videoClipTaskId": "",
		"traceInfo": map[string]string{
			"traceKey":       traceKey,
			"uploadCdnStart": m_timestamp,
			"uploadCdnEnd":   m_timestamp,
		},
		"objectDesc": map[string]interface{}{
			"mpTitle":     "",
			"description": "#",
			"extReading":  map[string]interface{}{},
			"mediaType":   2,
			"location": map[string]interface{}{
				"latitude":      latitude,
				"longitude":     longitude,
				"city":          city,
				"poiClassifyId": "",
			},
			"topic": map[string]string{
				"finderTopicInfo": "<finder><version>1</version><valuecount>1</valuecount><style><at></at></style><value0><![CDATA[#]]></value0></finder>",
			},
			"event":         map[string]interface{}{},
			"mentionedUser": []interface{}{},
			"media": []map[string]interface{}{
				{
					"url":          "https://finder.video.qq.com/251/20304/stodownload?bizid=1023&dotrans=0&encfilekey=rjD5jyTuFrIpZ2ibE8T7YmwgiahniaXswqzlcPHbg3keBzW4VS8lwPicXwrhNTBNPANzbmJy5eCYtU9wW84HkFg0SodsuItTiaJCntECoD6KjS40Xic9YW9F23pQ&hy=SH&idx=1&m=&scene=0&token=6xykWLEnztLZ2p2FvRHDeQllaUUlhtY4uib9kA3OuU1hn8128ltYoZ03MExib0nWAsib78xicM1SWibEdGia1hactfMw",
					"fileSize":     3654801,
					"thumbUrl":     "https://finder.video.qq.com/251/20350/stodownload?bizid=1023&dotrans=0&encfilekey=rjD5jyTuFrIpZ2ibE8T7YmwgiahniaXswqzP07ibSibiciaOBic48bnlJPOzGP19NJZPg9xqnVu1cFkQmQGibHTOSta4RjeAORYGoJWPW5mp97cFRmrlGGiasWIibhiaww&hy=SH&idx=1&m=&scene=0&token=cztXnd9GyrF2IupOnSmYjcggicYCNZSU0HfvmaQHkqcjCAapPClozqrSbaYIXazT1UXprqOqYMJYqMqsviamW5yg",
					"fullThumbUrl": "https://finder.video.qq.com/251/20350/stodownload?bizid=1023&dotrans=0&encfilekey=rjD5jyTuFrIpZ2ibE8T7YmwgiahniaXswqzP07ibSibiciaOBic48bnlJPOzGP19NJZPg9xqnVu1cFkQmQGibHTOSta4RjeAORYGoJWPW5mp97cFRmrlGGiasWIibhiaww&hy=SH&idx=1&m=&scene=0&token=cztXnd9GyrF2IupOnSmYjcggicYCNZSU0HfvmaQHkqcjCAapPClozqrSbaYIXazT1UXprqOqYMJYqMqsviamW5yg",
					"mediaType":    2,
					"videoPlayLen": 0,
					"width":        2600,
					"height":       3900,
					"md5sum":       "f025b779-4814-448c-840f-afb0bb51175e",
					"urlCdnTaskId": "",
				},
			},
			"member": map[string]interface{}{},
		},
		"postFlag":        0,
		"mode":            1,
		"clientid":        myUUID,
		"timestamp":       timestamp,
		"_log_finder_uin": "",
		"_log_finder_id":  finderUsername,
		"rawKeyBuff":      nil,
		"pluginSessionId": nil,
		"scene":           7,
		"reqScene":        7,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", headers["Content-Type"])
	req.Header.Set("X-Wechat-Uin", headers["X-Wechat-Uin"])

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	// Decode the response
	body, _ := ioutil.ReadAll(resp.Body)
	var res map[string]interface{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	fmt.Println(res)
}

// 获取视频列表
func getVideoList(finderUsername string) {
	print("获取视频列表 get_video_list")
	timestamp := int(time.Now().UnixNano() / 1000000)
	url := "https://channels.weixin.qq.com/cgi-bin/mmfinderassistant-bin/post/post_list"
	// headers := map[string]string{
	// 	"X-Wechat-Uin": uin, // 假设uin已定义
	// }

	data := map[string]interface{}{
		"pageSize":            10,
		"currentPage":         1,
		"onlyUnread":          false,
		"userpageType":        3,
		"needAllCommentCount": true,
		"forMcn":              false,
		"timestamp":           timestamp,
		"_log_finder_uin":     "",
		"_log_finder_id":      finderUsername,
		"rawKeyBuff":          nil,
		"pluginSessionId":     nil,
		"scene":               7,
		"reqScene":            7,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("***********************666*****************************************")

	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Wechat-Uin", uin)
	req.Header.Set("Cookie", Cookie)
	fmt.Println("******************************444**********************************")

	// reqBody, err := json.Marshal(data)
	// if err != nil {
	// 	fmt.Println("Failed to marshal data:", err)
	// 	return
	// }
	// print("获取视频列表********************************")

	// client := &http.Client{}
	// request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	// if err != nil {
	// 	fmt.Println("Failed to create request:", err)
	// 	return
	// }
	// print("获取视频列表********************************")
	// for k, v := range headers {
	// 	request.Header.Set(k, v)
	// }
	// request.Header.Set("Cookie", Cookie) // 在实际代码中替换为有效的Cookie值

	// response, err := client.Do(request)
	// if err != nil {
	// 	fmt.Println("Failed to send request:", err)
	// 	return
	// }
	// defer response.Body.Close()

	// bodyBytes, err := ioutil.ReadAll(response.Body)
	// if err != nil {
	// 	fmt.Println("Failed to read response body:", err)
	// 	return
	// }

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("***********************************123*****************************")

	// Decode the response
	body, _ := ioutil.ReadAll(resp.Body)

	var res map[string]interface{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		fmt.Println("Failed to parse JSON response:", err)
		return
	}
	fmt.Println(res)
	fmt.Println(res["errCode"])
	fmt.Println(res["errCode"] == 0)
	fmt.Println(res["errCode"] == "0")
	// if res["errCode"] != "0" {
	// 	fmt.Println("Failed to parse JSON response:", res["errMsg"])
	// 	return
	// }

	fmt.Println("****************************************************************")
	fmt.Println(res["data"])
	print(("****************************************************************"))

	list, ok := res["data"].(map[string]interface{})["list"].([]interface{})
	if !ok {
		fmt.Println("视频列表获取失败")
		return
	}

	for _, i := range list {
		item, ok := i.(map[string]interface{})
		if !ok {
			fmt.Println("无效的视频列表项格式")
			continue
		}
		exportId := item["exportId"].(string)
		// 获取评论列表
		getCommentList(finderUsername, exportId, item)
	}

}

// 获取评论列表
func getCommentList(finderUsername string, exportId string, videoItem interface{}) {
	fmt.Println("获取评论列表")
	timestamp := fmt.Sprintf("%d", int(time.Now().UnixNano()/1000000))
	// fmt.Println(uin) // 假设uin已定义

	url := "https://channels.weixin.qq.com/cgi-bin/mmfinderassistant-bin/comment/comment_list"
	headers := map[string]string{
		"X-Wechat-Uin": uin,
		"Content-Type": "application/json",
	}

	data := map[string]interface{}{
		"lastBuff":         "",
		"exportId":         exportId,
		"commentSelection": false,
		"forMcn":           false,
		"timestamp":        timestamp,
		"_log_finder_uin":  "",
		"_log_finder_id":   finderUsername,
		"rawKeyBuff":       nil,
		"pluginSessionId":  nil,
		"scene":            7,
		"reqScene":         7,
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Failed to marshal data:", err)
		return
	}

	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}

	for k, v := range headers {
		request.Header.Set(k, v)
	}
	request.Header.Set("Cookie", Cookie) // 在实际代码中替换为有效的Cookie值

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Failed to send request:", err)
		return
	}
	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return
	}

	var res map[string]interface{}
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		fmt.Println("Failed to parse JSON response:", err)
		return
	}
	fmt.Println("评论列表---res--")

	fmt.Println(res)
	commentData, ok := res["data"].(map[string]interface{})
	if !ok || commentData["comment"] == nil {
		fmt.Println("评论获取失败")
		return
	}
	fmt.Println("评论列表-----commentData")

	fmt.Println(commentData)
	comments, ok := commentData["comment"].([]interface{})

	if !ok {
		fmt.Println("无效的评论列表格式")
		return
	}
	fmt.Println("评论列表")

	fmt.Println(comments)

	for _, i := range comments {

		fmt.Println(i)
		comment, ok := i.(map[string]interface{})
		if !ok {
			fmt.Println("无效的评论格式")
			continue
		}
		content := comment["commentContent"].(string)
		fmt.Println(content)
		if content == "赞" {
			fmt.Println("回复评论")
			sendComment(finderUsername, exportId, comment)
		} else {
			fmt.Println("其他评论")
		}
	}
}

// 回复评论
func sendComment(finderUsername string, exportId string, comment interface{}) {
	myUUID := uuid.New().String()
	timestamp := fmt.Sprintf("%d", int(time.Now().UnixNano()/1000000))

	url := "https://channels.weixin.qq.com/cgi-bin/mmfinderassistant-bin/comment/create_comment"
	headers := map[string]string{
		"X-Wechat-Uin": uin, // 假设uin已定义
		"Content-Type": "application/json",
	}

	commentMap, ok := comment.(map[string]interface{})
	if !ok {
		fmt.Println("无效的评论对象格式")
		return
	}

	data := map[string]interface{}{
		"replyCommentId":  commentMap["commentId"],
		"content":         "哈哈",
		"clientId":        myUUID,
		"rootCommentId":   commentMap["commentId"],
		"comment":         comment,
		"exportId":        exportId,
		"timestamp":       timestamp,
		"_log_finder_uin": "",
		"_log_finder_id":  finderUsername,
		"rawKeyBuff":      nil,
		"pluginSessionId": nil,
		"scene":           7,
		"reqScene":        7,
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Failed to marshal data:", err)
		return
	}

	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}

	for k, v := range headers {
		request.Header.Set(k, v)
	}
	request.Header.Set("Cookie", Cookie) // 在实际代码中替换为有效的Cookie值

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Failed to send request:", err)
		return
	}
	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return
	}

	var res map[string]interface{}
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		fmt.Println("Failed to parse JSON response:", err)
		return
	}

	fmt.Println(res)
}

// 入口
func main() {
	// model.MongoDB()
	// return
	r := gin.Default()

	r.GET("/getIp", func(c *gin.Context) {
		model.GetIp()
		c.JSON(200, gin.H{
			"message": "GET",
		})
	})
	r.GET("/getQRCode", func(c *gin.Context) {
		getQRCode()
		c.JSON(200, gin.H{
			"message": "GET",
		})
	})

	r.POST("/book", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "POST",
		})
	})

	r.PUT("/book", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "PUT",
		})
	})

	//获取 私信列表
	r.GET("/getLoginCookie", func(c *gin.Context) {
		find_name := c.Query("find_name") // 获取query参数username
		getLoginCookie1(find_name)

		c.JSON(200, gin.H{
			"message": "DELETE",
		})
	})

	//获取 私信列表
	r.GET("/getnewcookie", func(c *gin.Context) {
		find_name := c.Query("find_name") // 获取query参数username
		cookie := c.Query("cookie")       // 获取query参数username
		getNewMsg_new(find_name, cookie)

		c.JSON(200, gin.H{
			"message": "DELETE",
		})
	})

	// 回复 私信
	r.GET("/sendPrivateMsg", func(c *gin.Context) {
		find_name := c.Query("find_name") // 获取query参数username
		sessionId := c.Query("sessionId")
		toUsername := c.Query("toUsername")
		fromUsername := c.Query("fromUsername")
		content := c.Query("content")

		code, msg := sendPrivateMsg1(
			find_name,
			sessionId,
			toUsername,
			fromUsername,
			content,
		)
		if msg != nil {
			c.JSON(200, gin.H{
				"msg":  msg.Error(),
				"code": code,
			})
		} else {
			c.JSON(200, gin.H{
				"msg":  "成功",
				"code": code,
			})
		}

	})
	// 回复私信

	// 页面扫码二维码登陆
	r.GET("/index", func(c *gin.Context) {
		url, token := getQRCode1()
		// 使用模板，第一个参数是状态码，第二个参数是模板名，第三个参数是传递给模板的数据
		c.HTML(200, "index.html", gin.H{
			"title": "Main website",
			"url":   url,
			"token": token,
		})
	})
	// 获取 _log_finder_id 和uin
	r.GET("/createSession", func(c *gin.Context) {
		token := c.Query("token") // 获取query参数username

		status, msg, uin := createSession1(token)
		// 使用模板，第一个参数是状态码，第二个参数是模板名，第三个参数是传递给模板的数据
		c.JSON(200, gin.H{
			"msg":            msg,
			"_log_finder_id": msg,
			"code":           status,
			"uin":            uin,
		})
	})

	r.Static("/images", "./images")
	r.LoadHTMLGlob("templates/*")

	r.Run()
	// ()
}

func yourGinHandler(c *gin.Context) {
	token := c.Query("token")
	go createSession1(token)

	c.JSON(200, gin.H{
		"message": "Session creation started",
	})
}

func createSession1(token string) (int64, string, string) {
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
	baseUrl := "https://channels.weixin.qq.com/cgi-bin/mmfinderassistant-bin/auth/auth_login_status"

	headers := map[string]string{
		"X-Wechat-Uin": "your_uin", // replace with your uin
	}
	data := map[string]interface{}{
		"token":           token,
		"timestamp":       timestamp,
		"_log_finder_uin": "",
		"_log_finder_id":  "",
		"rawKeyBuff":      nil,
		"pluginSessionId": nil,
		"scene":           7,
		"reqScene":        7,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return 0, err.Error(), ""
	}

	params := url.Values{}
	params.Add("token", token)
	params.Add("timestamp", timestamp)
	params.Add("_log_finder_uin", "")
	params.Add("_log_finder_id", "")
	params.Add("scene", "7")
	params.Add("reqScene", "7")
	url2 := fmt.Sprintf("%s?%s", baseUrl, params.Encode())

	req, err := http.NewRequest("POST", url2, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
		return 0, err.Error(), ""
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return 0, err.Error(), ""
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var res Response
	json.Unmarshal(body, &res)

	// go func() {
	switch {
	case res.Data.Status == 0 && res.Data.AcctStatus == 0:
		fmt.Println("未登录")
		return 0, "未登录", ""
		// time.Sleep(2 * time.Second)
		// createSession(token)
	case res.Data.Status == 5 && res.Data.AcctStatus == 1:
		// fmt.Println("已经扫码未确认")
		return 0, "已经扫码未确认", ""
		// time.Sleep(2 * time.Second)
		// createSession(token)
	case res.Data.Status == 1 && res.Data.AcctStatus == 1:
		fmt.Println("登录成功")
		Cookie = resp.Header.Get("Set-Cookie")
		if Cookie != "" {
			return getAuthData1()
		} else {
			return 0, "Cookie获取失败", ""
		}
		fmt.Println(res.Data.Status)
	case res.Data.Status == 5 && res.Data.AcctStatus == 2:
		return 0, "没有可登录的视频号", ""
	case res.Data.Status == 4:
		return 0, "二维码已经过期", ""
	default:
		return 0, "网络错误", ""
	}
	return 1, "登陆成功", ""
}

func getAuthData1() (int64, string, string) {
	fmt.Println("get auth data")
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
	url := "https://channels.weixin.qq.com/cgi-bin/mmfinderassistant-bin/auth/auth_data"

	// Create JSON data
	data := map[string]interface{}{
		"timestamp":       timestamp,
		"_log_finder_uin": "",
		"_log_finder_id":  "",
		"rawKeyBuff":      nil,
		"pluginSessionId": nil,
		"scene":           7,
		"reqScene":        7,
	}
	fmt.Println("jsonData")
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return 0, "数据异常", ""
	}
	fmt.Println("Create request")
	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
		return 0, "数据异常", ""
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Wechat-Uin", "your_uin")
	req.Header.Set("Cookie", Cookie)

	fmt.Println("client")
	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return 0, "数据异常", ""
	}
	defer resp.Body.Close()

	// Decode the response
	body, _ := ioutil.ReadAll(resp.Body)
	var res AuthDataResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)

		return 0, "数据异常", ""
	}
	fmt.Println(res)
	fmt.Println("finderUsername")
	// v2_060000231003b20faec8c6eb8b19cad4c807ec3cb0774ca0dd7b6df1e4aa4943f964c1951ff7
	finderUsername := res.Data.FinderUser.FinderUsername
	if finderUsername != "" {

		uin, err = get_x_wechat_uin1(finderUsername)
		if err != nil {
			return 0, "authData登陆失败", ""
		}
		return 1, finderUsername, uin

		// get_x_wechat_uin(finderUsername)
		// go getLoginCookie(finderUsername)
		// go getVideoList(finderUsername)
		// // upload_video()
		// // locaction, err := get_location(finderUsername)
		// // if err != nil {
		// // 	fmt.Println("get_location获取失败")
		// // }
		// // fmt.Println(locaction)
		// fmt.Println("authData获取ok")
		// // traceKey := get_trace_key(finderUsername)
		// // fmt.Println(traceKey)
		// // time.Sleep(2 * time.Second)
		// // publish_video(traceKey, locaction.Latitude, locaction.Longitude, locaction.City, finderUsername)
	} else {
		return 0, "authData登陆失败", ""
	}

}

func get_x_wechat_uin1(finderUsername string) (string, error) {
	fmt.Println("get_x_wechat_uin")
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
	url := "https://channels.weixin.qq.com/cgi-bin/mmfinderassistant-bin/helper/helper_upload_params"

	// Create JSON data
	data := map[string]interface{}{
		"timestamp":       timestamp,
		"_log_finder_uin": "",
		"_log_finder_id":  finderUsername,
		"rawKeyBuff":      nil,
		"pluginSessionId": nil,
		"scene":           7,
		"reqScene":        7,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Wechat-Uin", "0000000000")
	req.Header.Set("Cookie", Cookie)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer resp.Body.Close()

	// Decode the response
	body, _ := ioutil.ReadAll(resp.Body)
	var res UinResponse
	err = json.Unmarshal(body, &res)

	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return "", err
	}

	uin = strconv.Itoa(res.Data.Uin)
	// fmt.Println("****************uni")
	// fmt.Println(uin)
	// fmt.Println(res)
	// fmt.Println("****************uni")
	return uin, nil
}

type GetLoginCookieResponse struct {
	Data    GetLoginCookieData `json:"data"`
	ErrCode int                `json:"errCode"`
	ErrMsg  string             `json:"errMsg"`
}

type GetLoginCookieData struct {
	BaseResp GetLoginCookieBaseResp `json:"baseResp"`
	Cookie   string                 `json:"cookie"`
}

type GetLoginCookieBaseResp struct {
	Errcode int `json:"errcode"`
}

// 获取接收私信消息请求参数cookie
func getLoginCookie1(finderUsername string) {
	timestamp := int(time.Now().UnixNano() / 1000000)
	fmt.Println("get_login_cookie")

	url := "https://channels.weixin.qq.com/cgi-bin/mmfinderassistant-bin/private-msg/get-login-cookie"
	headers := map[string]string{
		"X-Wechat-Uin": uin,    // 假设uin已定义
		"Cookie":       Cookie, // 需要在实际代码中替换为有效的Cookie值
	}

	data := map[string]interface{}{
		"timestamp":       timestamp,
		"_log_finder_uin": "",
		"_log_finder_id":  finderUsername,
		"rawKeyBuff":      nil,
		"pluginSessionId": nil,
		"scene":           7,
		"reqScene":        7,
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Failed to marshal data:", err)
		return
	}

	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		request.Header.Set(k, v)
	}
	request.Header.Set("Cookie", Cookie) // 在实际代码中替换为有效的Cookie值

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Failed to send request:", err)
		return
	}
	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return
	}

	var res GetLoginCookieResponse
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		fmt.Println("Failed to parse JSON response:", err)
		return
	}
	fmt.Println("****************************************************************")
	fmt.Println(res)
	fmt.Println("状态码：：：")
	fmt.Println(res.ErrCode)
	fmt.Println(res.ErrMsg == "request failed")

	fmt.Println("****************************************************************")
	ok := res.Data.Cookie
	if ok != "" {
		getNewMsg_new(finderUsername, ok) // 接收私信消息，假设getNewMsg函数已定义
	} else {
		fmt.Println("登录cookie获取失败")
	}
}

func getNewMsg1(finderUsername string, cookie string) {

	timestamp := int(time.Now().UnixNano() / 1000000)
	fmt.Println("接收私信消息 get_new_msg1")

	url := "https://channels.weixin.qq.com/cgi-bin/mmfinderassistant-bin/private-msg/get-new-msg"
	headers := map[string]string{
		"X-Wechat-Uin": uin, // 假设uin已定义
		// "Cookie": Cookie, // 在实际代码中替换为有效的Cookie值
	}

	data := map[string]interface{}{
		"cookie":          cookie,
		"timestamp":       timestamp,
		"_log_finder_uin": "",
		"_log_finder_id":  finderUsername,
		"rawKeyBuff":      nil,
		"pluginSessionId": nil,
		"scene":           7,
		"reqScene":        7,
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Failed to marshal data:", err)
		return
	}

	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		request.Header.Set(k, v)
	}
	request.Header.Set("Cookie", Cookie) // 在实际代码中替换为有效的Cookie值

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Failed to send request:", err)
		return
	}
	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return
	}

	var res map[string]interface{}
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		fmt.Println("Failed to parse JSON response:", err)
		return
	}
	fmt.Println("*************************getNewMsg1***************************************")
	fmt.Println(res)

	fmt.Println("***************************getNewMsg1*************************************")
	msg, ok := res["data"].(map[string]interface{})["msg"].([]interface{})
	if !ok {
		fmt.Println("Failed to extract 'msg' from response")
		return
	}

	fmt.Println(msg)
	fmt.Println("***************************getNewMsg1----msg *************************************")

	// for _, i := range msg {
	// 	message, ok := i.(map[string]interface{})
	// 	if !ok {
	// 		fmt.Println("Invalid message format")
	// 		continue
	// 	}
	// 	if message["rawContent"] == "你好" {
	// 		time.Sleep(2 * time.Second)
	// 		// 调用send_private_msg函数发送回复（假设已定义）
	// 		sendPrivateMsg(
	// 			finderUsername,
	// 			message["sessionId"].(string),
	// 			message["toUsername"].(string),
	// 			message["fromUsername"].(string),
	// 		)
	// 	} else {
	// 		fmt.Println("其他消息")
	// 	}
	// }

	// // 每隔5秒请求一次
	// time.Sleep(5 * time.Second)
	// fmt.Println("每5秒自动请求555-----------------")

	// // 使用新cookie重新获取消息
	// newCookie, ok := res["data"].(map[string]interface{})["cookie"].(string)
	// if ok {
	// 	fmt.Println("重新获取")
	// 	getNewMsg(finderUsername, newCookie)
	// } else {
	// 	fmt.Println("无法从响应中获取新的cookie")
	// }
}

func getNewMsg2(finderUsername string, cookie string) {

	timestamp := int(time.Now().UnixNano() / 1000000)
	fmt.Println("接收私信消息 get_new_msg")

	url := "https://channels.weixin.qq.com/cgi-bin/mmfinderassistant-bin/private-msg/get-new-msg"
	headers := map[string]string{
		"X-Wechat-Uin": uin, // 假设uin已定义
		// "Cookie": Cookie, // 在实际代码中替换为有效的Cookie值
	}

	data := map[string]interface{}{
		"cookie":          cookie,
		"timestamp":       timestamp,
		"_log_finder_uin": "",
		"_log_finder_id":  finderUsername,
		"rawKeyBuff":      nil,
		"pluginSessionId": nil,
		"scene":           7,
		"reqScene":        7,
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Failed to marshal data:", err)
		return
	}

	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		request.Header.Set(k, v)
	}
	request.Header.Set("Cookie", Cookie) // 在实际代码中替换为有效的Cookie值

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Failed to send request:", err)
		return
	}
	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return
	}

	var res map[string]interface{}
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		fmt.Println("Failed to parse JSON response:", err)
		return
	}
	fmt.Println("*************************msglist***************************************")
	fmt.Println(res)

	fmt.Println("***************************msglist*************************************")
	// msg, ok := res["data"].(map[string]interface{})["msg"].([]interface{})
	// if !ok {
	// 	fmt.Println("Failed to extract 'msg' from response")
	// 	return
	// }

	// for _, i := range msg {
	// 	message, ok := i.(map[string]interface{})
	// 	if !ok {
	// 		fmt.Println("Invalid message format")
	// 		continue
	// 	}
	// 	if message["rawContent"] == "你好" {
	// 		time.Sleep(2 * time.Second)
	// 		// 调用send_private_msg函数发送回复（假设已定义）
	// 		sendPrivateMsg(
	// 			finderUsername,
	// 			message["sessionId"].(string),
	// 			message["toUsername"].(string),
	// 			message["fromUsername"].(string),
	// 		)
	// 	} else {
	// 		fmt.Println("其他消息")
	// 	}
	// }

	// // 每隔5秒请求一次
	// time.Sleep(5 * time.Second)
	// fmt.Println("每5秒自动请求555-----------------")

	// // 使用新cookie重新获取消息
	newCookie, ok := res["data"].(map[string]interface{})["cookie"].(string)
	if ok {
		fmt.Println("重新获取 新cookie")
		fmt.Println(newCookie)
		// getNewMsg(finderUsername, newCookie)
	} else {
		fmt.Println("无法从响应中获取新的cookie")
	}
}

// 专门 cokkie 获取 新私信
func getNewMsg_new(finderUsername string, cookie string) {

	timestamp := int(time.Now().UnixNano() / 1000000)
	fmt.Println("接收私信消息 get_new_msg")

	url := "https://channels.weixin.qq.com/cgi-bin/mmfinderassistant-bin/private-msg/get-new-msg"
	headers := map[string]string{
		"X-Wechat-Uin": uin, // 假设uin已定义
		// "Cookie": Cookie, // 在实际代码中替换为有效的Cookie值
	}

	data := map[string]interface{}{
		"cookie":          cookie,
		"timestamp":       timestamp,
		"_log_finder_uin": "",
		"_log_finder_id":  finderUsername,
		"rawKeyBuff":      nil,
		"pluginSessionId": nil,
		"scene":           7,
		"reqScene":        7,
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Failed to marshal data:", err)
		return
	}

	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		request.Header.Set(k, v)
	}
	request.Header.Set("Cookie", Cookie) // 在实际代码中替换为有效的Cookie值

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Failed to send request:", err)
		return
	}
	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return
	}

	var res MsgResponse

	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		fmt.Println("Failed to decode response body:", err)
		return
	}

	// var res map[string]interface{}
	// err = json.Unmarshal(bodyBytes, &res)
	// if err != nil {
	// 	fmt.Println("Failed to parse JSON response:", err)
	// 	return
	// }
	fmt.Println("获取的消息")
	fmt.Println(res)

	if res.ErrCode != 0 {
		fmt.Println("错误")
		fmt.Println(res.ErrCode)
	} else {

	}
	fmt.Println(res.Data.Msg)
	fmt.Println("获取的消息")
	if len(res.Data.Msg) > 0 {
		msg := res.Data.Msg
		for _, i := range msg {
			fmt.Println("--------------------------------")
			fmt.Println("sessionId", i.SessionId)
			fmt.Println("toUsername", i.ToUsername)
			fmt.Println("fromUsername", i.FromUsername)
			fmt.Println("conten", i.RawContent)
			fmt.Println("--------------------------------")

			sendPrivateMsg1(finderUsername, i.SessionId, i.ToUsername, i.FromUsername, i.RawContent)
		}

	}
	// msg, ok := res["data"].(map[string]interface{})["msg"].([]interface{})
	// if !ok {
	// 	fmt.Println("Failed to extract 'msg' from response")
	// 	return
	// }

	// for _, i := range msg {
	// 	message, ok := i.(map[string]interface{})
	// 	if !ok {
	// 		fmt.Println("Invalid message format")
	// 		continue
	// 	}
	// 	if message["rawContent"] == "你好" {
	// 		time.Sleep(2 * time.Second)
	// 		// 调用send_private_msg函数发送回复（假设已定义）
	// 		sendPrivateMsg(
	// 			finderUsername,
	// 			message["sessionId"].(string),
	// 			message["toUsername"].(string),
	// 			message["fromUsername"].(string),
	// 		)
	// 	} else {
	// 		fmt.Println("其他消息")
	// 	}
	// }

	// // 每隔5秒请求一次
	// time.Sleep(5 * time.Second)
	// fmt.Println("每5秒自动请求555-----------------")

	// // 使用新cookie重新获取消息
	ok := res.Data.Cookie
	if ok != "" {
		fmt.Println("重新获取 新cookie")
		fmt.Println(ok)
		// getNewMsg(finderUsername, newCookie)
	} else {
		fmt.Println("无法从响应中获取新的cookie")
	}
}

type MsgResponse struct {
	Data    MsgData `bson:"data"`
	ErrCode int     `bson:"errCode"`
	ErrMsg  string  `bson:"errMsg"`
}

type MsgData struct {
	BaseResp MsgBaseResp `bson:"baseResp"`
	Cookie   string      `bson:"cookie"`
	Msg      []Msg       `bson:"msg"`
}

type MsgBaseResp struct {
	Errcode int    `bson:"errcode"`
	ErrMsg  string `bson:"errMsg"`
}

type Msg struct {
	FromUsername string  `bson:"fromUsername"`
	MsgType      int     `bson:"msgType"`
	RawContent   string  `bson:"rawContent"`
	Seq          int     `bson:"seq"`
	SessionId    string  `bson:"sessionId"`
	SvrMsgId     string  `bson:"svrMsgId"` // 修改为string类型
	TextMsg      TextMsg `bson:"textMsg"`
	ToUsername   string  `bson:"toUsername"`
	Ts           float64 `bson:"ts"`
}

type TextMsg struct {
	Content string `bson:"content"`
}

// 6.回复私信消息
func sendPrivateMsg1(finderUsername string, sessionId string, toUsername string, fromUsername string, content string) (code int64, err error) {
	myUUID := uuid.New().String()
	timestamp := fmt.Sprintf("%d", int(time.Now().UnixNano()/1000000))

	url := "https://channels.weixin.qq.com/cgi-bin/mmfinderassistant-bin/private-msg/send-private-msg"
	headers := map[string]string{
		"X-Wechat-Uin": uin, // 假设uin已定义
		"Content-Type": "application/json",
	}

	data := map[string]interface{}{
		"msgPack": map[string]interface{}{
			"sessionId":    sessionId,
			"fromUsername": toUsername,
			"toUsername":   fromUsername,
			"msgType":      1,
			"textMsg":      map[string]interface{}{"content": content},
			"cliMsgId":     myUUID,
		},
		"timestamp":       timestamp,
		"_log_finder_uin": "",
		"_log_finder_id":  finderUsername,
		"rawKeyBuff":      nil,
		"pluginSessionId": nil,
		"scene":           7,
		"reqScene":        7,
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Failed to marshal data:", err)
		return 0, err
	}

	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(reqBody)))
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return 0, err
	}

	for k, v := range headers {
		request.Header.Set(k, v)
	}
	request.Header.Set("Cookie", Cookie) // 在实际代码中替换为有效的Cookie值

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Failed to send request:", err)
		return 0, err
	}
	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return 0, err
	}

	var res GetLoginCookieResponse
	err = json.Unmarshal(bodyBytes, &res)

	if err != nil {
		fmt.Println("Failed to parse JSON response:", err)
		return 0, err
	}
	if res.ErrCode != 0 {
		fmt.Println(res.ErrMsg)
		return 0, err
	}
	return 1, nil
}
