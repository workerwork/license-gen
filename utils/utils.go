package utils

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"gopkg.in/ffmt.v1"
	"math/big"
	"os"
	"os/exec"
	"strings"
)

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func Exec(cmd string) (output []string, err error) {
	//fmt.Println(cmd)
	_output, err := exec.Command("/bin/bash", "-c", cmd).Output()
	//time.Sleep(1 * time.Second)
	//fmt.Println("this :", _output)
	output = strings.Split(string(_output), "\n")
	return
}

func Run(cmd string) (err error) {
	//fmt.Println(cmd)
	var stderr bytes.Buffer
	command := exec.Command("/bin/bash", "-c", cmd)
	command.Stderr = &stderr
	err = command.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	return nil
}

func CreateRandomString(len int) string {
	var container string
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(str[randomInt.Int64()])
	}
	return container
}

func RandomString() {
	randomStr := CreateRandomString(15)
	ffmt.P(randomStr)
	//return str
}
