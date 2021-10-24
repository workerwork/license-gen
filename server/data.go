package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"license-gen/conf"
	"license-gen/utils"
	"net/http"
	"strings"
	"time"
	//"unsafe"
)

type Data struct {
	Oa_id      string `json:"oa_id"`
	Apply_type string `json:"apply_type"`
	Applicant  string `json:"applicant"`
	Reason     string `json:"reason"`
	Purpose    string `json:"purpose"`
	Po         string `json:"po"`
	Scope      string `json:"scope"`
	District   string `json:"district"`
	Customer   string `json:"customer"`
	Email      string `json:"email"`
	Product    string `json:"product"`
	Item_list  []Item `json:"item_list"`
}

type Item struct {
	Auth_code   string `json:"auth-code"`
	Version     string `json:"version"`
	Max_enb_num uint   `json:"max-enb-num"`
	Max_ue_num  uint   `json:"max-ue-num"`
	Total_time  uint   `json:"total-time"`
}

type AdviseResult struct {
	Oa_id         string `json:"oa_id"`
	Apply_type    string `json:"apply_type"`
	Create_result bool   `json:"create_result"`
	File_name     string `json:"file_name"`
	Msg           string `json:"msg"`
}

type ResultCode struct {
	Result_code uint   `json:"result_code"`
	Message     string `json:"message"`
}

func (data *Data) String() string {
	b, err := json.Marshal(*data)
	if err != nil {
		return fmt.Sprintf("%+v", *data)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("%+v", *data)
	}
	return out.String()
}

func (advise_result *AdviseResult) String() string {
	b, err := json.Marshal(*advise_result)
	if err != nil {
		return fmt.Sprintf("%+v", *advise_result)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("%+v", *advise_result)
	}
	return out.String()
}

// 写log文件
func (data Data) Log() {
	timeLayoutStr := "2006-01-02 15:04:05"
	timeStr := time.Now().Format(timeLayoutStr)
	formatStr := strings.Repeat("*", 20)
	str := fmt.Sprintf("%s%s%s\n%s\n", formatStr, timeStr, formatStr, data.String())
	utils.WriteFile(conf.ServerConf.Log, str)
}

func (advise_result AdviseResult) Log() {
	str := fmt.Sprintf("%s\n", advise_result.String())
	utils.WriteFile(conf.ServerConf.Log, str)
}

func ClientGetInfo() (Data, error) {
	data := Data{}
	client := &http.Client{}
	request, _ := http.NewRequest("GET", conf.URL_GET, nil)
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	response, _ := client.Do(request)
	if response.StatusCode == 200 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Error().Err(err).Str("func", "ioutil.ReadAll()").Msg("IO error!")
			return Data{}, err
		}
		//fmt.Println(string(body))
		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Debug().Str("func", "json.Unmarshal()").Msg("Unmarshal error!")
		} else {
			log.Info().Msgf("Get data from %s\n%+v", conf.URL_GET, data)
			return data, nil
		}
	} else {
		log.Error().Msg("http response error!")
	}
	return Data{}, errors.New("http something is wrong!")
}

func ClientPostAdviseResult(advise_result AdviseResult) error {
	result_code := ResultCode{}
	bytesData, err := json.Marshal(&advise_result)
	if err != nil {
		log.Error().Err(err).Str("func", "json.Marshal()").Msg("Marshal error!")
		return err
	}
	reader := bytes.NewReader(bytesData)
	request, err := http.NewRequest("POST", conf.URL_POST1, reader)
	if err != nil {
		log.Error().Err(err).Str("func", "http.NewReader()").Msg("http error!")
		return err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Error().Err(err).Str("func", "client.Do()").Msg("http error!")
		return err
	}
	//fmt.Println(resp.StatusCode)
	if resp.StatusCode == 200 {
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error().Err(err).Str("func", "ioutil.ReadAll()").Msg("IO error!")
			return err
		}
		err = json.Unmarshal(respBytes, &result_code)
		if err != nil {
			log.Debug().Str("func", "json.Unmarshal()").Msg("Unmarshal error!")
			return err
		} else {
			log.Info().Msgf("Get data from %s: %+v", conf.URL_POST1, result_code)
		}
	} else {
		log.Error().Msg("http response error!")
		return err
	}
	if result_code.Result_code != 2000 {
		return errors.New("Server response err!")
	}
	//fmt.Println("result_code:", result_code)
	return nil
	//byte数组直接转成string，优化内存
	//str := (*string)(unsafe.Pointer(&respBytes))
	//return str
}

func ClientUploadLicenseFile() {}
