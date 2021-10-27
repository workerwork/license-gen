package server

import (
	"encoding/xml"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"io/ioutil"
	"license-gen/conf"
	"license-gen/utils"
	"os"
	"strconv"
	"time"
)

const (
	DIR     string = "/tmp/license-gen"
	XML_OUT string = "output.xml"
	XML_ENC string = "output.xml.enc"
	LIC_BIN string = "licence.bin"
	EXEC    string = "licence"
	XML_NP  string = "http://www.baicells.com/product"
)

type Product struct {
	Code        string `xml:"code"`
	Version     string `xml:"version"`
	Auth_code   string `xml:"auth_code"`
	Create_time string `xml:"create_time"`
}

type Parameters struct {
	Version     string `xml:"version"`
	Type        string `xml:"type"`
	Sn          string `xml:"sn"`
	Total_time  uint   `xml:"total_time"`
	Use_time    uint   `xml:"use_time"`
	Max_ue_num  uint   `xml:"max_ue_num"`
	Max_enb_num uint   `xml:"max_enb_num"`
}

type License struct {
	XMLName        xml.Name   `xml:"license"`
	XMLNs          string     `xml:"xmlns,attr"`
	Product        Product    `xml:"product"`
	Parameters     Parameters `xml:"parameters"`
	Path_oa_id     string
	Path_auth_code string
}

// @function NewLic
// @description 从xml文件读取数据初始化License结构
// @param ""
// @return *License, error
func NewLic() (lic *License, err error) {
	XmlParam, err := ioutil.ReadFile(conf.LicenseConf.Src)
	if err != nil {
		log.Error().Err(err).Str("func", "ioutil.ReadFile()").Msg("xml文件读取失败")
		return nil, err
	}
	log.Debug().Msgf("Read XML file from %s\n%s", conf.LicenseConf.Src, string(XmlParam))

	err = xml.Unmarshal(XmlParam, &lic)
	if err != nil {
		log.Error().Err(err).Str("func", "xml.Unmarshal()").Msg("Unmarshal error")
		return nil, err
	}
	log.Info().Msg("Unmarshal from XML file success!")
	return lic, nil
}

// @function PathOaId
// @description 把oa_id字段值写入License实例
// @param *License, string
// @return *License
func (l *License) PathOaId(oa_id string) *License {
	l.Path_oa_id = oa_id
	return l
}

// @function PathAuthCode
// @description 把PathAuthCode字段值写入License实例
// @param *License, string
// @return *License
func (l *License) PathAuthCode(auth_code string) *License {
	l.Path_auth_code = auth_code
	return l
}

// @function XmlNs
// @description 把XmlNs字段值写入License实例
// @param *License, string
// @return *License
func (l *License) XmlNs(code string) *License {
	var ns string
	if code == "Bai5GC" {
		ns = "5gc"
	} else if code == "BaiWCG" {
		ns = "egw"
	} else if code == "BaiEPC" {
		ns = "epc"
	}
	l.XMLNs = fmt.Sprintf("%s/%s", XML_NP, ns)
	log.Debug().Str("xmlns", ns).Msg("Set xmlns success!")
	return l
}

// @function Code
// @description 把Code字段值写入License实例
// @param *License, string
// @return *License
func (l *License) Code(code string) *License {
	if code != "" {
		log.Debug().Str("code", code).Msg("Set code success!")
		l.Product.Code = code
	}
	return l
}

// @function ProductVersion
// @description 把Version字段值写入License实例
// @param *License, string
// @return *License
func (l *License) ProductVersion(version string) *License {
	if version != "" {
		log.Debug().Str("product-version", version).Msg("Set product-version success!")
		l.Product.Version = version
	}
	return l
}

// @function AuthCode
// @description 把Auth_code字段值写入License实例
// @param *License, string
// @return *License
func (l *License) AuthCode(auth_code string) *License {
	if auth_code != "" {
		log.Debug().Str("auth_code", auth_code).Msg("Set auth_code success!")
		l.Product.Auth_code = auth_code
	}
	return l
}

// @function CreateTime
// @description 把Create_time字段值写入License实例
// @param *License, string
// @return *License
func (l *License) CreateTime() *License {
	create_time := time.Now().Format("20060102")
	log.Debug().Str("create_time", create_time).Msg("Set create_time success!")
	l.Product.Create_time = create_time
	return l
}

// @function ParametersVersion
// @description 把Version字段值写入License实例
// @param *License, string
// @return *License
func (l *License) ParametersVersion(version string) *License {
	if version != "" {
		log.Debug().Str("parameters-version", version).Msg("Set parameters-version success!")
		l.Parameters.Version = version
	}
	return l
}

