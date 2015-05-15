package handle

import (
	"fmt"
	// log "github.com/cihub/seelog"
	redisgo "github.com/garyburd/redigo/redis"
	model "github.com/jiangzhx/receive/model"
	"github.com/jiangzhx/receive/redis"
	"github.com/jiangzhx/receive/util/date"
	// "math"
	"strconv"
	"strings"
)

// var (
// 	retentions []int = []int{-1, -3, -7, -14, -30}
// )

func FullRealTimeSendTask() {

}

func EasyRealTimeSendTask(model model.Model) {
	channelHandle(model)
	serverHandle(model)
	switch {
	case strings.EqualFold("install", model.What):
		installHandle(model)
	case strings.EqualFold("register", model.What):
		registerHandle(model)
		dauHandle(model)
	case strings.EqualFold("loggedin", model.What):
		dauHandle(model)
	case strings.EqualFold("payment", model.What):
		paymentHandle(model)
		dauHandle(model)
	case strings.EqualFold("heartbeat", model.What):
		heartbeatHandle(model)
	}
}
func channelHandle(model model.Model) {
	conn := redis.GetConn()
	defer conn.Close()
	conn.Do("sadd", model.GetChannelsKey(), model.Context["channelid"])
}

func serverHandle(model model.Model) {
	conn := redis.GetConn()
	defer conn.Close()
	conn.Do("sadd", "SERVERS-ALL:"+model.Appid, model.Context["serverid"])

}

func installHandle(model model.Model) {
	if redis.Bexist(model.GetInstallALLFactKey(), model.Context["deviceid"]) != true {
		return
	}
	conn := redis.GetConn()
	defer conn.Close()
	// service for count
	conn.Do("pfadd", fmt.Sprintf("INSTALL-HOUR:%s:%s", model.Appid, date.Format(model.When, date.HourFormat)), model.Context["deviceid"])

	// service for server
	conn.Do("pfadd", fmt.Sprintf("INSTALL-SERVER-HOUR:%s:%s:%s", model.Appid, model.Context["serverid"], date.Format(model.When, date.HourFormat)), model.Context["deviceid"])

	// service for channel
	conn.Do("pfadd", fmt.Sprintf("INSTALL-CHANNEL-HOUR:%s:%s:%s", model.Appid, model.Context["channelid"], date.Format(model.When, date.HourFormat)), model.Context["deviceid"])

	redis.Badd(model.GetInstallALLFactKey(), model.Context["deviceid"])
	redis.Badd(model.GetInstallChannelALLFactKey(), model.Context["deviceid"])
}

func registerHandle(model model.Model) {
	if redis.Bexist(model.GetRegisterAllFactKey(), model.Who) != true {
		return
	}
	// service for retention
	redis.Badd(model.GetRegisterAllFactKey(), model.Who)
	redis.Badd(model.GetRegisterChannelDayFactKey(), model.Who)
	redis.Badd(model.GetRegisterServerDayFactKey(), model.Who)

	conn := redis.GetConn()
	defer conn.Close()
	// service for count
	conn.Do("pfadd", fmt.Sprintf("REGISTER-HOUR:%s:%s", model.Appid, date.Format(model.When, date.HourFormat)), model.Who)
	// service for server
	conn.Do("pfadd", fmt.Sprintf("REGISTER-SERVER-HOUR:%s:%s:%s", model.Appid, model.Context["serverid"], date.Format(model.When, date.HourFormat)), model.Who)
	// service for channel
	conn.Do("pfadd", fmt.Sprintf("REGISTER-CHANNEL-HOUR:%s:%s:%s", model.Appid, model.Context["channelid"], date.Format(model.When, date.HourFormat)), model.Who)
	// service for all fact data
	redis.Badd(model.GetRegisterAllFactKey(), model.Who)
	// service for register by channel fact data
	redis.Badd(model.GetRegisterAllChannelFactKey(), model.Who)
}

