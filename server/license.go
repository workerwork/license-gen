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

var (
	Tmp_dir  string
	Out_exec string
	Out_xml  string
	Out_bin  string
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
	XMLName xml.Name `xml:"license"`
	XMLNs      string     `xml:"xmlns,attr"`
	Product    Product    `xml:"product"`
	Parameters Parameters `xml:"parameters"`
}

func NewLic() (lic *License, err error) {
	XmlParam, err := ioutil.ReadFile(conf.LicenseConf.Src)
	if err != nil {
		log.Error().Err(err).Str("func", "New()").Msg("xml文件读取失败")
		return nil, err
	}
	log.Debug().Msgf("Read XML file from %s\n%s", conf.LicenseConf.Src, string(XmlParam))

	err = xml.Unmarshal(XmlParam, &lic)
	if err != nil {
		log.Error().Err(err).Str("func", "New()").Msg("Unmarshal error")
		return nil, err
	}
	log.Info().Msg("Unmarshal from XML file success!")
	return lic, nil
}

func (l *License) XmlNs(code string) *License {
	var ns string
	if code == "Bai5GC" {
		ns = "5gc"
	} else if code == "BaiWCG" {
		ns = "egw"
	}
	log.Debug().Str("xmlns", ns).Msg("Set xmlns success!")
	l.XMLNs = fmt.Sprintf("%s/%s", XML_NP, ns)
	return l
}

func (l *License) Code(code string) *License {
	if code != "" {
		log.Debug().Str("code", code).Msg("Set code success!")
		l.Product.Code = code
	}
	return l
}

func (l *License) ProductVersion(version string) *License {
	if version != "" {
		log.Debug().Str("product-version", version).Msg("Set product-version success!")
		l.Product.Version = version
	}
	return l
}

func (l *License) AuthCode(auth_code string) *License {
	if auth_code != "" {
		log.Debug().Str("auth_code", auth_code).Msg("Set auth_code success!")
		l.Product.Auth_code = auth_code
	}
	return l
}

func (l *License) CreateTime() *License {
	create_time := time.Now().Format("20060102")
	log.Debug().Str("create_time", create_time).Msg("Set create_time success!")
	l.Product.Create_time = create_time
	return l
}

func (l *License) ParametersVersion(version string) *License {
	if version != "" {
		log.Debug().Str("parameters-version", version).Msg("Set parameters-version success!")
		l.Parameters.Version = version
	}
	return l
}

func (l *License) Type(ptype string) *License {
	if ptype != "" {
		log.Debug().Str("type", ptype).Msg("Set type success!")
		l.Parameters.Type = ptype
	}
	return l
}

func (l *License) Sn(sn string) *License {
	if sn != "" {
		log.Debug().Str("sn", sn).Msg("Set sn success!")
		l.Parameters.Sn = sn
	}
	return l
}

func (l *License) TotalTime(total_time uint) *License {
	if total_time != 0 {
		log.Debug().Str("total_time", string(strconv.Itoa(int(total_time)))).Msg("Set total_time success!")
		l.Parameters.Total_time = total_time
	}
	return l
}

func (l *License) UseTime(use_time uint) *License {
	if use_time != 0 {
		log.Debug().Str("use_time", string(strconv.Itoa(int(use_time)))).Msg("Set use_time success!")
		l.Parameters.Use_time = use_time
	}
	return l
}

func (l *License) MaxUeNum(max_ue_num uint) *License {
	if max_ue_num != 0 {
		log.Debug().Str("max_ue_num", string(strconv.Itoa(int(max_ue_num)))).Msg("Set max_ue_num success!")
		l.Parameters.Max_ue_num = max_ue_num
	}
	return l
}

func (l *License) MaxEnbNum(max_enb_num uint) *License {
	if max_enb_num != 0 {
		log.Debug().Str("max_enb_num", string(strconv.Itoa(int(max_enb_num)))).Msg("Set max_enb_num success!")
		l.Parameters.Max_enb_num = max_enb_num
	}
	return l
}

func (l *License) ToXML() (*License, error) {
	output, err := xml.Marshal(l)
	if err != nil {
		//log.Fatal().Str("func", "ToXML()").Msg("Marshal error!")
		log.Error().Err(err).Str("func", "ToXML()").Msg("Marshal error!")
		return nil, err
	}
	//log.Debug().Msgf("New XML str:\n%s", string(output))
	//path := conf.LicenseConf.Dst
	str := utils.CreateRandomString(6)
	Tmp_dir = DIR + "/" + str
	err = os.MkdirAll(Tmp_dir, 0777) //此处未判断文件夹是否已经存在
	if err != nil {
		log.Error().Err(err).Str("func", "ToXML()").Msg("创建文件路径失败!")
		return nil, err
	}
	Out_xml = Tmp_dir + "/" + XML_OUT
	err = ioutil.WriteFile(Out_xml, output, 0666)
	if err != nil {
		log.Error().Err(err).Str("func", "ToXML()").Msg("WriteFile error!")
		return nil, err
	}
	log.Info().Str("path", Out_xml).Msg("Marshal to XML file success!")
	return l, nil
}

func (l *License) GenLic() error {
	Out_exec = Tmp_dir + "/" + EXEC
	exec, _ := os.Create(Out_exec)
	os.Chmod(Out_exec, 0755)
	f, _ := os.OpenFile(conf.LicenseConf.Exec, os.O_APPEND, 0666)
	io.Copy(exec, f)
	exec.Close()
	cmd := Out_exec + " -E " + Out_xml
	log.Info().Str("cmd", cmd).Msgf("Starting exec cmd!...")
	//time.Sleep(time.Duration(2)*time.Second)
	err := utils.Run(cmd)
	if err != nil {
		log.Error().Err(err).Str("func", "GenLic()").Msg("GenLic error!")
		return err
	}
	Out_bin = Tmp_dir + "/" + LIC_BIN
	os.Rename(Tmp_dir+"/"+XML_ENC, Out_bin)
	return nil
}
