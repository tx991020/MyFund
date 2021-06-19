package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/astaxie/beego/logs"
	resty "github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"github.com/tx991020/utils"
	"github.com/tx991020/utils/fx"
)

var (
	httpClient *resty.Request
	fundDetail = `https://fundmobapi.eastmoney.com/FundMNewApi/FundMNContrastPKList?deviceid=app_h5_zixuanduibi&version=5.8.5&product=efund&plat=Iphone&FCODES=%s&callback=jsonp_1613115116483_50119%27`
	fundStock  = `https://fund.eastmoney.com/API/FundDXGJJ.ashx?callback=jQuery183001088009341492624_1612592150749&r=1612596186000&m=8&pageindex=%d&sorttype=desc&SFName=RATIO&IsSale=1&_=1612596186706`
	fundList   = `http://fund.eastmoney.com/API/FundDXGJJ.ashx?m=8&pageindex=1&pagesize=6000&sorttype=desc&SFName=SUMPLACE&IsSale=0&_=1611753827717`
	zuhe       = `https://danjuanfunds.com/djapi/plan/last_detail?plan_code=%s`
	plans      = `https://danjuanfunds.com/djapi/fundx/portfolio/v3/plan/united/page?tab=4&page=1&size=20&default_order=3&invest_strategy=&type=&manager_type=&yield_between=&mz_between=`
	planDetail = `https://danjuanfunds.com/djapi/plan/%s`
	pe         = `https://danjuanfunds.com/djapi/index_eva/dj`
	dayRank    = `curl -H 'gtoken: ceaf-99b78df4fab899ced2cc7e10ce64bccb' -H 'clientInfo: ttjj-JSN-AL00a-Android-10' -H 'Host: fundmobapi.eastmoney.com' -H 'User-Agent: okhttp/3.12.0' --data "appType=ttjj&Sort=desc&product=EFund&gToken=ceaf-99b78df4fab899ced2cc7e10ce64bccb&version=6.4.0&DataConstraintType=0&ctoken=ca-ueu1cjdeue6-dredkjf18uuqaed18.9&ESTABDATE=&deviceid=884a04b668a18d4fb77dec2b17efe2e4%7C%7C346719440025268&ENDNAV=&FundType=0&BUY=true&pageIndex=AAAAA&RLEVEL_SZ=&RISKLEVEL=&DISCOUNT=&utoken=rhd8rqkr6cq1r--6h8eqe6eed1edha8eccfac8b8.9&CompanyId=&SortColumn=RZDF&pageSize=30&MobileKey=884a04b668a18d4fb77dec2b17efe2e4%7C%7C346719440025268&TOPICAL=&plat=Android&ISABNORMAL=true" --compressed 'https://fundmobapi.eastmoney.com/FundMNewApi/FundMNRank'`
)

func init() {
	httpClient = resty.New().R().SetHeaders(map[string]string{"User-Agent": "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.1 (KHTML, like Gecko) Chrome/22.0.1207.1 Safari/537.1"})

}

type FCode struct {
	Id string `json:"_id"`
}

type Shui struct {
	Time string `json:"time"`
	Rate string `json:"rate"`
}

type Fund struct {
	Id        string  `json:"_id"`
	SHORTNAME string  `json:"SHORTNAME"`
	STKNUM    int     `json:"STKNUM"`
	RATIO     float64 `json:"RATIO"`
}

type Stock struct {
	Scode string  `json:"Scode"`
	Sname string  `json:"Sname"`
	Ratio float64 `json:"Ratio"`
}

type Trading struct {
	Code    string  `json:"fd_code"`
	Name    string  `json:"fd_name"`
	Portion float64 `json:"portion"`
}

