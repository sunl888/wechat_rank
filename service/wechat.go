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

type Data struct {
	VerifyName    string   `json:"verify_name"`
	WxName        string   `json:"wx_name"`
	AddTime       string   `json:"add_time"`
	WxVip         string   `json:"wx_vip"`
	WxNote        string   `json:"wx_note"`
	WxLogo        string   `json:"wx_logo"`
	Wci           float64  `json:"wci,omitempty"`
	WxNickname    string   `json:"wx_nickname,omitempty"`
	NicknameId    string   `json:"nickname_id"`
	WxQrcode      string   `json:"wx_qrcode"`
	WxBiz         string   `json:"wx_biz"`
	WxAccountTags []string `json:"wx_account_tags,omitempty"`
}

type WechatResponse struct {
	DataResp    []Data `json:"data"`
	Url         string `json:"url"`
	Application string `json:"application"`
}

func (w *wechatService) WechatCreate(wechat *model.Wechat) error {
	tmpWechat, err := w.WechatStore.WechatLoad(wechat.WxName)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			wechatResp, err := w.client.GetOfficialAccount(wechat.WxName)
			if err != nil {
				return err
			}
			if len(wechatResp.DataResp) <= 0 {
				return errors.BadRequest("公众号不存在", nil)
			}
			convert2WechatModel(wechatResp, wechat)
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

func convert2WechatModel(response *utils.AccountResponse, wechat *model.Wechat) {
	wechat.WxName = response.DataResp[0].WxName
	wechat.VerifyName = response.DataResp[0].VerifyName
	wechat.WxLogo = response.DataResp[0].WxLogo
	wechat.WxNote = response.DataResp[0].WxNote
	wechat.WxQrcode = response.DataResp[0].WxQrcode
	wechat.WxVip = response.DataResp[0].WxVip
}

func NewWechatService(ws model.WechatStore, client *utils.OfficialAccount) model.WechatService {
	return &wechatService{ws, client}
}
