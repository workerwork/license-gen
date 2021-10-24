package server

import (
	"fmt"
	"license-gen/conf"
	"license-gen/utils"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup
var data = Data{}

// 写log文件
func Log(data Data) {
	timeLayoutStr := "2006-01-02 15:04:05"
	timeStr := time.Now().Format(timeLayoutStr)
	formatStr := strings.Repeat("*", 20)
	str := fmt.Sprintf("%s%s%s\n%s\n", formatStr, timeStr, formatStr, data.String())
	utils.WriteFile(conf.ServerConf.Log, str)
}

func Serve() {
	for {
		go run()
		time.Sleep(conf.ServerConf.Timer * time.Second)
	}
}

func run() {
	data, err := ClientGetInfo()
	if err != nil {
		return
	}
	for _, item := range data.Item_list {
		wg.Add(1)
		go func() {
			lic, err := NewLic()
			if err != nil {
				wg.Done()
				return
			}
			lic.XmlNs(data.Product).
				Code(data.Product).
				ProductVersion(item.Version).
				AuthCode(item.Auth_code).
				CreateTime().
				TotalTime(item.Total_time).
				MaxUeNum(item.Max_ue_num).
				MaxEnbNum(item.Max_enb_num).
				PathOaId(data.Oa_id).
				PathAuthCode(item.Auth_code)
			if err := lic.ToXML(); err != nil {
				wg.Done()
				return
			}
			if err := lic.GenLic(); err != nil {
				wg.Done()
				return
			}
			wg.Done()
		}()
		wg.Wait()
		//TODO
		Log(data)
	}
	advise_result := AdviseResult{
		Oa_id:         data.Oa_id,
		Apply_type:    data.Apply_type,
		Create_result: true,
		File_name:     "test",
		Msg:           "success",
	}
	ClientPostAdviseResult(advise_result)
}
