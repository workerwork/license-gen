package server

import (
	"encoding/json"
	//"fmt"
	"errors"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"license-gen/conf"
	"net/http"
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

func ClientGetInfo() (Data, error) {
	data := Data{}
	client := &http.Client{}
	request, _ := http.NewRequest("GET", conf.URL, nil)
	request.Header.Set("Connection", "keep-alive")
	response, _ := client.Do(request)
	if response.StatusCode == 200 {
		body, _ := ioutil.ReadAll(response.Body)
		//fmt.Println(string(body))
		err := json.Unmarshal(body, &data)
		if err != nil {
			log.Debug().Str("func", "json.Unmarshal()").Msg("Unmarshal error!")
		} else {
			log.Info().Msgf("Get data from %s\n%+v", conf.URL, data)
			return data, nil
		}
	} else {
		log.Error().Msg("http response error!")
	}
	return Data{}, errors.New("http something is wrong!")
}
