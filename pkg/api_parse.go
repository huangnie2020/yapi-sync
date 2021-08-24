package pkg

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"encoding/json"

	"github.com/goinggo/mapstructure"
)

func Read(fileName string, path *string) map[string]interface{} {
	var bytes []byte
	var err error
	if strings.Index(fileName, "http") != -1 {

		client := &http.Client{}
		uri := fileName
		req, err := http.NewRequest("GET", uri, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Close = true
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("client.Do failed err=%v", err)
			return nil
		}
		defer resp.Body.Close()
		bytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("read responde failed err=%v", err)
			return nil
		}
	} else {
		filePath := *path + "/" + fileName
		bytes, err = ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Printf("read file[%s] failed err=%v", filePath, err)
			return nil
		}
	}

	js := strings.Replace(string(bytes), "$ref", "ref", -1)

	var data map[string]interface{}

	json.Unmarshal([]byte(js), &data)

	return data
}

func Parse(data map[string]interface{}) []YapiDocument {

	var apiList []YapiDocument

	// map -> struct
	var swagJson SwagJson
	mapstructure.Decode(data, &swagJson)

	for addrPath, documentMap := range swagJson.Paths {

		var yapiDocument = YapiDocument{
			Token:               "",
			Id:                  0,
			CatId:               0,
			CatName:             "",
			Tag:                 []string{},
			Title:               "",
			Desc:                "",
			Method:              "",
			Path:                "",
			ProjectId:           0,
			ReqBodyType:         "json",
			ReqBodyIsJsonSchema: "true",
		}

		// 每个路由可能有多个方法
		for method, document := range documentMap {

			var swagPathsDocument SwagPathsDocument

			// map -> struct
			mapstructure.Decode(document, &swagPathsDocument)

			yapiDocument.Path = addrPath
			yapiDocument.Method = method
			yapiDocument.Title = swagPathsDocument.Summary
			yapiDocument.Desc = swagPathsDocument.Description
			yapiDocument.Tag = swagPathsDocument.Tags

			// 参数
			for _, parameter := range swagPathsDocument.Parameters {
				var isRequired int
				if parameter.Required == "true" {
					isRequired = 1
				} else {
					isRequired = 0
				}

				if parameter.In == "body" {
					// body 对象
					yapiDocument.ReqBodyOther = parseSwagScheme(parameter.Schema, swagJson.Definitions)
				} else {

					reqField := &YapiQueryFieldItem{
						Name:     parameter.Name,
						Desc:     parameter.Description,
						Required: isRequired,
						Value:    "",
					}

					if parameter.In == "path" {
						yapiDocument.ReqParams = append(yapiDocument.ReqParams, reqField)
					} else {
						if parameter.Type == "array" {
							reqField.Items = parseSwagItems(parameter.Items, swagJson.Definitions)
						}
						yapiDocument.ReqQuery = append(yapiDocument.ReqQuery, reqField)
					}
				}
			}

			// 可能有多种响应
			for status, response := range swagPathsDocument.Responses {
				if status != "default" {
					yapiDocument.ResBody = parseSwagScheme(response.Schema, swagJson.Definitions)
					apiList = append(apiList, yapiDocument)
				}
			}
		}
	}

	return apiList
}

func StructToMap(st interface{}) map[string]interface{} {
	data := make(map[string]interface{})
	j, _ := json.Marshal(st)
	json.Unmarshal(j, &data)
	return data
}

func parseSwagRefKey(Ref string) string {
	return strings.Replace(Ref, "#/definitions/", "", -1)
}

func parseSwagObjectProperties(properties map[string]interface{}, definitions map[string]SwagFieldObject) map[string]interface{} {

	for field, property := range properties {

		var propertyStruct SwagFieldObjectPropertyItem
		mapstructure.Decode(property, &propertyStruct)

		if propertyStruct.Ref != "" {
			// 对象字段解析
			key := parseSwagRefKey(propertyStruct.Ref)
			obj := definitions[key]
			obj.Properties = parseSwagObjectProperties(obj.Properties, definitions)
			if obj.Title == "" {
				obj.Title = propertyStruct.Title
			}
			if obj.Description == "" {
				obj.Description = propertyStruct.Description
			}

			properties[field] = StructToMap(obj)
		} else {
			// 数组字段解析
			var itemsStruct SwagFieldArray
			mapstructure.Decode(propertyStruct.Items, &itemsStruct)
			if &itemsStruct != nil {
				propertyStruct.Items = parseSwagItems(propertyStruct.Items, definitions)
			} else if propertyStruct.Type == "string" {
				propertyStruct.Format = "string"
			}

			properties[field] = StructToMap(propertyStruct)
		}
	}

	return properties
}

func parseSwagItems(items map[string]interface{}, definitions map[string]SwagFieldObject) map[string]interface{} {
	var itemsStruct SwagFieldArray
	mapstructure.Decode(items, &itemsStruct)
	if &itemsStruct == nil || itemsStruct.Ref == "" {
		return items
	}

	// 对象元素
	key := parseSwagRefKey(itemsStruct.Ref)
	obj := definitions[key]
	obj.Properties = parseSwagObjectProperties(obj.Properties, definitions)

	return StructToMap(obj)
}

func parseSwagScheme(scheme map[string]interface{}, definitions map[string]SwagFieldObject) map[string]interface{} {

	var schemeStruct SwagFieldSchema
	mapstructure.Decode(scheme, &schemeStruct)

	if &schemeStruct == nil {
		return nil
	}

	if schemeStruct.Ref != "" {
		// 关联对象元素
		key := parseSwagRefKey(schemeStruct.Ref)
		obj := definitions[key]
		obj.Properties = parseSwagObjectProperties(obj.Properties, definitions)
		if obj.Title == "" {
			obj.Title = schemeStruct.Title
		}
		if obj.Description == "" {
			obj.Description = schemeStruct.Description
		}

		return StructToMap(obj)
	} else if schemeStruct.Type == "array" {
		// 实体数组
		schemeStruct.Items = parseSwagItems(schemeStruct.Items, definitions)
		return StructToMap(schemeStruct)
	} else {
		return nil
	}
}
