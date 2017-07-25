package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

var (
	err error
)

// HandlerFunc is a custom implementation of the http.HandlerFunc
type HandlerFunc func(http.ResponseWriter, *http.Request, AppContext)

// makeHandler allows us to pass an environment struct to our handlers, without resorting to global
// variables. It accepts an environment (Env) struct and our own handler function. It returns
// a function of the type http.HandlerFunc so can be passed on to the HandlerFunc in main.go.
func makeHandler(ctx AppContext, fn func(http.ResponseWriter, *http.Request, AppContext)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, ctx)
	}
}

// HealthcheckHandler returns useful info about the app
func HealthcheckHandler(w http.ResponseWriter, req *http.Request, ctx AppContext) {
	check := Healthcheck{
		AppName: "go-rest-api-template",
		Version: ctx.Version,
	}
	ctx.Render.JSON(w, http.StatusOK, check)
}
//

type Result struct {
	Status bool        `json:"status"`
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
}

//
func GetFuncHandler(w http.ResponseWriter, r *http.Request,ctx AppContext) {
	// Vars returns the route variables for the current request, if any.
	//返回请求url中参数 如果有多个 则全部返回
	vars := mux.Vars(r)
	id := vars["id"]
	res := Result{
		Status: true,
		Data:   "get id  is" + id,
		Msg:    "none",
	}
	rc :=ctx.RedisClient.Get()
	defer rc.Close()
	rc.Do("set","getPram",id)
	ctx.Render.JSON(w, http.StatusOK, res)
	//返回json结果集 encode(任意类型参数interface{})
	//json.NewEncoder(w).Encode(res)
	//println("id:",id)

}


func PostFuncHandler(w http.ResponseWriter, r *http.Request,ctx AppContext) {
	//get param from url
	//vars :=mux.Vars(r)
	//id :=vars["id"]
	//get post data from request body
	decoder := json.NewDecoder(r.Body)
	var res interface{}
	err = decoder.Decode(&res)
	if err != nil {
		println(err.Error())
	}
	//ctx.DB.Exec("")
	//return response by json
	ctx.Render.JSON(w, http.StatusOK, res)
}

func PutFuncHandler(w http.ResponseWriter, r *http.Request,ctx AppContext) {
	vars := mux.Vars(r)
	id := vars["id"]
	json.NewEncoder(w).Encode("id is:" + id)
	//处理未知的json类型 使用interface{}来进行解析
	var json_data interface{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&json_data)
	json.NewEncoder(w).Encode(json_data)

}

func DeleteFuncHandler(w http.ResponseWriter, r *http.Request,ctx AppContext) {
	//vars := mux.Vars(r)
	//id := vars["id"]
}

