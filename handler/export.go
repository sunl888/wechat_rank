package handler

import (
	"code.aliyun.com/zmdev/wechat_rank/errors"
	"code.aliyun.com/zmdev/wechat_rank/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
	"time"
)

type Export struct {
}

const (
	AccountRank = iota
	ArticleRank
)

func headers(n int) []string {
	h := make([][]string, 2)
	h[AccountRank] = []string{
		"公众号", "帐号名", "文章总数", "头条文章总数", "阅读总数", "平均阅读数", "点赞总数", "平均点赞数",
		"头条文章阅读量", "头条文章点赞数", "最大阅读数", "最大点赞数", "点赞率", "WCI", "总排名",
	}
	h[ArticleRank] = []string{
		"公众号", "帐号名", "标题", "摘要", "URL", "发布时间", "阅读数", "点赞数", "文章序号",
	}
	return h[n]
}

func setHeaders(ctx *gin.Context) {
	fileName := time.Now().Format(service.DATE_FORMAT)
	ctx.Writer.Header().
		Add("Content-Disposition", "attachment;filename=\""+fileName+".xlsx\"")
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
	setHeaders(ctx)
	builder := xlsx.NewStreamFileBuilder(ctx.Writer)
	cellStrType := xlsx.CellTypeString
	cellNumType := xlsx.CellTypeNumeric
	ct := []*xlsx.CellType{
		cellStrType.Ptr(), cellStrType.Ptr(),
		cellStrType.Ptr(), cellStrType.Ptr(), cellStrType.Ptr(),
		cellStrType.Ptr(), cellNumType.Ptr(), cellNumType.Ptr(), cellNumType.Ptr(),
	}
	if err := builder.AddSheet("文章排名", headers(ArticleRank), ct); err != nil {
		_ = ctx.Error(err)
		return
	}
	streamFile, err := builder.Build()
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	defer streamFile.Close()
	articles, _, err := service.ArticleRank(ctx, l.StartDate, l.EndDate, l.CategoryId, 0, l.Top)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	var records [][]string
	for i := 0; i < len(articles); i++ {
		records = append(records, []string{
			articles[i].WxNickname, // 公众号
			articles[i].WxName,     // 帐号名
			articles[i].Title,      // 标题
			articles[i].Desc,       // 摘要
			articles[i].Url,        // Url
			articles[i].PublishedAt.Format("2006-01-02 15:04:05"), // 发布时间
			fmt.Sprintf("%d", articles[i].ReadCount),              // 阅读数
			fmt.Sprintf("%d", articles[i].LikeCount),              // 点赞数
			fmt.Sprintf("%d", articles[i].Top),                    // 文章序号
		})
	}
	_ = streamFile.WriteAll(records)
	streamFile.Flush()
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
	setHeaders(ctx)
	builder := xlsx.NewStreamFileBuilder(ctx.Writer)
	cellStrType := xlsx.CellTypeString
	cellNumType := xlsx.CellTypeNumeric
	ct := []*xlsx.CellType{
		cellStrType.Ptr(), cellStrType.Ptr(),
		cellNumType.Ptr(), cellNumType.Ptr(), cellNumType.Ptr(),
		cellNumType.Ptr(), cellNumType.Ptr(), cellNumType.Ptr(), cellNumType.Ptr(),
		cellNumType.Ptr(), cellNumType.Ptr(), cellNumType.Ptr(), cellNumType.Ptr(), cellNumType.Ptr(),
	}
	if err := builder.AddSheet("公众号排名", headers(AccountRank), ct); err != nil {
		_ = ctx.Error(err)
		return
	}
	streamFile, err := builder.Build()
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	defer streamFile.Close()
	ranks, _, err := service.RankDetail(ctx, l.RankId, l.CategoryId, l.Top, 0)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	var records [][]string
	for i := 0; i < len(ranks); i++ {
		records = append(records, []string{
			ranks[i].WxNickname,                      // 公众号
			ranks[i].WxName,                          // 帐号名
			fmt.Sprintf("%d", ranks[i].ArticleCount), // 文章总数
			fmt.Sprintf("%d", ranks[i].TopCount),     // 头条文章总数
			fmt.Sprintf("%d", ranks[i].ReadCount),    // 阅读总数
			fmt.Sprintf("%d", ranks[i].AvgReadCount), // 平均阅读数
			fmt.Sprintf("%d", ranks[i].LikeCount),    // 点赞总数
			fmt.Sprintf("%d", ranks[i].AvgLikeCount), // 平均点赞数
			fmt.Sprintf("%d", ranks[i].TopReadCount), // 头条文章阅读量
			fmt.Sprintf("%d", ranks[i].TopLikeCount), // 头条文章点赞数
			fmt.Sprintf("%d", ranks[i].MaxReadCount), // 最大阅读数
			fmt.Sprintf("%d", ranks[i].MaxLikeCount), // 最大点赞数
			fmt.Sprintf("%.5f", ranks[i].LikeRate),   // 点赞率
			fmt.Sprintf("%.5f", ranks[i].Wci),        // WCI
			fmt.Sprintf("%d", ranks[i].TotalRank),    // 总排名
		})
	}
	_ = streamFile.WriteAll(records)
	streamFile.Flush()
	return
}

func NewExport() *Export {
	return &Export{}
}
