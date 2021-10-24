package server

import (
	"fmt"
	"license-gen/conf"
	"sync"
	"time"
)

var wg sync.WaitGroup
var data = Data{}

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
	} else {
		data.Log()
	}
	// response msg
	advise_result := AdviseResult{
		Oa_id:         data.Oa_id,
		Apply_type:    data.Apply_type,
		Create_result: false,
		File_name:     "",
		Msg:           "",
	}
	//result := make(map[string]string)
	//slice := make([]string, len(data.Item_list))
	for _, item := range data.Item_list {
		wg.Add(1)
		go func() {
			lic, err := NewLic()
			if err != nil {
				advise_result.Msg += fmt.Sprintf("%s: License gen failed!\n", item.Auth_code)
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
				advise_result.Msg += fmt.Sprintf("%s:License gen failed! ", item.Auth_code)
				wg.Done()
				return
			}
			if err := lic.GenLic(); err != nil {
				advise_result.Msg += fmt.Sprintf("%s:License gen failed! ", item.Auth_code)
				wg.Done()
				return
			}
			advise_result.Create_result = true
			advise_result.File_name = data.Oa_id + ".zip"
			wg.Done()
		}()
		wg.Wait()
		if advise_result.Msg == "" {
			advise_result.Msg = "success"
		}
	}
	advise_result.Log()
	ClientPostAdviseResult(advise_result)
}
