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
			lic.XmlNs(data.Product)
			lic.Code(data.Product)
			lic.ProductVersion(item.Version)
			lic.AuthCode(item.Auth_code)
			lic.CreateTime()
			lic.TotalTime(item.Total_time)
			lic.MaxUeNum(item.Max_ue_num)
			lic.MaxEnbNum(item.Max_enb_num)
			//fmt.Println(lic)
			lic.ToXML()
			lic.GenLic()
		}()
	}
}
