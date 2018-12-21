package main

import (
	"fmt"
	"github.com/tealeg/xlsx"
)

func headers(n int) []string {
	h := make([][]string, 2)
	h[0] = []string{
		"公众号", "帐号名", "文章总数", "头条文章总数", "阅读总数", "平均阅读数", "点赞总数", "平均点赞数",
		"头条文章阅读量", "头条文章点赞数", "最大阅读数", "最大点赞数", "点赞率", "WCI", "总排名",
	}
	h[1] = []string{
		"公众号", "帐号名", "标题", "摘要", "URL", "发布时间", "阅读数", "点赞数", "文章序号",
	}
	return h[n]
}

func writeHeader(sheet *xlsx.Sheet) {
	r := sheet.AddRow()
	c := r.AddCell()
	c.Merge(14, 1)
	c.Value = "淮南榜"
}
func writeFooter(sheet *xlsx.Sheet) {
	r := sheet.AddRow()
	c := r.AddCell()
	c.Merge(14, 1)
	c.SetValue("温馨提示：排名算法来自“清博大数据")
}
func main() {
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
	writeHeader(sheet)

	row = sheet.AddRow()
	cell = row.AddCell()
	//cell.Merge(3, 0)
	cell.Value = "I am a cell!"
	writeFooter(sheet)
	err = file.Save("MyXLSXFile.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
}
