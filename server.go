package main

import (
	// "encoding/json"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/codegangsta/martini-contrib/render"
	handle "github.com/jiangzhx/receive/handle"
	model "github.com/jiangzhx/receive/model"
	"runtime"
	// redis "github.com/jiangzhx/receive/redis"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	logger, _ := log.LoggerFromConfigAsFile("seelog.xml")
	log.ReplaceLogger(logger)
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Get("/gettime", gettime)
	m.Post("/:what", binding.Json(model.Model{}), event)
	m.Run()
}

func gettime(r render.Render) {
	t := time.Now().Local()
	ts, _ := strconv.Atoi(t.Format("20060102150405"))
	r.JSON(200, map[string]interface{}{"status": 0, "ts": ts})
}

func event(req *http.Request, params martini.Params, model model.Model, r render.Render) {
	what := params["what"]
	ip := getIP(req)
	if !strings.EqualFold(what, "event") {
		model.What = what
	}

	context := model.Context
	tmp := make(map[string]string)
	for key, value := range context {
		tmp[strings.ToLower(key)] = value
	}
	model.Context = tmp
	model.Where = what
	model.Context["ip"] = ip

	go handle.EasyRealTimeSendTask(model)

	// b, _ := json.Marshal(model)
	// log.Debug(string(b))

	r.JSON(200, map[string]interface{}{"status": 0})
}

func formatData(model model.Model) string {
	what, topic := model.What, model.What

	if strings.EqualFold("loggedin", strings.TrimSpace(what)) {
		topic = "dau"
	} else if strings.EqualFold("register", strings.TrimSpace(what)) {
		what, topic = "reged", "reged"
	} else if strings.EqualFold("heartbeat", strings.TrimSpace(what)) {
		what, topic = "hb", "hb"
	}

	// data.append(model.getWho())
	// 	.append("\t")
	// 	.append(model.getWhen())
	// 	.append("\t")
	// 	.append(model.getWhere())
	// 	.append("\t")
	// 	.append(what)
	// 	.append("\t")
	// 	.append(this.formatContext(model, ip))
	// 	.append("\t")
	// 	.append(model.getAppid())
	// 	.append("\t")
	// 	.append(model.getDs());
	log.Debug(topic)
	return fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t", model.Who, model.When, model.Where, model.What, "", model.Appid, model.When[:strings.Index(model.When, " ")])
}

func getIP(req *http.Request) string {
	ip := req.Header.Get("X-Forwarded-For")
	if ip == "" || strings.EqualFold(ip, "unknown") || len(ip) == 0 {
		ip = req.Header.Get("Proxy-Client-IP")
	}
	if ip == "" || strings.EqualFold(ip, "unknown") || len(ip) == 0 {
		ip = req.Header.Get("WL-Proxy-Client-IP")
	}
	if ip == "" || strings.EqualFold(ip, "unknown") || len(ip) == 0 {
		ip = req.Header.Get("HTTP_CLIENT_IP")
	}
	if ip == "" || strings.EqualFold(ip, "unknown") || len(ip) == 0 {
		ip = req.Header.Get("HTTP_X_FORWARDED_FOR")
	}

	return ip
}
