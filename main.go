package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"encoding/json"
	"pkg"
)

var (
	path = flag.String("path", "", "yapi-sync.json的目录")
	token = flag.String("token", "", "默认使用配置的token")
)


func main() {
	flag.Parse()
	if *path == "" {
		fmt.Println("path cannot be empty!")
		return
	}
	configPath := *path + "/yapi-sync.json"
	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Printf("read file[%s] failed err=%v", configPath, err)
		return
	}
	var config pkg.Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		fmt.Printf("read file[%s] Unmarshal json failed err=%v content=%s", *path, err, string(bytes))
		return
	}
	fmt.Printf("load config:%s \ncontent: %+v \n", configPath, config)
	upload(config)
}


func upload(config pkg.Config)  {

	files := strings.Split(config.File, ",")
	for _, fileName := range files {

		data := pkg.Read(fileName, path)
		apiList := pkg.Parse(data)
		for _, api := range apiList{
			syncApiDocument(api, config)
		}
	}
}

func syncApiDocument(apiDocument pkg.YapiDocument, config pkg.Config) {

	apiDocument.Id = 0
	apiDocument.ReqBodyType = "json"

	apiDocument.Token = config.Token
	if *token == "" {
		apiDocument.Token = *token
	}

	//var values = pkg.StructToMap(apiDocument)
	//fmt.Println("body", values)

	bt, _ := json.Marshal(apiDocument)
	data := string(bt)

	//fmt.Println(data, strings.NewReader(data))

	client := &http.Client{}
	uri := config.Server + "/api/interface/save"
	fmt.Println("uri", uri, data)
	req, err := http.NewRequest("POST", uri, strings.NewReader(data))
	if err != nil {
		fmt.Printf("http.NewRequest failed err=%v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Close = true
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("client.Do failed err=%v", err)
		return
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read responde failed err=%v", err)
		return
	}
	fmt.Println(string(result))
}

func newApiDocument()  {
	
}


func saveApiDocument()  {

}