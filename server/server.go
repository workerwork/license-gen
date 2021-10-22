package server

import (
	//"fmt"
	"sync"
	"time"
)

var lock sync.Mutex
var wg sync.WaitGroup
var data = Data{}

func Serve() {
	for {
		go run()
		time.Sleep(time.Duration(2) * time.Second)
	}
}

func run() {
	data, err := ClientGetInfo()
	if err != nil {
		return
	}
    /*
	for _, item := range data.Item_list {
		wg.Add(1)
		lock.Lock()
		go func() {
			lic, err := NewLic()
			if err != nil {
				lock.Unlock()
				wg.Done()
				return
			}
			//fmt.Println(lic)
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
				lock.Unlock()
				wg.Done()
				return
			}
			if err := lic.GenLic(); err != nil {
				lock.Unlock()
				wg.Done()
				return
			}
			lock.Unlock()
			wg.Done()
		}()
		wg.Wait()
		//TODO
	}*/
	advise_result := AdviseResult{
		Oa_id:         data.Oa_id,
		Apply_type:    data.Apply_type,
		Create_result: true,
		File_name:     "test",
		Msg:           "success",
	}
	ClientPostAdviseResult(advise_result)
}
