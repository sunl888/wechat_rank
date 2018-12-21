package handler

import (
	"code.aliyun.com/zmdev/wechat_rank/errors"
	"code.aliyun.com/zmdev/wechat_rank/service"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"strconv"
)

type Export struct {
}

const (
	ArticleHeader = iota
	AccountHeader
	Table = "Sheet1"
)

func setHeader(headIndex int, xlsx *excelize.File, tableName string) {
	headers := map[int]map[string]string{
		AccountHeader: {
			"A3": "公众号", "B3": "帐号名", "C3": "文章总数",
			"D3": "头条文章总数", "E3": "阅读总数", "F3": "平均阅读数",
			"G3": "点赞总数", "H3": "平均点赞数", "I3": "头条文章阅读量",
			"J3": "头条文章点赞数", "K3": "最大阅读数", "L3": "最大点赞数",
			"M3": "点赞率", "N3": "WCI", "O3": "总排名",
		},
		ArticleHeader: {
			"A3": "公众号", "B3": "帐号名", "C3": "标题",
			"D3": "摘要", "E3": "URL", "F3": "发布时间",
			"G3": "阅读数", "H3": "点赞数", "I3": "文章序号",
		},
	}
	if _, ok := headers[headIndex]; !ok {
		fmt.Println("指定的 index 不存在")
		return
	}
	// 水平垂直居中,加粗, 字号:14
	titleId, _ := xlsx.NewStyle(`{"font":{"size":14,"family":"微软雅黑"},"alignment":{"horizontal":"center","vertical":"center"}}`)
	headerId, _ := xlsx.NewStyle(`{"font":{"family":"微软雅黑"}}`)

	xlsx.SetCellValue(tableName, "A1", "淮南政务微信排行榜")
	if headIndex == ArticleHeader {
		xlsx.MergeCell(tableName, "A1", "I2")
		xlsx.SetCellStyle(tableName, "A1", "I2", titleId)
		xlsx.SetCellStyle(tableName, "A3", "I3", headerId)
	} else if headIndex == AccountHeader {
		xlsx.MergeCell(tableName, "A1", "O2")
		xlsx.SetCellStyle(tableName, "A1", "O2", titleId)
		xlsx.SetCellStyle(tableName, "A3", "O3", headerId)
	}
	for k, v := range headers[headIndex] {
		xlsx.SetCellValue(tableName, k, v)
	}
}

func setFooter(headIndex int, row string, xlsx *excelize.File, tableName string) {
	xlsx.SetCellValue(tableName, "A"+row, "清博指数是透明、学术、权威的第三方评价，数据和公式可公开查询：www.gsdata.cn/site/usage")
	styleId, _ := xlsx.NewStyle(`{"font":{"italic":true}}`)
	if headIndex == ArticleHeader {
		xlsx.MergeCell(tableName, "A"+row, "I"+row)
		xlsx.SetCellStyle(tableName, "A"+row, "I"+row, styleId)
	} else if headIndex == AccountHeader {
		xlsx.MergeCell(tableName, "A"+row, "O"+row)
		xlsx.SetCellStyle(tableName, "A"+row, "O"+row, styleId)
	}
}

func setOutStreamHeaders(ctx *gin.Context, filename string) {
	ctx.Writer.Header().
		Add("Content-Disposition", "attachment;filename=\""+filename+".xlsx\"")
	ctx.Writer.Header().
		Add("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet;charset=utf-8")
	ctx.Writer.Header().
		Set("X-Content-Type-Options", "nosniff")
}

