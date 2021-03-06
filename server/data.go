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
	"mime/multipart"
	//"net/url"
	"io"
	"os"
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

// @function String
// @description 格式化输出Data结构实例
// @param *Data
// @return string
func (data *Data) String() string {
	str, err := utils.OutString(data)
	if err != nil {
		log.Error().Err(err).Str("func", "Data::utils.OutString()").Msg("marshal error!")
		return ""
	}
	return str
}

// @function String
// @description 格式化输出AdviseResult结构实例
// @param *AdviseResult
// @return string
func (advise_result *AdviseResult) String() string {
	str, err := utils.OutString(advise_result)
	if err != nil {
		log.Error().Err(err).Str("func", "AdviseResult::utils.OutString()").Msg("marshal error!")
		return ""
	}
	return str
}

// @function Log
// @description 记录Data结构实例到指定路径
// @param Data
// @return ""
func (data Data) Log() {
	timeLayoutStr := "2006-01-02 15:04:05"
	timeStr := time.Now().Format(timeLayoutStr)
	formatStr := strings.Repeat("*", 20)
	str := fmt.Sprintf("%s%s%s\n%s\n", formatStr, timeStr, formatStr, data.String())
	utils.WriteFile(conf.ServerConf.Log, str)
}

// @function Log
// @description 记录AdviseResult结构实例到指定路径
// @param AdviseResult
// @return ""
func (advise_result AdviseResult) Log() {
	str := fmt.Sprintf("%s\n", advise_result.String())
	utils.WriteFile(conf.ServerConf.Log, str)
}

// @function ClientGetInfo
// @description 获取data数据
// @param ""
// @return Data, error
func ClientGetInfo() (Data, error) {
	data := Data{}
	client := &http.Client{}
	request, err := http.NewRequest("GET", conf.URL_GET, nil)
	if err != nil {
		log.Error().Err(err).Str("func", "http.NewRequest()").Msg("http error!")
	}
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	response, err := client.Do(request)
	if err != nil {
		log.Error().Err(err).Str("func", "client.Do()").Msg("http error!")
		return Data{}, err
	}
	if response.StatusCode == 200 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Error().Err(err).Str("func", "ioutil.ReadAll()").Msg("IO error!")
			return Data{}, err
		}
		if len(body) == 0 {
			log.Debug().Msg("can't get data from server yet!")
			return data, errors.New("can't get data from server yet!")
		}
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

// @function ClientPostAdviseResult
// @description 返回License制作结果
// @param AdviseResult
// @return error
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

// @function ClientUploadLicense
// @description 上传License文件
// @param Data
// @return error
func ClientUploadLicense(data Data) error {
	result_code := ResultCode{}
	params := make(map[string]string)
	params["oa_id"] = data.Oa_id
	params["apply_type"] = data.Apply_type
	path := DIR + "/" + data.Oa_id + ".zip"
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("uploadFile", path)
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return err
	}
	request, err := http.NewRequest("POST", conf.URL_POST2, body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(request)
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
			log.Info().Msgf("Get data from %s: %+v", conf.URL_POST2, result_code)
		}
	} else {
		log.Error().Msg("http response error!")
		return err
	}
	if result_code.Result_code != 2000 {
		return errors.New("Server response err!")
	}
	defer resp.Body.Close()
	return nil
}
