package server

import (
	"fmt"
	"license-gen/conf"
	"sync"
	"time"
)

var (
	wg   sync.WaitGroup
	data = Data{}
)

const (
	FORMAT  = ".zip"
	SUCCESS = "success"
)

// @function Serve
// @description 服务入口循环,协程runtime
// @param ""
// @return ""
func Serve() {
	for {
		go run()
		time.Sleep(conf.ServerConf.Timer * time.Second)
	}
}

// @function run
// @description 1)获取data 2)生成lic 3)上传文件
// @param ""
// @return ""
func run() {
	//从LicenseCenter获取源数据
	data, err := ClientGetInfo()
	if err != nil {
		return
	} else {
		//记录获取的data数据
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
		go func() { //每个License鉴权码使用一个单独的协程
			lic, err := NewLic()
			if err != nil { //没有细分err，统一回复server: License gen failed!
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
			if err := lic.ToXML(); err != nil { //没有细分err，统一回复server: License gen failed!
				advise_result.Msg += fmt.Sprintf("%s:License gen failed! ", item.Auth_code)
				wg.Done()
				return
			}
			if err := lic.GenLic(); err != nil { //没有细分err，统一回复server: License gen failed!
				advise_result.Msg += fmt.Sprintf("%s:License gen failed! ", item.Auth_code)
				wg.Done()
				return
			}
			advise_result.Create_result = true
			advise_result.File_name = data.Oa_id + FORMAT
			wg.Done()
		}()
		wg.Wait()
		if advise_result.Msg == "" { //如果最后Msg中没有写入err，则写入SUCCESS
			advise_result.Msg = SUCCESS
		}
	}
	//记录返回的结果信息
	advise_result.Log()
	//返回License生成结果
	if err := ClientPostAdviseResult(advise_result); err != nil {
		return
	}
	//上传License文件
	if err := ClientUploadLicense(data); err != nil {
		return
	}
}