func (*Export) ArticleRank(ctx *gin.Context) {
	l := struct {
		StartDate  string `json:"start_date" form:"start_date"`
		EndDate    string `json:"end_date" form:"end_date"`
		CategoryId int64  `json:"category_id" form:"category_id"`
		Top        int    `json:"top" form:"top"`
	}{}
	if err := ctx.ShouldBind(&l); err != nil {
		_ = ctx.Error(errors.BindError(err))
		return
	}
	setOutStreamHeaders(ctx, fmt.Sprintf("文章排名榜单(%s~%s)", l.StartDate, l.EndDate))
	articles, _, err := service.ArticleRank(ctx, l.StartDate, l.EndDate, l.CategoryId, 0, l.Top)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	xlsx := excelize.NewFile()
	sheetIndex := xlsx.NewSheet(Table)
	setHeader(ArticleHeader, xlsx, Table)
	for i := 0; i < len(articles); i++ {
		r := strconv.Itoa(i + 4)
		xlsx.SetCellValue(Table, "A"+r, articles[i].WxNickname)
		xlsx.SetCellValue(Table, "B"+r, articles[i].WxName)
		xlsx.SetCellValue(Table, "C"+r, articles[i].Title)
		xlsx.SetCellValue(Table, "D"+r, articles[i].Desc)
		xlsx.SetCellValue(Table, "E"+r, articles[i].Url)
		xlsx.SetCellValue(Table, "F"+r, articles[i].PublishedAt.Format("2006-01-02 15:04:05"))
		xlsx.SetCellValue(Table, "G"+r, fmt.Sprintf("%d", articles[i].ReadCount))
		xlsx.SetCellValue(Table, "H"+r, fmt.Sprintf("%d", articles[i].LikeCount))
		xlsx.SetCellValue(Table, "I"+r, fmt.Sprintf("%d", articles[i].Top))
	}
	//TODO 自动调整列宽

	nextRow := strconv.Itoa(len(articles) + 4)
	setFooter(ArticleHeader, nextRow, xlsx, Table)
	xlsx.SetActiveSheet(sheetIndex)
	err = xlsx.Write(ctx.Writer)
	if err != nil {
		_ = ctx.Error(err)
	}
	return
}

func (*Export) AccountRank(ctx *gin.Context) {
	l := struct {
		RankId     int64 `json:"rank_id" form:"rank_id"`
		CategoryId int64 `json:"category_id" form:"category_id"`
		Top        int   `json:"top" form:"top"`
	}{}
	if err := ctx.ShouldBind(&l); err != nil {
		_ = ctx.Error(errors.BindError(err))
		return
	}
	rank, _ := service.RankLoad(ctx, l.RankId)
	t := map[string]string{
		"week":  "周榜",
		"month": "月榜",
		"year":  "年榜",
	}
	setOutStreamHeaders(ctx, fmt.Sprintf("微信排名-%s(%s~%s)", t[rank.Period], rank.StartDate, rank.EndDate))
	ranks, _, err := service.RankDetail(ctx, l.RankId, l.CategoryId, l.Top, 0)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	xlsx := excelize.NewFile()
	sheetIndex := xlsx.NewSheet(Table)
	setHeader(AccountHeader, xlsx, Table)
	for i := 0; i < len(ranks); i++ {
		r := strconv.Itoa(i + 4)
		xlsx.SetCellValue(Table, "A"+r, ranks[i].WxNickname)
		xlsx.SetCellValue(Table, "B"+r, ranks[i].WxName)
		xlsx.SetCellValue(Table, "C"+r, ranks[i].ArticleCount)
		xlsx.SetCellValue(Table, "D"+r, ranks[i].TopCount)
		xlsx.SetCellValue(Table, "E"+r, ranks[i].ReadCount)
		xlsx.SetCellValue(Table, "F"+r, ranks[i].AvgReadCount)
		xlsx.SetCellValue(Table, "G"+r, ranks[i].LikeCount)
		xlsx.SetCellValue(Table, "H"+r, ranks[i].AvgLikeCount)
		xlsx.SetCellValue(Table, "I"+r, ranks[i].TopReadCount)
		xlsx.SetCellValue(Table, "J"+r, ranks[i].TopLikeCount)
		xlsx.SetCellValue(Table, "K"+r, ranks[i].MaxReadCount)
		xlsx.SetCellValue(Table, "L"+r, ranks[i].MaxLikeCount)
		xlsx.SetCellValue(Table, "M"+r, fmt.Sprintf("%.5f", ranks[i].LikeRate))
		xlsx.SetCellValue(Table, "N"+r, fmt.Sprintf("%.5f", ranks[i].Wci))
		xlsx.SetCellValue(Table, "O"+r, ranks[i].TotalRank)
	}
	nextRow := strconv.Itoa(len(ranks) + 4)
	setFooter(AccountHeader, nextRow, xlsx, Table)
	xlsx.SetActiveSheet(sheetIndex)
	err = xlsx.Write(ctx.Writer)
	if err != nil {
		_ = ctx.Error(err)
	}
	return
}

func NewExport() *Export {
	return &Export{}
}