func ReqFundList(url string) (arr []string) {
	get, err := resty.New().SetRetryCount(3).R().Get(url)
	if err != nil {
		logs.Error(err)
		return arr
	}
	array := gjson.Get(get.String(), "Datas").String()
	result := make([]*FCode, 0)
	err = json.Unmarshal([]byte(array), &result)
	if err != nil {
		logs.Error(err)
		return
	}
	for _, fund := range result {
		arr = append(arr, fund.Id)
	}
	return
}
func ReqFundBasicInfo(url string) string {

	res, err := httpClient.Get(url)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	data := strings.TrimLeft(res.String(), "jsonp_1613115116483_50119'(")
	data2 := strings.TrimRight(data, ")")

	detail := gjson.Get(data2, "Datas").String()
	no := gjson.Get(detail, "FCODE.0").String()
	rankw := gjson.Get(detail, "RANKW.0").String()
	rankm := gjson.Get(detail, "RANKM.0").String()
	name := gjson.Get(detail, "SHORTNAME.0").String()
	week := gjson.Get(detail, "SYL_Z.0").String()
	one_m := gjson.Get(detail, "SYL_Y.0").String()
	three_m := gjson.Get(detail, "SYL_3Y.0").String()
	six_m := gjson.Get(detail, "SYL_6Y.0").String()
	one_year := gjson.Get(detail, "SYL_1N.0").String()
	bodong := gjson.Get(detail, "STDDEV_1N.0").String()
	shap_1n := gjson.Get(detail, "SHARP_1N.0").String()
	MAXRETRA_1N := gjson.Get(detail, "MAXRETRA_1N.0").String()

	x, _ := utils.TypeConversion(one_year, "float64")
	y, _ := utils.TypeConversion(MAXRETRA_1N, "float64")
	huiche := func(a, b float64) float64 {
		if b == 0 {
			return 0
		}
		return a / b
	}(x.Float(), y.Float())
	RATE := gjson.Get(detail, "RATE.0").String()
	rate12 := gjson.Get(detail, "RATEINFO.0.sh").String()
	shuilv := make([]*Shui, 0)
	err = json.Unmarshal([]byte(rate12), &shuilv)
	if err != nil {
		fmt.Println(err)
	}
	shs := make([]string, 0)
	for _, shui := range shuilv {
		sh := fmt.Sprintf("%s:%s", shui.Time, shui.Rate)
		shs = append(shs, sh)
	}

	rate13 := strings.Join(shs, "|")
	rate14 := strings.ReplaceAll(rate13, ",", "")
	rate15 := strings.ReplaceAll(rate14, "\r", "")
	rate1 := strings.ReplaceAll(rate15, "\n", "")

	NEWSDCC := gjson.Get(detail, "NEWSDCC.0").String()
	result := make([]*Stock, 0)
	err = json.Unmarshal([]byte(NEWSDCC), &result)
	if err != nil {
		return ""
	}
	arr1 := make([]string, 0)
	arr := make([]string, 0)
	for _, stock := range result {
		//s1 := fmt.Sprintf("%s-%s-%.2f", stock.Scode, stock.Sname, stock.Ratio)
		s1 := stock.Sname
		s := stock.Scode

		arr = append(arr, s)
		arr1 = append(arr1, s1)

	}
	stocks1 := strings.Join(arr1, ",")
	stocks := strings.Join(arr, ",")
	info := fmt.Sprintf(`%s,%s,%s,%.2f,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,"{%s}","{%s}"`, no, rankw, rankm, huiche, name, week, one_m, three_m, six_m, one_year, bodong, shap_1n, MAXRETRA_1N, RATE, rate1, stocks, stocks1)
	return info

}

func ReqFundStocks(url string, path string) {

	fx.From(func(source chan<- interface{}) {
		for i := 1; i < 221; i++ {

			source <- fmt.Sprintf(url, i)
		}
	}).Map(func(item interface{}) interface{} {
		info := ""

		res, err := httpClient.Get(item.(string))

		if err != nil {
			fmt.Println(err)

		} else {
			data := strings.TrimLeft(res.String(), "jQuery183001088009341492624_1612592150749(")
			info = gjson.Get(strings.TrimRight(data, ")"), "Datas").String()

		}

		return info
	}).ForEach(func(item interface{}) {

		result := make([]*Fund, 0)
		err := json.Unmarshal([]byte(item.(string)), &result)
		if err != nil {
			fmt.Println(err)

		} else {
			for _, fund := range result {

				info := fmt.Sprintf("%s,%s,%d,%.2f", fund.Id, fund.SHORTNAME, fund.STKNUM, fund.RATIO)

				utils.PutContentsAppend(path, info)
				utils.PutContentsAppend(path, "\n")

			}
		}
	})
}

func ReadCSVFund(url, tmp, to string, list []string) (err error) {

	fx.From(func(source chan<- interface{}) {
		for _, i2 := range list {
			source <- fmt.Sprintf(url, i2)
		}
	}).Map(func(item interface{}) interface{} {
		info := ReqFundBasicInfo(item.(string))

		return info
	}, fx.WithWorkers(8)).ForEach(func(item interface{}) {
		utils.PutContentsAppend(tmp, item.(string))
		utils.PutContentsAppend(tmp, "\n")

		//统计排名和重仓
	})
	sp := fmt.Sprintf(`grep -v '^$' %s > %s`, tmp, to)
	fmt.Println(sp)
	cmd := exec.Command("/bin/sh", "-c", sp)
	_, err = cmd.Output()

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
		os.Remove(tmp)
	}()

	return
}

func GetZuHeLatest(fcode, url string, to string) string {
	fmt.Println(url)
	res, err := httpClient.Get(fmt.Sprintf(url, fcode))
	if err != nil {
		fmt.Println(err)
		return ""
	}

	trads := make([]*Trading, 0)
	plname := gjson.Get(res.String(), "data.plan_name").String()
	data := gjson.Get(res.String(), "data.info.trading_elements").String()
	fmt.Println(data)
	err = json.Unmarshal([]byte(data), &trads)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	arr := make([]string, 0)
	arr1 := make([]string, 0)
	arr2 := make([]string, 0)
	for _, v := range trads {
		s := v.Code
		s1 := v.Name
		s2 := v.Portion
		arr = append(arr, s)
		arr1 = append(arr1, s1)
		arr2 = append(arr2, fmt.Sprintf("%0.2f", s2))

	}
	codes := strings.Join(arr, ",")
	names := strings.Join(arr1, ",")
	rates := strings.Join(arr2, ",")
	item := fmt.Sprintf(`%s,%s,"{%s}","{%s}","{%s}"`, fcode, plname, codes, names, rates)
	return item

}

