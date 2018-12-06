package handler

import (
	"code.aliyun.com/zmdev/wechat_rank/errors"
	"github.com/gin-gonic/gin"
)

type Export struct {
}

func (*Export) ExportData(ctx *gin.Context) {
	//headers := []string{"groupName", "groupID", "groupPath", "groupOwner"}
	//ctx.Writer.Header().Add("Content-Disposition", "attachment;filename=\""+"test"+".xlsx\"")
	//ctx.Writer.Header().Add("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet;charset=utf-8")
	//ctx.Writer.Header().Set("X-Content-Type-Options", "nosniff")
	//b := xlsx.NewStreamFileBuilder(ctx.Writer)
	//cellType := xlsx.CellTypeString
	//ct := []*xlsx.CellType{cellType.Ptr(), cellType.Ptr(), cellType.Ptr(), cellType.Ptr()}
	//
	//if err := b.AddSheet("Sheet1", headers, ct); err != nil {
	//	fmt.Println("end")
	//	return
	//}
	//sf, err := b.Build()
	//if err != nil {
	//	fmt.Println("end")
	//	return
	//}
	//defer sf.Close()
	//for i := 0; i < 100; i++ {
	//	_ = sf.Write([]string{"GroupName", "GroupID", "PathName", "OwnerAccount"})
	//	sf.Flush()
	//}
	_ = ctx.Error(errors.BadRequest("test", nil))
	return
}

func NewExport() *Export {
	return &Export{}
}