func dauHandle(model model.Model) {
	model.Context["channelid"] = GetFirstTimeChannelByWho(model)

	conn := redis.GetConn()
	defer conn.Close()
	// server for hour count
	conn.Do("pfadd", fmt.Sprintf("DAU-HOUR:%s:%s", model.Appid, date.Format(model.When, date.HourFormat)), model.Who)

	// server for server
	conn.Do("pfadd", fmt.Sprintf("DAU-SERVER-HOUR:%s:%s:%s", model.Appid, model.Context["serverid"], date.Format(model.When, date.HourFormat)), model.Who)

	// server for channel
	conn.Do("pfadd", fmt.Sprintf("DAU-CHANNEL-HOUR:%s:%s:%s", model.Appid, model.Context["channelid"], date.Format(model.When, date.HourFormat)), model.Who)
}

func paymentHandle(model model.Model) {
	model.Context["channelid"] = GetFirstTimeChannelByWho(model)

	conn := redis.GetConn()
	defer conn.Close()
	paymenttype := model.Context["paymenttype"]
	amount := float64(0)
	val, _ := model.Context["currencyamount"]
	amount, _ = strconv.ParseFloat(val, 64)

	if strings.EqualFold("free", paymenttype) && amount > 0 {
		// ---------------calc for payment amount-----------------
		conn.Do("incrByFloat", fmt.Sprintf("PAYMENT-HOUR:%s:%s", model.Appid, date.Format(model.When, date.HourFormat)), amount)

		conn.Do("incrByFloat", fmt.Sprintf("PAYMENT-SERVER-HOUR:%s:%s", model.Appid, model.Context["serverid"], date.Format(model.When, date.HourFormat)), amount)

		conn.Do("incrByFloat", fmt.Sprintf("PAYMENT-CHANNEL-HOUR:%s:%s", model.Appid, model.Context["channelid"], date.Format(model.When, date.HourFormat)), amount)
		// ---------------calc for payers-----------------
		conn.Do("pfadd", fmt.Sprintf("PAYER-DAY:%s:%s", model.Appid, date.Format(model.When, date.DayFormat)), model.Who)

		conn.Do("pfadd", fmt.Sprintf("PAYER-HOUR:%s:%s", model.Appid, date.Format(model.When, date.HourFormat)), model.Who)

		conn.Do("pfadd", fmt.Sprintf("PAYER-SERVER-HOUR:%s:%s", model.Appid, model.Context["serverid"], date.Format(model.When, date.HourFormat)), model.Who)

		conn.Do("pfadd", fmt.Sprintf("PAYER-CHANNEL-HOUR:%s:%s", model.Appid, model.Context["channelid"], date.Format(model.When, date.HourFormat)), model.Who)

		// ------calc first time payment and payers-----------------------
		if redis.Bexist(model.GetPayerAllFactKey(), model.Who) == false {
			conn.Do("pfadd", fmt.Sprintf("PAYER-FT-DAY:%s:%s", model.Appid, date.Format(model.When, date.DayFormat)), model.Who)
			conn.Do("incrByFloat", fmt.Sprintf("PAYMENT-FT-DAY:%s:%s", model.Appid, date.Format(model.When, date.DayFormat)), amount)
			redis.Badd(model.GetPayerAllFactKey(), model.Who)
		}

		if redis.Bexist(model.GetPayerServerAllFactKey(), model.Who) == false {
			conn.Do("pfadd", fmt.Sprintf("PAYER-FT-SERVER-DAY:%s:%s:%s", model.Appid, model.Context["serverid"], date.Format(model.When, date.DayFormat)), model.Who)
			conn.Do("incrByFloat", fmt.Sprintf("PAYMENT-FT-DAY:%s:%s", model.Appid, model.Context["serverid"], date.Format(model.When, date.DayFormat)), amount)
			redis.Badd(model.GetPayerServerAllFactKey(), model.Who)
		}
	}
}