// @function Type
// @description 把Type字段值写入License实例
// @param *License, string
// @return *License
func (l *License) Type(ptype string) *License {
	if ptype != "" {
		log.Debug().Str("type", ptype).Msg("Set type success!")
		l.Parameters.Type = ptype
	}
	return l
}

// @function Sn
// @description 把Sn字段值写入License实例
// @param *License, string
// @return *License
func (l *License) Sn(sn string) *License {
	if sn != "" {
		log.Debug().Str("sn", sn).Msg("Set sn success!")
		l.Parameters.Sn = sn
	}
	return l
}

// @function TotalTime
// @description 把Total_time字段值写入License实例
// @param *License, string
// @return *License
func (l *License) TotalTime(total_time uint) *License {
	if total_time != 0 {
		log.Debug().Str("total_time", string(strconv.Itoa(int(total_time)))).Msg("Set total_time success!")
		l.Parameters.Total_time = total_time
	}
	return l
}

// @function UseTime
// @description 把Use_time字段值写入License实例
// @param *License, string
// @return *License
func (l *License) UseTime(use_time uint) *License {
	if use_time != 0 {
		log.Debug().Str("use_time", string(strconv.Itoa(int(use_time)))).Msg("Set use_time success!")
		l.Parameters.Use_time = use_time
	}
	return l
}

// @function MaxUeNum
// @description 把Max_ue_num字段值写入License实例
// @param *License, string
// @return *License
func (l *License) MaxUeNum(max_ue_num uint) *License {
	if max_ue_num != 0 {
		log.Debug().Str("max_ue_num", string(strconv.Itoa(int(max_ue_num)))).Msg("Set max_ue_num success!")
		l.Parameters.Max_ue_num = max_ue_num
	}
	return l
}

// @function MaxEnbNum
// @description 把Max_enb_num字段值写入License实例
// @param *License, string
// @return *License
func (l *License) MaxEnbNum(max_enb_num uint) *License {
	if max_enb_num != 0 {
		log.Debug().Str("max_enb_num", string(strconv.Itoa(int(max_enb_num)))).Msg("Set max_enb_num success!")
		l.Parameters.Max_enb_num = max_enb_num
	}
	return l
}

// @function ToXML
// @description 把License实例转换为xml文件
// @param *License
// @return error
func (l *License) ToXML() error {
	output, err := xml.Marshal(l)
	if err != nil {
		//log.Fatal().Str("func", "ToXML()").Msg("Marshal error!")
		log.Error().Err(err).Str("func", "xml.Marshal()").Msg("Marshal error!")
		return err
	}
	//log.Debug().Msgf("New XML str:\n%s", string(output))
	//path := conf.LicenseConf.Dst
	//str := utils.CreateRandomString(6)

	//清理上次遗留文件
	os.RemoveAll(DIR)
	//定义此次auth_code文件路径
	dir := DIR + "/" + l.Path_oa_id + "/" + l.Path_auth_code
	err = os.MkdirAll(dir, 0777) //此处未判断文件路径是否已经存在
	if err != nil {
		log.Error().Err(err).Str("func", "os.MkdirAll()").Msg("创建文件路径失败!")
		return err
	}
	out_xml := dir + "/" + XML_OUT
	err = ioutil.WriteFile(out_xml, output, 0666)
	if err != nil {
		log.Error().Err(err).Str("func", "ioutil.WriteFile()").Msg("WriteFile error!")
		return err
	}
	log.Info().Str("path", out_xml).Msg("Marshal to XML file success!")
	return nil
}

// @function GenLic
// @description 从xml文件制作License
// @param *License
// @return error
func (l *License) GenLic() error {
	dir := DIR + "/" + l.Path_oa_id + "/" + l.Path_auth_code
	out_exec := dir + "/" + EXEC
	exec, _ := os.Create(out_exec)
	os.Chmod(out_exec, 0755)
	f, _ := os.OpenFile(conf.LicenseConf.Exec, os.O_APPEND, 0666)
	io.Copy(exec, f)
	exec.Close()
	out_xml := dir + "/" + XML_OUT
	cmd := out_exec + " -E " + out_xml
	log.Info().Str("cmd", cmd).Msgf("Starting exec cmd!...")
	//time.Sleep(time.Duration(2)*time.Second)
	err := utils.Run(cmd)
	if err != nil {
		log.Error().Err(err).Str("func", "utils.Run()").Msg("GenLic error!")
		return err
	}
	out_bin := dir + "/" + LIC_BIN
	os.Rename(dir+"/"+XML_ENC, out_bin)
	os.Remove(out_exec)
	os.Remove(out_xml)
	cur_dir, _ := os.Getwd()
	os.Chdir(DIR)
	utils.Zip(l.Path_oa_id, l.Path_oa_id+FORMAT)
	os.Chdir(cur_dir)
	return nil
}
