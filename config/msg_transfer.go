package config

import (
	"fmt"
	"time"

	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/common"
	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/pkg/util"
)

// SendTransfer 发送转账消息
func (c *Context) SendTransfer(msg MsgTransfer) error {

	return c.SendMessage(&MsgSendReq{
		Header: MsgHeader{
			RedDot: 1,
		},
		ChannelID:   msg.ToUID,
		ChannelType: common.ChannelTypePerson.Uint8(),
		FromUID:     msg.FromUID,
		Payload: []byte(util.ToJson(map[string]interface{}{
			"record_no": msg.RecordNo,
			"remark":    msg.Remark,
			"amount":    msg.Amount,
			"type":      common.Transfer,
			"status":    msg.Status,
		})),
	})
}

// SendTransferRecover 发送转账回收消息
func (c *Context) SendTransferRecover(msg MsgTransferRecover) error {
	return c.SendTradeSystemNotifyTemplate(MsgTradeSystemNotifyTemplate{
		ChannelID: msg.Receiver,
		// ChannelID:      msg.Creater,
		ChannelType:    common.ChannelTypePerson.Uint8(),
		LeftTitle:      "转账退款到账通知",
		CenterTitle:    fmt.Sprintf("¥%0.2f", util.CentToYuan(msg.Amount)),
		CenterSubtitle: "退款金额",
		URLTitle:       "查看详情",
		Notice:         "转账退款",
		Attrs: map[string]string{
			"退款方式":         "退回零钱",
			"退款原因":         "转账超过24小时未被领取",
			"到账时间":         util.ToyyyyMMddHHmmss(time.Unix(msg.ExpiredAt, 0)),
			"备注":           "退款金额已到账",
			"creater":      msg.Creater,
			"creater_name": "",
		},
	})
}

// MsgTransfer 转账消息
type MsgTransfer struct {
	RecordNo string `json:"record_no"` // 转账记录编号
	Remark   string `json:"remark"`    // 转账备注
	Amount   int64  `json:"amount"`    // 转账金额
	FromUID  string `json:"from_uid"`  // 转账者uid
	ToUID    string `json:"to_uid"`    // 接受者uid
	Status   string `json:"status"`    // 转账状态 noaccept: 未接收 accepted：已收款
}

// MsgTransferRecover 转账回收
type MsgTransferRecover struct {
	Creater   string `json:"uid"`        // 转账创建者
	Receiver  string `json:"receiver"`   // 接收者
	Amount    int64  `json:"amount"`     // 退款金额
	ExpiredAt int64  `json:"expired_at"` // 过期时间
}
