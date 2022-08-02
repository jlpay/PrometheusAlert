package controllers

import (
	"PrometheusAlert/models"
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func PostToBkalert(text, BkUrl, logsign string) string {
	logs.Info(logsign, "[Bkalert]", "接收的消息：", text)
	// 去除前后空格
	//text := strings.TrimSpace(string(text))

	// 去除换行符
	//text = strings.Replace(tmpText, "\n", "", -1)

	//logs.Info(logsign, "[Bkalert]", "处理后的消息：", tmpText)
	open := beego.AppConfig.String("bkalert")
	if open != "1" {
		logs.Info(logsign, "[bkalert]", "蓝鲸告警接口未配置未开启状态,请先配置bkalert为1")
		return "蓝鲸告警接口未配置未开启状态,请先配置bkalert为1"
	}

	// 头 -H 'X-Secret:
	secret := beego.AppConfig.String("Bkalert_secret")

	JsonMsg := bytes.NewReader([]byte(text))
	var tr *http.Transport
	if proxyUrl := beego.AppConfig.String("proxy"); proxyUrl != "" {
		proxy := func(_ *http.Request) (*url.URL, error) {
			return url.Parse(proxyUrl)
		}
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           proxy,
		}
	} else {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("POST", BkUrl, JsonMsg)
	if err != nil {
		logs.Error(logsign, "[Webhook]", err.Error())
	}
	//logs.Info(logsign, "[BkUrl]", "添加Header: X-Secret: ", secret)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Secret", secret)

	res, err := client.Do(req)
	//defer res.Body.Close()
	//result, err := ioutil.ReadAll(res.Body)

	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logs.Error(logsign, "[BkUrl]", err.Error())
	}
	defer res.Body.Close()
	models.AlertToCounter.WithLabelValues("BkUrl").Add(1)
	ChartsJson.Bkalert += 1
	logs.Info(logsign, "[BkUrl]", string(result))
	return string(result)
}
