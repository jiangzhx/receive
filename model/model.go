package model

import (
	"fmt"
	"github.com/jiangzhx/receive/util/date"
)

type Model struct {
	Appid   string            `json:"appid"`
	When    string            `json:"when"`
	Who     string            `json:"who"`
	What    string            `json:"what"`
	Where   string            `json:"where"`
	Context map[string]string `json:"context"`
}

func (model *Model) GetChannelsKey() string {
	return fmt.Sprintf("CHANNELS-ALL:%s", model.Appid)
}

func (model *Model) GetInstallChannelALLFactKey() string {
	return fmt.Sprintf("INSTALL-CHANNEL-ALL-FACT:%s:%s", model.Appid, model.Context["channelid"])
}
func (model *Model) GetInstallChannelALLFactKeyWithChannel(channel string) string {
	return fmt.Sprintf("INSTALL-CHANNEL-ALL-FACT:%s:%s", model.Appid, channel)
}

func (model *Model) GetInstallALLFactKey() string {
	return fmt.Sprintf("INSTALL-ALL-FACT:%s", model.Appid)
}

func (model *Model) GetRegisterAllFactKey() string {
	return fmt.Sprintf("REGISTER-ALL-FACT:%s", model.Appid)
}

func (model *Model) GetRegisterAllChannelFactKey() string {
	return fmt.Sprintf("REGISTER-CHANNEL-ALL-FACT:%s:%s", model.Appid, model.Context["channelid"])
}

func (model *Model) GetRegisterDayFactKey() string {
	return fmt.Sprintf("REGISTER-DAY-FACT:%s:%s", model.Appid,
		date.Format(model.When, date.DayFormat))
}

func (model *Model) GetRegisterDayFactKeyWithRet(ret int) string {
	d := date.PreDay(model.When, ret).Format(date.TimeFormat)

	return fmt.Sprintf("REGISTER-DAY-FACT:%s:%s", model.Appid,
		date.Format(d, date.DayFormat))
}

func (model *Model) GetRegisterChannelDayFactKey() string {
	return fmt.Sprintf("REGISTER-CHANNEL-DAY-FACT:%s:%s:%s", model.Appid,
		model.Context["channelid"], date.Format(model.When, date.DayFormat))
}

func (model *Model) GetRegisterChannelDayFactKeyWithRet(ret int) string {
	d := date.PreDay(model.When, ret).Format(date.TimeFormat)

	return fmt.Sprintf("REGISTER-CHANNEL-DAY-FACT:%s:%s:%s", model.Appid, model.Context["channelid"], date.Format(d, date.DayFormat))
}

func (model *Model) GetRegisterServerDayFactKey() string {
	return fmt.Sprintf("REGISTER-SERVER-DAY-FACT:%s:%s:%s", model.Appid,
		model.Context["serverid"], date.Format(model.When, date.DayFormat))
}

func (model *Model) GetRegisterServerDayFactKeyWithRet(ret int) string {
	d := date.PreDay(model.When, ret).Format(date.TimeFormat)

	return fmt.Sprintf("REGISTER-SERVER-DAY-FACT:%s:%s:%s", model.Appid, model.Context["serverid"], date.Format(d, date.DayFormat))
}

func (model *Model) GetRegisterChannelAllFactKey(channel string) string {
	return fmt.Sprintf("REGISTER-CHANNEL-ALL-FACT:%s:%s", model.Appid, channel)
}

func (model *Model) GetPayerAllFactKey() string {
	return fmt.Sprintf("PAYER-ALL-FACT:%s", model.Appid)
}

func (model *Model) GetPayerServerAllFactKey() string {
	return fmt.Sprintf("PAYER-SERVER-ALL-FACT:%s", model.Appid, model.Context["serverid"])
}
