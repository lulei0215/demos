package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type Product struct {
	Name        string
	Code        string
	Link        string
	SalePrice   string
	MarketPrice string
	Size        string
}

func main() {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// 用于存储页面HTML的变量
	// 	var htmlContent string

	// 	// 使用chromedp.Run运行一系列任务
	// 	err := chromedp.Run(ctx,
	// 		// 使用chromedp.Navigate导航到一个URL
	// 		chromedp.Navigate(`https://www.yougou.com`),
	// 		// 获取页面HTML
	// 		chromedp.OuterHTML("html", &htmlContent),
	// 	)
	// 	if err != nil {
	// 		fmt.Println("获取页面HTML时发生错误:", err)
	// 		return
	// 	}

	// 	// 使用goquery解析HTML
	// 	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	// 	if err != nil {
	// 		fmt.Println("解析HTML时发生错误:", err)
	// 		return
	// 	}
	// // 查找所有的a标签
	// doc.Find("div.nav-container").Find("div.nav").Find("a").Each(func(i int, s *goquery.Selection) {
	// 	// 获取a标签的href属性
	// 	href, _ := s.Attr("href")
	// 	// 打印a标签的href属性
	// 	fmt.Println(href)
	// })

	// `https://www.yougou.com/f-belle-MXZ-0-1.html`
	// 从键盘读取URL
	fmt.Print("请输入网址URL: ")
	var url string
	_, err := fmt.Scanln(&url)
	if err != nil {
		fmt.Println("读取输入时发生错误:", err)
		return
	}

	err = chromedp.Run(ctx,
		network.Enable(),
		chromedp.Navigate(url),        // 替换为你的网站URL
		chromedp.Sleep(1*time.Second), // 等待页面加载
	)

	if err != nil {
		log.Fatal(err)
	}
	var product []Product
	for {
		var html string
		err := chromedp.Run(ctx,
			chromedp.EvaluateAsDevTools(`window.scrollTo(0, document.body.scrollHeight)`, nil), // 滚动到页面底部
			chromedp.Sleep(1*time.Second),     // 等待新内容加载
			chromedp.OuterHTML("html", &html), // 获取HTML
		)
		if err != nil {
			log.Fatal(err)
		}

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {
			log.Fatal(err)
		}
		found := false

		doc.Find("ul#proList").Find("li").Find("div.srchlst-wrap").Each(func(i int, s *goquery.Selection) {
			fmt.Println("开始：************")
			aSelection := s.Find("a").First()

			href, _ := aSelection.Attr("href")
			// fmt.Printf("Found %d a tags\n", aSelection.Length())
			// fmt.Println(href)
			if !strings.HasPrefix(href, "https://www.yougou.com/") {
				href = "https://www.yougou.com/" + href
			}

			// s.Find("a").Each(func(i int, s *goquery.Selection) {
			// 	var href string
			// 	href, _ = s.Attr("href")
			// title, _ := s.Attr("title")

			// if exists {
			// fmt.Println("商品名："+title+"链接：https://www.yougou.com"+href)

			var detail_href = href

			// 打开新的页面
			// 	ctx1, cancel1 := chromedp.NewContext(context.Background())
			// defer cancel1()
			err := chromedp.Run(ctx, chromedp.Navigate(detail_href))
			if err != nil {
				log.Fatal(err)
			}
			// 查找特定的元素
			var detail, price, price2, title, sku string

			// start := strings.Index(href, "sku-") + len("sku-") // 找到"sku-"后面的位置
			// end := strings.Index(href[start:], "-") + start    // 找到下一个"-"的位置
			// if start < end {
			// 	sku = href[start:end]
			// }

			err = chromedp.Run(ctx, chromedp.OuterHTML("html", &detail))
			if err != nil {
				fmt.Println("chromedp.Run(ctx, error")
				log.Fatal(err)
			}
			// yitianPrice
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(detail))
			if err != nil {
				log.Fatal(err)
			}
			doc.Find("h1").Each(func(i int, s *goquery.Selection) {
				title = s.Text()

			})
			doc.Find("#yitianPrice").Find("i").Each(func(i int, s *goquery.Selection) {
				price = s.Text()

			})
			doc.Find("del").Each(func(i int, s *goquery.Selection) {
				price2 = s.Text()

			})
			// doc.Find("#shopping-container").Find("h1").Each(func(i int, s *goquery.Selection) {
			// 	title = s.Text()

			// })

			var mashu string
			doc.Find(".fl.prodSpec.size").First().Find("a").Each(func(i int, s *goquery.Selection) {
				dataName, exists := s.Attr("data-name")
				if exists {
					// fmt.Println(data
					// Name)

					mashu += dataName + ","

				}

				// fmt.Println("chromedp.Run(ctx1, dataName")
			})
			re := regexp.MustCompile(`[A-Za-z0-9]+$`) // 匹配结尾的字母和数字
			sku = re.FindString(title)

			fmt.Println("商品：" + title + "货号：" + sku + "链接：" + detail_href + "售价：" + price + "市场价：" + price2 + "尺码：" + mashu)
			p := Product{
				Name:        title,
				Code:        sku,
				Link:        detail_href,
				SalePrice:   price,
				MarketPrice: price2,
				Size:        mashu,
			}

			product = append(product, p)
			found = true

			// }

			// })

		})

		// found = false
		if !found {
			// fmt.Println("没有a标签了")

			if len(product) == 0 {
				fmt.Println("没有数据打印")
				break
			}
			fmt.Println("正在导出中")
			f := excelize.NewFile()
			// 创建一个新的Sheet
			index := f.NewSheet("Sheet1")

			// 设置单元格的值
			f.SetCellValue("Sheet1", "A1", "商品")
			f.SetCellValue("Sheet1", "B1", "货号")
			f.SetCellValue("Sheet1", "C1", "链接")
			f.SetCellValue("Sheet1", "D1", "售价")
			f.SetCellValue("Sheet1", "E1", "市场价")
			f.SetCellValue("Sheet1", "F1", "尺码")
			// 将产品数据写入Excel

			row := 2
			for _, product1 := range product {
				// 将Size字段转换为切片
				sizes := strings.Split(strings.TrimSuffix(product1.Size, ","), ",")

				// 写入每个尺寸
				for _, size := range sizes {
					f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), product1.Name)
					f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), product1.Code)
					f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), product1.Link)
					f.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), product1.SalePrice)
					f.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), product1.MarketPrice)
					f.SetCellValue("Sheet1", fmt.Sprintf("F%d", row), size)
					row++
				}
			}

			// 设置默认打开的Sheet
			f.SetActiveSheet(index)

			t := time.Now()
			filename := t.Format("20060102150405") + ".xlsx"

			// 保存文件
			err := f.SaveAs(filename)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("导出成功" + filename)
			break
		}
	}
	fmt.Println("结束：************")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
