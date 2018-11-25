package service

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
	"code.aliyun.com/zmdev/wechat_rank/utils"
	"encoding/json"
	"fmt"
)

type wechatService struct {
	model.WechatStore
	qingboClient *utils.QingboClient
}

type WechatMaterialResp struct {
	VerifyName string `json:"verify_name"`
	WxName     string `json:"wx_name"`
	AddTime    string `json:"add_time"`
	WxVip      string `json:"wx_vip"`
	WxNote     string `json:"wx_note"`
	WxLogo     string `json:"wx_logo"`
	Wci        int64  `json:"wci"`
	WxNickname string `json:"wx_nickname"`
	NicknameId string `json:"nickname_id"`
	WxQrcode   string `json:"wx_qrcode"`
}

func (w *wechatService) WechatCreate(wechat *model.Wechat) error {
	resp, err := w.qingboClient.Get("users", "wx_name="+wechat.Name)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)

	fullResp := &struct {
		WechatResp struct {
			MapData []*WechatMaterialResp `json:"data"`
		}
	}{}
	err = json.Unmarshal([]byte(resp), fullResp)
	if err != nil {
		return err
	}
	fmt.Println(fullResp.WechatResp.MapData[0])
	return nil
	//return w.WechatStore.WechatCreate(wechat)
}

func NewWechatService(ws model.WechatStore, qc *utils.QingboClient) model.WechatService {
	return &wechatService{ws, qc}
}
