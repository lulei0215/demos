package main

import (
	"encoding/json"
	"fmt"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"net/http"
	"time"
)

type ora struct {
	Id    string `json:"id"`
	Qu    string `json:"qu"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Type  string `json:"type"`
	Level string `json:"level"`
	Local string `json:"local"`
}
type res struct {
	Status   string `json:"status"`
	Info     string `json:"info"`
	Infocode string `json:"infocode"`
	Pois     []Pois `json:"pois"`
	Count    string `json:"count"`
}
type BaiduRes struct {
	Status     int            `json:"status"`
	Message    string         `json:"message"`
	ResultType string         `json:"result_type"`
	Results    []BaiduResults `json:"results"`
}
type BaiduResults struct {
	Name     string `json:"name"`
	Location struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	} `json:"location"`
	Address   string `json:"address"`
	Province  string `json:"province"`
	City      string `json:"city"`
	Area      string `json:"area"`
	StreetId  string `json:"street_id"`
	Telephone string `json:"telephone"`
	Detail    int    `json:"detail"`
	Uid       string `json:"uid"`
}

type Pois struct {
	Name     string `json:"name"`
	Id       string `json:"id"`
	Location string `json:"location"`
	Type     string `json:"type"`
	Pname    string `json:"pname"`
	Cityname string `json:"cityname"`
	Adname   string `json:"adname"`
	Address  string `json:"address"`
	Pcode    string `json:"pcode"`
	Citycode string `json:"citycode"`
	Adcode   string `json:"adcode"`
}

func readXlsx(filename string) []ora {
	var listOra []ora
	xlFile, err := xlsx.OpenFile(filename)
	if err != nil {
		fmt.Printf("open failed: %s\n", err)
	}
	for _, sheet := range xlFile.Sheets {

		tmpOra := ora{}
		// 获取标签页(时间)
		//tmpOra.TIME = sheet.Name
		for _, row := range sheet.Rows {

			var strs []string

			for _, cell := range row.Cells {
				text := cell.String()
				strs = append(strs, text)
			}
			//tmpOra.Id = num
			tmpOra.Qu = strs[1]
			tmpOra.Name = strs[2]
			tmpOra.Phone = strs[3]
			if strs[4] == "" {
				tmpOra.Type = ""
			} else {
				tmpOra.Type = strs[4]
			}
			if strs[5] == "" {
				tmpOra.Level = ""
			} else {
				tmpOra.Level = strs[5]
			}
			if strs[2] != "" {
				//url := "https://restapi.amap.com/v5/place/text?key=6530f54b07ad1786c9e5ca8508bc3a3a&keywords=" + strs[2]
				url := "https://api.map.baidu.com/place/v2/search?region=合肥&output=json&ak=yB1B3rlZaw0O8h29u5fEXTVHQpP4aPjb&query=" + strs[2]

				resp, err := GetHttp(url)
				res := BaiduRes{}
				if err != nil {
					fmt.Println(err)
				}
				json.Unmarshal(resp, &res)
				//fmt.Println(key, res)
				if len(res.Results) == 0 {
					tmpOra.Local = ""
				} else {
					tmpOra.Local = fmt.Sprintf("%f", res.Results[0].Location.Lng) + "," + fmt.Sprintf("%f", res.Results[0].Location.Lat)
				}
			} else {
				tmpOra.Local = ""
			}

			//fmt.Println(tmpOra)

			listOra = append(listOra, tmpOra)
		}
	}
	return listOra
}

func main() {
	var name string
	fmt.Printf("Please enter your file name: ")
	fmt.Scanf("%s", &name)
	//excelFileName := "C:\\Users\\llw98\\Desktop\\灾备数据库复制总量\\2019-04-26-2019-05-02Lag延时数据.xlsx"
	excelFileName := "/Users/ll/Documents/合肥幼儿园20219更新.xlsx"
	//excelFileName := name
	oraList := readXlsx(excelFileName)
	//fmt.Println(oraList)

	marshal, err := json.Marshal(oraList)
	if err != nil {
		fmt.Println("序列化失败，err=", err)
		return
	}
	fmt.Println(string(marshal))

	return
	//writingXlsx(oraList)

}

func writingXlsx(oraList []ora) {
	fmt.Println("开始写入")
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}

	row = sheet.AddRow()
	row.SetHeightCM(0.5)
	cell = row.AddCell()
	cell.Value = "序号"
	cell = row.AddCell()
	cell.Value = "县（市）区"
	cell = row.AddCell()
	cell.Value = "幼儿园名称"
	cell = row.AddCell()
	cell.Value = "联系电话"
	cell = row.AddCell()
	cell.Value = "办园性质"
	cell = row.AddCell()
	cell.Value = "办园等级"
	cell = row.AddCell()
	cell.Value = "经纬度"
	cell = row.AddCell()

	for j, i := range oraList {
		if j < 1 {
			continue
		}

		// 判断是否为-9999，是的变为0.0
		var row1 *xlsx.Row
		row1 = sheet.AddRow()
		row1.SetHeightCM(0.5)
		cell = row1.AddCell()
		cell.Value = string(j)
		cell = row1.AddCell()
		cell.Value = i.Qu
		cell = row1.AddCell()
		cell.Value = i.Name
		cell = row1.AddCell()
		cell.Value = i.Phone
		cell = row1.AddCell()
		cell.Value = i.Type
		cell = row1.AddCell()
		cell.Value = i.Level
		cell = row1.AddCell()
		cell.Value = i.Local

		//if i.v1000 == "-9999" {
		//	i.v1000 = "0.0"
		//}
		//if i.v2000 == "-9999" {
		//	i.v2000 = "0.0"
		//}
		//if i.H == "-9999" {
		//	i.H = "0.0"
		//}
		//if i.L == "-9999" {
		//	i.L = "0.0"
		//}

		// 打印时间
		//cell = row1.AddCell()
	}

	err = file.Save("2019-_-_-2019-_-_Lag延时数据.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
}
func GetHttp(url string) (body []byte, err error) {

	// 创建 client 和 resp 对象
	var client http.Client
	var resp *http.Response

	// 这里博主设置了10秒钟的超时
	client = http.Client{Timeout: 10 * time.Second}

	// 这里使用了 Get 方法，并判断异常
	resp, err = client.Get(url)
	if err != nil {
		return nil, err
	}
	// 释放对象
	defer resp.Body.Close()

	// 把获取到的页面作为返回值返回
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// 释放对象
	defer client.CloseIdleConnections()

	return body, nil
}
