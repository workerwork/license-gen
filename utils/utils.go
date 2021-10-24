package utils

import (
	"archive/zip"
	//"bufio"
	"bytes"
	"crypto/rand"
	"fmt"
	"gopkg.in/ffmt.v1"
	"io"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
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

// srcFile could be a single file or a directory
func Zip(srcFile string, destZip string) error {
	zipfile, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	filepath.Walk(srcFile, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = path
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
		}
		return err
	})

	return err
}

func Unzip(zipFile string, destDir string) error {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		fpath := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			inFile, err := f.Open()
			if err != nil {
				return err
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func WriteFile(path string, buf string) {
	var f *os.File
	var err error
	if len(buf) == 0 {
		return
	}
	if !CheckFileIsExist(path) {
		f, err = os.Create(path)
	} else {
		f, err = os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0666)
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	//使用完毕，需要关闭文件
	defer f.Close()
	_, err = f.WriteString(buf)
	if err != nil {
		fmt.Println("err = ", err)
	}
}
