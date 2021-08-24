package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"encoding/json"
	"pkg"
)

var (
	path  = flag.String("path", "", "yapi-sync.json的目录")
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

func upload(config pkg.Config) {

	var projectId int = 11
	menuList := getCatMenuList(projectId, config)
	oldApiList := getApiList(1, 10000, projectId, config)

	files := strings.Split(config.File, ",")
	for _, fileName := range files {

		data := pkg.Read(fileName, path)
		apiList := pkg.Parse(data)
		for _, apiDocument := range apiList {

			fmt.Println("cc:", apiDocument.Tag)
			catName := apiDocument.Tag[0]
			menu := menuList[catName]
			if menu.Name == "" {
				menuList[catName] = *createCateMenu(projectId, catName, config)
			}

			fmt.Println("cc:", menu.Name, catName)

			key := apiDocument.Method + apiDocument.Path
			old := oldApiList[key]

			syncApiDocument(menu, old, apiDocument, config)
		}
	}
}

func syncApiDocument(menu pkg.YapiMenuItem, apiOld pkg.YapiApiListItem, apiDocument pkg.YapiDocument, config pkg.Config) {

	apiDocument.Id = apiOld.Id
	apiDocument.ProjectId = 11
	apiDocument.ReqBodyType = "json"

	apiDocument.Token = config.Token
	if apiDocument.Token == "" {
		apiDocument.Token = *token
	}

	apiDocument.CatId = menu.Id

	var values = pkg.StructToMap(apiDocument)

	params := make(map[string]interface{})

	for k, v := range values {
		if v == nil && (k == "req_query" || k == "req_params" || k == "req_headers" || k == "req_body_form") {
			v = []string{}
		} else if v == nil && k == "req_body_other" {
			v = "{}"
		} else if k == "res_body" {
			vb, _ := json.Marshal(v)
			v = string(vb)
			// v = strconv.Quote(v)
		}
		params[k] = v
	}

	fmt.Println("body", params)

	bt, _ := json.Marshal(params)
	body := bytes.NewBuffer(bt)

	client := &http.Client{}
	uri := config.Server + "/api/interface/save"
	req, err := http.NewRequest("POST", uri, body)
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

func getApiList(page int, limit int, projectId int, config pkg.Config) map[string]pkg.YapiApiListItem {

	client := &http.Client{}
	uri := arrayToString([]string{
		config.Server,
		"/api/interface/getCatMenu?page=",
		strconv.Itoa(page),
		"&limit=",
		strconv.Itoa(limit),
		"&project_id=",
		strconv.Itoa(projectId),
		"&token=",
		config.Token,
	})

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		fmt.Printf("http.NewRequest failed err=%v", err)
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	req.Close = true
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("client.Do failed err=%v", err)
		return nil
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read responde failed err=%v", err)
		return nil
	}
	fmt.Println(string(result))

	var response pkg.YapiListResponse
	json.Unmarshal(result, &response)

	var responseApiList pkg.YapiApiList

	j, _ := json.Marshal(response.Data)
	json.Unmarshal(j, &responseApiList)

	var apiList map[string]pkg.YapiApiListItem

	for _, v := range responseApiList.List {

		var apiItem pkg.YapiApiListItem
		b, _ := json.Marshal(v)
		json.Unmarshal(b, &apiItem)

		key := apiItem.Method + apiItem.Path

		apiList[key] = apiItem
	}

	return apiList

}

func createCateMenu(projectId int, catName string, config pkg.Config) *pkg.YapiMenuItem {

	params := make(map[string]interface{})
	params["project_id"] = projectId
	params["name"] = catName
	params["desc"] = nil
	params["token"] = config.Token

	fmt.Println("body", params)

	bt, _ := json.Marshal(params)
	body := bytes.NewBuffer(bt)

	client := &http.Client{}
	uri := config.Server + "/api/interface/add_cat"
	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		fmt.Printf("http.NewRequest failed err=%v", err)
		return nil
	}

	req.Header.Set("Content-Type", "application/json")
	req.Close = true
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("client.Do failed err=%v", err)
		return nil
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read responde failed err=%v", err)
		return nil
	}
	fmt.Println(string(result))

	var response pkg.YapiDetailResponse
	json.Unmarshal(result, &response)

	var menuItem pkg.YapiMenuItem

	j, _ := json.Marshal(response.Data)
	json.Unmarshal(j, &menuItem)

	return &menuItem
}

func getCatMenuList(prodjectId int, config pkg.Config) map[string]pkg.YapiMenuItem {

	client := &http.Client{}

	uri := arrayToString([]string{
		config.Server,
		"/api/interface/getCatMenu?project_id=",
		strconv.Itoa(prodjectId),
		"&token=",
		config.Token,
	})
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		fmt.Printf("http.NewRequest failed err=%v", err)
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	req.Close = true
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("client.Do failed err=%v", err)
		return nil
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read responde failed err=%v", err)
		return nil
	}
	fmt.Println(string(result))

	var response pkg.YapiListResponse
	json.Unmarshal(result, &response)

	var memuList = make(map[string]pkg.YapiMenuItem)

	for _, v := range response.Data {

		var menu pkg.YapiMenuItem
		b, _ := json.Marshal(v)
		json.Unmarshal(b, &menu)

		memuList[menu.Name] = menu
	}

	return memuList
}

func arrayToString(arr []string) string {
	var result string
	//遍历数组中所有元素追加成string
	for _, i := range arr {
		result += string(i)
	}
	return result
}