func heartbeatHandle(model model.Model) {
	// when := "2015-04-23 18:04:05"
	cm := date.Parse(model.When).Minute()
	cm = cm / 5 * 5
	minute := "00"
	if cm < 10 {
		minute = fmt.Sprintf("0%d", cm)
	} else {
		minute = fmt.Sprintf("%d", cm)
	}
	minute = fmt.Sprintf("%s:%sZ", date.Format(model.When, "2006-01-02T15"), minute)
	conn := redis.GetConn()
	defer conn.Close()
	conn.Do("pfadd", fmt.Sprintf("HEARTBEAT-5M:%s:%s", model.Appid, minute), model.Who)
	conn.Do("pfadd", fmt.Sprintf("HEARTBEAT-SERVER-5M:%s:%s:%s", model.Appid, model.Context["serverid"], minute), model.Who)
}

func retentionHandle(model model.Model) {
	model.Context["channelid"] = GetFirstTimeChannelByWho(model)
	retentions := []int{-1, -3, -7, -14, -30}

	conn := redis.GetConn()
	defer conn.Close()

	for i := 0; i < len(retentions); i++ {
		ret := retentions[i]
		when := date.PreDay(model.When, ret).Format(date.DayFormat)
		if redis.Bexist(model.GetRegisterDayFactKeyWithRet(ret), model.Who) {
			conn.Do("pfadd", fmt.Sprintf("RETENTION-DAY%d:%s:%s", CalcAbs(ret), model.Appid, when), model.Who)
		}
		if redis.Bexist(model.GetRegisterChannelDayFactKeyWithRet(ret), model.Who) {
			conn.Do("pfadd", fmt.Sprintf("RETENTION-CHANNEL-DAY%d:%s:%s:%s", CalcAbs(ret), model.Appid, model.Context["channelid"], when), model.Who)

		}
		if redis.Bexist(model.GetRegisterServerDayFactKeyWithRet(ret), model.Who) {
			conn.Do("pfadd", fmt.Sprintf("RETENTION-SERVER-DAY%d:%s:%s:%s", CalcAbs(ret), model.Appid, model.Context["serverid"], when), model.Who)
		}
	}
}

func roiHandle(model model.Model) {
	model.Context["channelid"] = GetFirstTimeChannelByWho(model)
	// retentions := []int{-1, -3, -7, -14, -30}
	retentions := make([]int, 30)
	for i := 0; i < len(retentions); i++ {
		retentions[i] = 0 - i
	}
	conn := redis.GetConn()
	defer conn.Close()

	paymenttype := model.Context["paymenttype"]
	amount := float64(0)
	val, _ := model.Context["currencyamount"]
	amount, _ = strconv.ParseFloat(val, 64)

	if strings.EqualFold("free", paymenttype) && amount > 0 {
		for i := 0; i < len(retentions); i++ {
			ret := retentions[i]
			when := date.PreDay(model.When, ret).Format(date.DayFormat)
			if redis.Bexist(model.GetRegisterDayFactKeyWithRet(ret), model.Who) {
				conn.Do("pfadd", fmt.Sprintf("ROIER-DAY%d:%s:%s", CalcAbs(ret), model.Appid, when), model.Who)
				conn.Do("incrByFloat", fmt.Sprintf("ROI-DAY%d:%s:%s", CalcAbs(ret), model.Appid, when), amount)

			}
			if redis.Bexist(model.GetRegisterChannelDayFactKeyWithRet(ret), model.Who) {
				conn.Do("pfadd", fmt.Sprintf("ROIER-CHANNEL-DAY%d:%s:%s:%s", CalcAbs(ret), model.Appid, model.Context["channelid"], when), model.Who)
				conn.Do("incrByFloat", fmt.Sprintf("ROI-CHANNEL-DAY%d:%s:%s:%s", CalcAbs(ret), model.Appid, model.Context["channelid"], when), amount)

			}
		}
	}
}

func GetFirstTimeChannelByWho(model model.Model) string {
	conn := redis.GetConn()
	defer conn.Close()

	keys, _ := redisgo.Strings(conn.Do("smembers", model.GetChannelsKey()))

	for i := 0; i < len(keys); i++ {
		if redis.Bexist(model.GetRegisterChannelAllFactKey(keys[i]), model.Who) {
			return keys[i]
		}
	}
	return "UNKNOWN"
}

func CalcAbs(a int) (ret int) {
	ret = (a ^ a>>31) - a>>31
	return
}
