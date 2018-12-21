package handler

import (
	"code.aliyun.com/zmdev/wechat_rank/errors"
	"code.aliyun.com/zmdev/wechat_rank/handler/middleware"
	"code.aliyun.com/zmdev/wechat_rank/model"
	"code.aliyun.com/zmdev/wechat_rank/service"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
)

type Category struct {
}

type CategoryResp struct {
	Id          int64     `json:"id"`
	Title       string    `json:"title"`
	IsPrivate   bool      `json:"is_private"`
	WechatCount int       `json:"wechat_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (*Category) Update(ctx *gin.Context) {
	cId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		_ = ctx.Error(errors.BadRequest("id 格式不正确", nil))
		return
	}
	l := struct {
		Title     string `json:"title" form:"title"`
		IsPrivate bool   `json:"is_private" form:"is_private"`
	}{}
	if err := ctx.ShouldBind(&l); err != nil {
		_ = ctx.Error(errors.BindError(err))
		return
	}
	err = service.CategoryUpdate(ctx, &model.Category{
		Id:        cId,
		Title:     l.Title,
		IsPrivate: l.IsPrivate,
	})
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.Status(204)
}

func (*Category) Create(ctx *gin.Context) {
	l := struct {
		Title     string `json:"title" form:"title"`
		IsPrivate bool   `json:"is_private"`
	}{}
	if err := ctx.ShouldBind(&l); err != nil {
		_ = ctx.Error(errors.BindError(err))
		return
	}
	category := &model.Category{
		Title:     l.Title,
		IsPrivate: l.IsPrivate,
	}
	err := service.CategoryCreate(ctx, category)
	if err != nil {
		_ = ctx.Error(errors.BindError(err))
		return
	}
	ctx.JSON(200, convert2CategoryResp(category, ctx))
}

func (*Category) List(ctx *gin.Context) {
	isLogin := middleware.CheckLogin(ctx)
	categories, err := service.CategoryList(ctx, isLogin)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(200, convert2CategoriesResp(categories, ctx))
}

func (*Category) Show(ctx *gin.Context) {
	cId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		_ = ctx.Error(errors.BadRequest("id 格式不正确", nil))
		return
	}
	category, err := service.CategoryLoad(ctx, cId)
	if err != nil {
		_ = ctx.Error(errors.BadRequest("分类不存在", nil))
		return
	}
	ctx.JSON(200, convert2CategoryResp(category, ctx))
}

func (c *Category) Delete(ctx *gin.Context) {
	var err error
	cId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		_ = ctx.Error(errors.BadRequest("id 格式不正确", nil))
		return
	}
	_, err = service.CategoryLoad(ctx, cId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = errors.BadRequest("分类不存在", nil)
		}
		_ = ctx.Error(err)
		return
	}
	_, count, err := service.WechatListByCategory(ctx, cId, 1, 0)
	if count > 0 {
		_ = ctx.Error(errors.BadRequest("该分类下有公众号,不能删除", nil))
		return
	}
	err = service.CategoryDelete(ctx, cId)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.Status(204)
}

func convert2CategoryResp(c *model.Category, ctx *gin.Context) *CategoryResp {
	_, count, _ := service.WechatListByCategory(ctx, c.Id, 0, 0)
	return &CategoryResp{
		Id:          c.Id,
		Title:       c.Title,
		IsPrivate:   c.IsPrivate,
		WechatCount: int(count),
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

func convert2CategoriesResp(cs []*model.Category, ctx *gin.Context) []*CategoryResp {
	categoriesResp := make([]*CategoryResp, 0, len(cs))
	for _, c := range cs {
		categoriesResp = append(categoriesResp, convert2CategoryResp(c, ctx))
	}
	return categoriesResp
}

func NewCategory() *Category {
	return &Category{}
}
