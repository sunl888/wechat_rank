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

func headers() []string {
	return []string{
		"公众号", "帐号名", "文章总数", "头条文章总数", "阅读总数", "平均阅读数", "点赞总数", "平均点赞数",
		"头条文章阅读量", "头条文章点赞数", "最大阅读数", "最大点赞数", "点赞率", "WCI", "总排名",
	}
}

func (*Export) ExportData(ctx *gin.Context) {
	l := struct {
		RankId     int64 `json:"rank_id" form:"rank_id"`
		CategoryId int64 `json:"category_id" form:"category_id"`
		Top        int   `json:"top" form:"top"`
	}{}
	if err := ctx.ShouldBind(&l); err != nil {
		_ = ctx.Error(errors.BindError(err))
		return
	}
	fileName := time.Now().Format(service.DATE_FORMAT)
	ctx.Writer.Header().Add("Content-Disposition", "attachment;filename=\""+fileName+".xlsx\"")
	ctx.Writer.Header().Add("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet;charset=utf-8")
	ctx.Writer.Header().Set("X-Content-Type-Options", "nosniff")

	builder := xlsx.NewStreamFileBuilder(ctx.Writer)
	cellStrType := xlsx.CellTypeString
	cellNumType := xlsx.CellTypeNumeric
	ct := []*xlsx.CellType{
		cellStrType.Ptr(), cellStrType.Ptr(),
		cellNumType.Ptr(), cellNumType.Ptr(), cellNumType.Ptr(),
		cellNumType.Ptr(), cellNumType.Ptr(), cellNumType.Ptr(), cellNumType.Ptr(),
		cellNumType.Ptr(), cellNumType.Ptr(), cellNumType.Ptr(), cellNumType.Ptr(), cellNumType.Ptr(),
	}
	if err := builder.AddSheet("排名", headers(), ct); err != nil {
		_ = ctx.Error(err)
		return
	}
	streamFile, err := builder.Build()
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	defer streamFile.Close()
	ranks, count, err := service.RankDetail(ctx, l.RankId, l.CategoryId, l.Top, 0)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	var records [][]string
	for i := 0; i < int(count); i++ {
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
			fmt.Sprintf("%d", i+1),                   // 总排名
		})
	}
	_ = streamFile.WriteAll(records)
	streamFile.Flush()
	return
}

func NewExport() *Export {
	return &Export{}
}