func GetPlans(url string) (zuhe []string) {
	res, err := httpClient.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	array := gjson.Get(res.String(), "data.items").Array()
	for _, v := range array {
		s := gjson.Get(v.String(), "plan_code").String()
		zuhe = append(zuhe, s)
	}
	return

}

func GetPlanDetail(url string) string {
	res, err := httpClient.Get(url)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	week := gjson.Get(res.String(), "data.plan_derived.nav_grl1w").String()
	one_m := gjson.Get(res.String(), "data.plan_derived.nav_grl1m").String()
	three_m := gjson.Get(res.String(), "data.plan_derived.nav_grl3m").String()
	six_m := gjson.Get(res.String(), "data.plan_derived.nav_grl6m").String()
	one_year := gjson.Get(res.String(), "data.plan_derived.nav_grl1y").String()
	return fmt.Sprintf("%s,%s,%s,%s,%s", week, one_m, three_m, six_m, one_year)
}

func DanJuan(plans, zuhe, to string) {

	arr := GetPlans(plans)
	fx.From(func(source chan<- interface{}) {
		for _, i2 := range arr {
			source <- i2
		}
	}).Map(func(item interface{}) interface{} {
		stock := GetZuHeLatest(item.(string), zuhe, to)
		shouyi := GetPlanDetail(fmt.Sprintf(planDetail, item.(string)))
		return fmt.Sprintf("%s,%s", stock, shouyi)
	}, fx.WithWorkers(1)).ForEach(func(item interface{}) {
		utils.PutContentsAppend(to, item.(string))
		utils.PutContentsAppend(to, "\n")
	})

}

func GetPE(url, to string) {
	res, err := httpClient.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	array := gjson.Get(res.String(), "data.items").Array()
	for _, v := range array {
		fmt.Println(v.String())
		id := gjson.Get(v.String(), "id").String()

		name := gjson.Get(v.String(), "name").String()
		p := gjson.Get(v.String(), "pb_percentile").Float()
		link := gjson.Get(v.String(), "url").String()
		matchString, err := utils.MatchString(`funding\/(\d+)`, link)
		if err != nil {
			fmt.Println(err)
			return
		}

		fcode := ""
		if len(matchString) == 2 {
			fcode = matchString[1]
		}
		item := fmt.Sprintf(`%s,%s,%0.3f,%s`, id, name, p, fcode)
		utils.PutContentsAppend(to, item)
		utils.PutContentsAppend(to, "\n")
	}
}

func PostDayRank(url, to string) {
	for i := 1; i < 230; i++ {
		replace := strings.Replace(url, "AAAAA", strconv.Itoa(i), -1)
		fmt.Println(replace)
		cmd := exec.Command("/bin/sh", "-c", replace)
		data, err := cmd.Output()
		if err != nil {
			fmt.Println(err)
			continue
		}
		array := gjson.Get(string(data), "Datas").Array()

		for _, result := range array {
			fcode := gjson.Get(result.String(), "FCODE").String()
			name := gjson.Get(result.String(), "SHORTNAME").String()
			RQ := gjson.Get(result.String(), "FSRQ").String()
			RZ := gjson.Get(result.String(), "RZDF").String()
			SYL_Z := gjson.Get(result.String(), "SYL_Z").String()
			SYL_Y := gjson.Get(result.String(), "SYL_Y").String()
			SYL_3Y := gjson.Get(result.String(), "SYL_3Y").String()
			SYL_6Y := gjson.Get(result.String(), "SYL_6Y").String()

			item := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s", fcode, name, RQ, RZ, SYL_Z, SYL_Y, SYL_3Y, SYL_6Y)
			utils.PutContentsAppend(to, item)
			utils.PutContentsAppend(to, "\n")

		}
	}

}

////postgreql 连接，数据导入
//func main() {
//
//	//ReqFundStocks(fundStock, "/Users/andy/GoLang/src/mgo-boot/a70/stock_rate.csv")

//	//ReqFundBasicInfo("https://fundmobapi.eastmoney.com/FundMNewApi/FundMNContrastPKList?deviceid=app_h5_zixuanduibi&version=5.8.5&product=efund&plat=Iphone&FCODES=004041&callback=jsonp_1613115116483_50119%27")
//	DanJuan(plans, zuhe, "/Users/andy/GoLang/src/mgo-boot/a70/danjuan.csv")
//	//
//	//GetPE(pe,"/Users/andy/GoLang/src/mgo-boot/a70/pe.csv")
//
//	//PostDayRank(dayRank, "/Users/andy/GoLang/src/mgo-boot/a70/paihang.csv")
//
//
//
//
//}
