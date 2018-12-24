package service

import (
	"code.aliyun.com/zmdev/wechat_rank/errors"
	"code.aliyun.com/zmdev/wechat_rank/model"
	"code.aliyun.com/zmdev/wechat_rank/pkg/qingbo"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
)

type wechatService struct {
	model.WechatStore
	*qingbo.WxAccount
	*qingbo.WxGroup
	*redis.Client
}

const (
	LimitRequest = 3
	IntervalTime = time.Second * 3600 * 24
)

func (w *wechatService) WechatSync(wxName string) (wechat *model.Wechat, err error) {
	if w.Client.Exists(wxName).Val() > 0 {
		val, err := strconv.Atoi(w.Client.Get(wxName).Val())
		if err != nil {
			return nil, err
		}
		if val >= LimitRequest {
			return nil, errors.BadRequest(fmt.Sprintf("拒绝更新,%.3f小时内最多只能请求%d次该接口", IntervalTime.Hours(), LimitRequest), nil)
		} else {
			w.Client.Incr(wxName)
		}
	} else {
		w.Client.Set(wxName, 1, IntervalTime)
	}
	resp, err := w.WxAccount.GetAccount(wxName)
	if err != nil {
		return nil, err
	}
	wechatData := resp.Data[0]
	wechat = new(model.Wechat)
	convert2WechatModel(wechatData, wechat)
	err = w.WechatStore.WechatUpdate(wechat)
	if err != nil {
		return nil, err
	}
	return
}

func (w *wechatService) WechatCreate(wechat *model.Wechat) error {
	_, err := w.WechatStore.WechatLoad(wechat.WxName)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			wechatResp, err := w.WxAccount.GetAccount(wechat.WxName)
			if err != nil {
				return err
			}
			if len(wechatResp.Data) <= 0 {
				return errors.BadRequest("公众号不存在", nil)
			}
			wechatData := wechatResp.Data[0]
			convert2WechatModel(wechatData, wechat)
			err = w.WechatStore.WechatCreate(wechat)
			if err != nil {
				return err
			}
			// 添加公众号到自定义榜单上去
			_, err = w.WxGroup.AddWx2Group(wechatData.WxName)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		return errors.BadRequest("公众号已经存在", nil)
	}
	return nil
}

func WechatCreate(ctx *gin.Context, wechat *model.Wechat) error {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.WechatCreate(wechat)
	}
	return ServiceError
}

func WechatSync(ctx *gin.Context, wxName string) (wechat *model.Wechat, err error) {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.WechatSync(wxName)
	}
	return nil, ServiceError
}

func WechatUpdate(ctx *gin.Context, wechat *model.Wechat) error {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.WechatUpdate(wechat)
	}
	return ServiceError
}

func WechatLoad(ctx *gin.Context, wxName string) (wechat *model.Wechat, err error) {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.WechatLoad(wxName)
	}
	return nil, ServiceError
}

func WechatList(ctx *gin.Context, limit, offset int) (wechats []*model.Wechat, count int64, err error) {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.WechatList(limit, offset)
	}
	return nil, 0, ServiceError
}

func WechatSearch(ctx *gin.Context, keyword string, limit, offset int) (wechats []*model.WechatAndCategory, count int64, err error) {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.WechatSearch(keyword, limit, offset)
	}
	return nil, 0, ServiceError
}

func WechatDelete(ctx *gin.Context, id int64) error {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.WechatDelete(id)
	}
	return ServiceError
}

func WechatListByCategory(ctx *gin.Context, cId int64, limit, offset int) (wechats []*model.Wechat, count int64, err error) {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.WechatListByCategory(cId, limit, offset)
	}
	return nil, 0, ServiceError
}

func WechatCountByCategory(ctx *gin.Context, cId int64) (count int64, err error) {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.WechatCountByCategory(cId)
	}
	return 0, ServiceError
}

func convert2WechatModel(account *qingbo.AccountData, wechat *model.Wechat) {
	wechat.WxName = account.WxName
	wechat.VerifyName = account.VerifyName
	wechat.WxLogo = account.WxLogo
	wechat.WxNote = account.WxNote
	wechat.WxQrcode = account.WxQrcode
	wechat.WxVip = account.WxVip
	wechat.WxNickname = account.WxNickname
}

func NewWechatService(ws model.WechatStore, wxAccount *qingbo.WxAccount, wxGroup *qingbo.WxGroup, redisClient *redis.Client) model.WechatService {
	return &wechatService{ws, wxAccount, wxGroup, redisClient}
}
