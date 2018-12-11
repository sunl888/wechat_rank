package qingbo

import (
	"code.aliyun.com/zmdev/wechat_rank/errors"
	"encoding/json"
)

type WxGroup struct {
	*client
	groupId string
}

// 这是个杀马特结构体(居然返回的json有大写有小写)
type WxGroupData struct {
	WxNickname string `json:"WX_NICKNAME"`
	WxName     string `json:"WX_NAME"`
	WxBiz      string `json:"WX_BIZ"`
	NicknameId string `json:"ID"`
}

type WxGroupResponse struct {
	Data        WxGroupData `json:"data"`
	Url         string      `json:"url"`
	Application string      `json:"application"`
}

// http://api.gsdata.cn/weixin/v1/groups/acounts
// 添加公众号到清博大数据的自定义分组里,这样有助于清博大数据监控此公众号
func (w *WxGroup) AddWx2Group(wxName string) (wxGroupResponse *WxGroupResponse, err error) {
	resp, err := w.client.post("groups/acounts", map[string]string{
		"wx_name":  wxName,
		"group_id": w.groupId,
	}, "weixin")
	if err != nil {
		return nil, err
	}
	errResp := &ErrorResponse{}
	err = json.Unmarshal([]byte(resp), errResp)
	if err != nil {
		return nil, err
	}
	if errResp.Status != 201 && errResp.Status != 0 {
		return nil, errors.QingboError(errResp.Name, errResp.Message, errResp.Code, errResp.Status)
	}
	wxGroupResponse = &WxGroupResponse{}
	err = json.Unmarshal([]byte(resp), wxGroupResponse)
	if err != nil {
		return nil, err
	}
	return wxGroupResponse, nil
}

func NewWxGroup(client *client, groupId string) *WxGroup {
	return &WxGroup{client, groupId}
}
