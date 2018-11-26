package service

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
	"code.aliyun.com/zmdev/wechat_rank/utils"
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
	//resp, err := w.qingboClient.Get("users", "wx_name="+wechat.WxName)
	//if err != nil {
	//	return err
	//}
	//	resp := `{
	//	"data": [{
	//		"verify_name": "淮南师范学院",
	//		"wx_name": "hnnu1958",
	//		"add_time": "2016-08-26 16:37:34",
	//		"wx_vip": "认证",
	//		"wx_note": "淮南师范学院官方公众平台",
	//		"wx_logo": "http://wx.qlogo.cn/mmhicuiazg1JfFib0CuP4gqPh8ghac8XgliaAibCILOECA/132",
	//		"wx_account_tags": null,
	//		"wci": 369.21,
	//		"wx_nickname": "淮南师范学院",
	//		"nickname_id": "6050891",
	//		"wx_qrcode": "https://mp.weixin.qq.com/mp/?scene=10000001&size=100&__biz=MzI2MTMyMTk4NA==&mid=2247491211&idx=1&sn=630b51e9367d76593de1074e474e97d4&scene=0",
	//		"types": null,
	//		"wx_biz": "MzI2MTMyMTk4NA=="
	//	}],
	//	"application": "4ee4d04d-9312-466a-ae86-3bfb85aa3fda",
	//	"url": "weixin/v1/users"
	//}`
	//	wechatResp := &WechatResponse{}
	//	err := json.Unmarshal([]byte(resp), wechatResp)
	//	if err != nil {
	//		return err
	//	}
	wechatResp, err := w.client.GetOfficialAccount(wechat.WxName)
	if err != nil {
		return err
	}
	convert2WechatModel(wechatResp, wechat)
	return w.WechatStore.WechatCreate(wechat)
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
