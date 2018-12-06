package service

import (
	"code.aliyun.com/zmdev/wechat_rank/errors"
	"code.aliyun.com/zmdev/wechat_rank/model"
	"code.aliyun.com/zmdev/wechat_rank/utils"
	"github.com/jinzhu/gorm"
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
			convert2WechatModel(wechatData, wechat)
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

func convert2WechatModel(account *utils.AccountData, wechat *model.Wechat) {
	wechat.WxName = account.WxName
	wechat.VerifyName = account.VerifyName
	wechat.WxLogo = account.WxLogo
	wechat.WxNote = account.WxNote
	wechat.WxQrcode = account.WxQrcode
	wechat.WxVip = account.WxVip
	wechat.WxNickname = account.WxNickname
}

func NewWechatService(ws model.WechatStore, client *utils.OfficialAccount) model.WechatService {
	return &wechatService{ws, client}
}
