package service

import (
	"code.aliyun.com/zmdev/wechat_rank/errors"
	"code.aliyun.com/zmdev/wechat_rank/model"
	"code.aliyun.com/zmdev/wechat_rank/pkg/qingbo"
	"github.com/jinzhu/gorm"
)

type wechatService struct {
	model.WechatStore
	*qingbo.WxAccount
	*qingbo.WxGroup
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

func convert2WechatModel(account *qingbo.AccountData, wechat *model.Wechat) {
	wechat.WxName = account.WxName
	wechat.VerifyName = account.VerifyName
	wechat.WxLogo = account.WxLogo
	wechat.WxNote = account.WxNote
	wechat.WxQrcode = account.WxQrcode
	wechat.WxVip = account.WxVip
	wechat.WxNickname = account.WxNickname
}

func NewWechatService(ws model.WechatStore, wxAccount *qingbo.WxAccount, wxGroup *qingbo.WxGroup) model.WechatService {
	return &wechatService{ws, wxAccount, wxGroup}
}
