package server

import (
	//"fmt"
	"time"
)

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
	for _, item := range data.Item_list {
		go func() {
			lic, err := NewLic()
			if err != nil {
				return
			}
			lic.XmlNs(data.Product).
				Code(data.Product).
				ProductVersion(item.Version).
				AuthCode(item.Auth_code).
				CreateTime().
				TotalTime(item.Total_time).
				MaxUeNum(item.Max_ue_num).
				MaxEnbNum(item.Max_enb_num)
			//fmt.Println(lic)
			tmp_dir, err := lic.ToXML()
			if err != nil {
				return
			}
			GenLic(tmp_dir)
		}()
	}
}
