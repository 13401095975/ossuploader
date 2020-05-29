package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type ossConfig struct {
	EndPoint        string
	AccessKeyId     string
	AccessKeySecret string
}

var (
	h bool

	bucket string
	config string
	dir    string

	OssConfig ossConfig
)

func init() {

	flag.StringVar(&bucket, "bucket", "", "bucket name")
	flag.StringVar(&config, "config", "./.oss.config", "config file path, default path ./.oss.config")
	flag.StringVar(&dir, "dir", "", "upload fir")

	flag.Usage = usage
}

func HandleError(err error) {
	fmt.Println("Error:", err)
	os.Exit(-1)
}

func main() {
	flag.Parse()
	if h {
		flag.Usage()
		os.Exit(0)
	}

	if bucket == "" {
		fmt.Println("Error: bucketName required")
		os.Exit(-1)
	}
	if len(dir) == 0 {
		fmt.Println("Error: dirPath required")
		os.Exit(-1)
	}
	allBytes, err := ioutil.ReadFile(config)
	if err != nil {
		HandleError(err)
	}

	err = json.Unmarshal(allBytes, &OssConfig)
	if err != nil {
		HandleError(err)
	}

	client := ossClient(OssConfig.EndPoint, OssConfig.AccessKeyId, OssConfig.AccessKeySecret)
	clearBucketFiles(client, bucket)
	UploadAllFile(client, bucket, dir)
}

func usage() {
	fmt.Fprintf(os.Stderr, `ossuploader
Usage: ossuploader [-h] [-bucket bucketName] [-config configFile] [-dir dirPath]

Options:
`)
	flag.PrintDefaults()
}

func ossClient(endPoint, keyId, keySecret string) *oss.Client {
	client, err := oss.New(endPoint, keyId, keySecret)
	if err != nil {
		HandleError(err)
	}
	return client
}

func UploadAllFile(client *oss.Client, bucketName string, dir string) {

	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		HandleError(err)
	}

	xfiles, _ := GetAllFiles(dir)
	for _, file := range xfiles {

		key := strings.Replace(file, dir+"/", "", 1)
		fmt.Printf("获取的文件为[%s], key:[%s]\n", file, key)
		err = bucket.PutObjectFromFile(key, file)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(-1)
		}
	}
	fmt.Printf("上传成功，文件数为[%d]\n", len(xfiles))
}

func GetAllFiles(dirPth string) (files []string, err error) {
	var dirs []string
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	PthSep := string(os.PathSeparator)

	for _, fi := range dir {
		if fi.IsDir() { // 目录, 递归遍历
			dirs = append(dirs, dirPth+PthSep+fi.Name())
			GetAllFiles(dirPth + PthSep + fi.Name())
		} else {
			if fi.Name()[0] != '.' {
				fullPath := filepath.ToSlash(dirPth + PthSep + fi.Name())
				files = append(files, fullPath)
			}
		}
	}

	// 读取子目录下文件
	for _, table := range dirs {
		temp, _ := GetAllFiles(table)
		for _, temp1 := range temp {
			fullPath := filepath.ToSlash(temp1)
			files = append(files, fullPath)
		}
	}

	return files, nil
}

func clearBucketFiles(client *oss.Client, bucketName string) {

	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		HandleError(err)
	}

	files := []string{}
	// 列举所有文件。
	marker := ""
	for {
		lsRes, err := bucket.ListObjects(oss.Marker(marker))
		if err != nil {
			HandleError(err)
		}

		// 打印列举文件，默认情况下一次返回100条记录。
		for _, object := range lsRes.Objects {
			files = append(files, object.Key)
		}

		if lsRes.IsTruncated {
			marker = lsRes.NextMarker
		} else {
			break
		}
	}

	if len(files) > 0 {
		// 返回删除成功的文件。
		delRes, err := bucket.DeleteObjects(files)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(-1)
		}
		fmt.Println("Delete File:", delRes.DeletedObjects)
	}

}
