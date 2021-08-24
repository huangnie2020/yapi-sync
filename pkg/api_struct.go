package pkg

// config

type Config struct {
	Type   string `json:"type"`
	Token  string `json:"token"`
	File   string `json:"file"`
	Merge  string `json:"merge"`
	Server string `json:"server"`
}

// swag

type SwagFieldArray struct {
	Type   string `json:"type"`
	Format string `json:"format"`
	Ref    string `json:"ref"` // 数组元素可能是对象，Ref不为空时，便是object数组,且 items = ref
}

type SwagFieldObjectPropertyItem struct {
	Ref         string                 `json:"ref"`   // ，Ref不为空时，该属性便是object元素,且 属性 = ref
	Items       map[string]interface{} `json:"items"` // 当type=array 若为对象数组则items=items[ref]， 若为普通数组则items=items
	Description string                 `json:"description"`
	Title       string                 `json:"title"`
	Format      string                 `json:"format"`
	Type        string                 `json:"type"`
}

// body 才会有对像字段

type SwagFieldObject struct {
	Type        string                 `json:"type"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Properties  map[string]interface{} `json:"properties"` // 属性是混合类型：当元素不存在 Ref时，为普通字段。 当元素存在 Ref时，若有type=array便是数组，若无type值就是object
}

// Parameters 才会有的一个对象载体

type SwagFieldSchema struct {
	SwagFieldObject
	Items map[string]interface{} `json:"items"`
	Ref   string                 `json:"ref"`
}

type SwagPathsDocumentParametersItem struct {
	Ref         string                 `json:"ref"`
	Items       map[string]interface{} `json:"items"`
	Schema      map[string]interface{} `json:"schema"` // 在 in=body 时，post数据都在 schema 对象属性列表， 也即是 yapi 的 req_body_other
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Required    string                 `json:"required"`
	Format      string                 `json:"format"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	In          string                 `json:"in"`
}

type SwagPathsDocumentResponsesItem struct {
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Schema      map[string]interface{} `json:"schema"`
}

type SwagPathsDocument struct {
	Summary     string                                    `json:"summary"`
	Description string                                    `json:"description"`
	OperationId string                                    `json:"operationId"`
	Tags        []string                                  `json:"tags"`
	Parameters  []SwagPathsDocumentParametersItem         `json:"parameters"` // 即是 yapi 的 req_query + req_req_params + req_body_other
	Responses   map[string]SwagPathsDocumentResponsesItem `json:"responses"`  // 即是 yapi 的 req_body
}

type SwagJson struct {
	Paths       map[string]map[string]interface{} `json:"paths"`
	Definitions map[string]SwagFieldObject        `json:"definitions"`
}

// yapi

//req_param 和 req_query  都是 QueryField 组成
//req_body 和 req_body_other  都是 BodyField + BodyFieldObject + BodyFieldArray 混合组成
// 其中 BodyFieldObject 和 BodyFieldArray 可能由 BodyField + BodyFieldObject + BodyFieldArray 混合组成

type YapiQueryFieldItem struct {
	Name     string                 `json:"name"`
	Desc     string                 `json:"desc"`
	Value    string                 `json:"value"`
	Required int                    `json:"required"`
	Items    map[string]interface{} `json:"items"`
}

type YapiBodyFieldItem struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Format      string `json:"format"`
}

// 类型 type=object

type YapiBodyFieldObject struct {
	Ref         string                 `json:"ref"`
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Properties  map[string]interface{} `json:"properties"`
}

// 类型 type=array

type YapiBodyFieldArray struct {
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Items       map[string]interface{} `json:"items"`
}

type YapiDocument struct {
	Token               string                 `json:"token"`
	Id                  int                    `json:"_id"`
	CatId               int                    `json:"catid"`
	CatName             string                 `json:"catname"`
	Tag                 []string               `json:"tag"`
	Title               string                 `json:"title"`
	Desc                string                 `json:"desc"`
	Method              string                 `json:"method"`
	Path                string                 `json:"path"`
	ProjectId           int                    `json:"project_id"`
	ReqBodyIsJsonSchema string                 `json:"req_body_is_json_schema"` // true | false
	ReqBodyType         string                 `json:"req_body_type"`           // json | form
	ReqBodyOther        map[string]interface{} `json:"req_body_other"`          //`json:"req_body_other"`
	ResBody             map[string]interface{} `json:"res_body"`
	ReqBodyForm         []interface{}          `json:"req_body_form"`
	ReqHeaders          []interface{}          `json:"req_headers"`
	ReqParams           []interface{}          `json:"req_params"`
	ReqQuery            []interface{}          `json:"req_query"`
}

// yapi 平台开放api

type YapiDetailResponse struct {
	Errcode int                    `json:"errcode"`
	Errmsg  string                 `json:"errmsg"`
	Data    map[string]interface{} `json:"data"`
}

type YapiListResponse struct {
	Errcode int           `json:"errcode"`
	Errmsg  string        `json:"errmsg"`
	Data    []interface{} `json:"data"`
}

type YapiMenuItem struct {
	Id        int    `json:"_id"`
	ProductId int    `json:"project_id"`
	Name      string `json:"name"`
	Desc      string `json:"desc"`
	Uid       int    `json:"uid"`
	AddTime   int    `json:"add_time"`
	UpTime    int    `json:"up_time"`
}

type YapiApiList struct {
	Count int                    `json:"count"`
	List  map[string]interface{} `json:"list"`
}

type YapiApiListItem struct {
	Id        int    `json:"_id"`
	CatId     int    `json:"CatId"`
	ProductId int    `json:"project_id"`
	Path      string `json:"path"`
	Method    string `json:"method"`
	Title     string `json:"title"`
	Status    string `json:"status"`
	Tag       int    `json:"tag"`
	AddTime   int    `json:"add_time"`
	UpTime    int    `json:"up_time"`
}

type YapiApiDetail struct {
	Id                  int                    `json:"_id"`
	CatId               int                    `json:"CatId"`
	ProductId           int                    `json:"project_id"`
	QueryPath           map[string]interface{} `json:"query_path"`
	Method              string                 `json:"method"`
	Tag                 []string               `json:"tag"`
	Title               string                 `json:"title"`
	Desc                string                 `json:"desc"`
	ReqBodyIsJsonSchema string                 `json:"req_body_is_json_schema"` // true | false
	ReqBodyType         string                 `json:"req_body_type"`           // json | form
	ReqBodyOther        map[string]interface{} `json:"req_body_other"`          //`json:"req_body_other"`
	ResBody             map[string]interface{} `json:"res_body"`
	ReqBodyForm         []interface{}          `json:"req_body_form"`
	ReqHeaders          []interface{}          `json:"req_headers"`
	ReqParams           []interface{}          `json:"req_params"`
	ReqQuery            []interface{}          `json:"req_query"`
	Status              string                 `json:"status"`
	Uid                 int                    `json:"uid"`
	Type                string                 `json:"type"`
	Username            string                 `json:"username"`
	AddTime             int                    `json:"add_time"`
	UpTime              int                    `json:"up_time"`
}

type YapiApiDetailQueryPath struct {
	Path   string        `json:"path"`
	Params []interface{} `json:"params"`
}
