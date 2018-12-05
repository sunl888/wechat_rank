package service

import (
	"code.aliyun.com/zmdev/wechat_rank/errors"
	"code.aliyun.com/zmdev/wechat_rank/model"
	"code.aliyun.com/zmdev/wechat_rank/utils"
	"github.com/jinzhu/gorm"
	"time"
)

type wechatService struct {
	model.WechatStore
	client *utils.OfficialAccount
}

func (w *wechatService) WechatCreate(wechat *model.Wechat) error {
	tmpWechat, err := w.WechatStore.WechatLoad(wechat.WxName)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			wechatResp, err := w.client.GetAccount(wechat.WxName)
			if err != nil {
				return err
			}
			if len(wechatResp.Data) <= 0 {
				return errors.BadRequest("公众号不存在", nil)
			}
			wechatData := wechatResp.Data[0]

			lastDate := time.Now().AddDate(0, 0, -1).Format(DATE_FORMAT)
			rankDay, err := w.client.GetRankDays(wechat.WxName, lastDate)
			if err != nil {
				return err
			}
			nickname := ""
			if rankDay.Data[0] != nil {
				nickname = rankDay.Data[0].WxNickname
			}
			convert2WechatModel(wechatData, wechat, nickname)
			err = w.WechatStore.WechatCreate(wechat)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	wechat = tmpWechat
	return nil
}

func convert2WechatModel(account *utils.AccountData, wechat *model.Wechat, nickname string) {
	wechat.WxName = account.WxName
	wechat.VerifyName = account.VerifyName
	wechat.WxLogo = account.WxLogo
	wechat.WxNote = account.WxNote
	wechat.WxQrcode = account.WxQrcode
	wechat.WxVip = account.WxVip
	if nickname == "" {
		wechat.WxNickname = wechat.VerifyName
	} else {
		wechat.WxNickname = nickname
	}
}

func NewWechatService(ws model.WechatStore, client *utils.OfficialAccount) model.WechatService {
	return &wechatService{ws, client}
}
