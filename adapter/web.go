package adapter

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Handler func(req Request) (interface{}, error)

func RunWebserver(
	handler Handler,
	port string,
) {
	if port == "" {
		port = "8080"
	}
	srv := NewHTTPService(handler)
	err := srv.Router.Run(fmt.Sprintf(":%v", port))
	if err != nil {
		fmt.Println(err)
	}
}

type HttpService struct {
	Router  *gin.Engine
	Handler Handler
}

func NewHTTPService(
	handler Handler,
) *HttpService {
	srv := HttpService{
		Router:  gin.Default(),
		Handler: handler,
	}
	srv.createRouter()
	return &srv
}


func (srv *HttpService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.Router.ServeHTTP(w, r)
}

func (srv *HttpService) createRouter() {
	r := gin.Default()
	r.POST("/", srv.Call)

	srv.Router = r
}

type JobReq struct {
	Id     string `json:"id"`
	JobId  string `json:"jobId"`
	Result Result `json:"result"`
}

type resp struct {
	JobRunID   string      `json:"jobRunID"`
	StatusCode int         `json:"status_code"`
	Status     string      `json:"status"`
	Data       interface{} `json:"data"`
	Error      interface{} `json:"error"`
}

func errorJob(c *gin.Context, statusCode int, jobId, error string) {
	c.JSON(statusCode, resp{
		JobRunID:   jobId,
		StatusCode: statusCode,
		Status:     "errored",
		Error:      error,
	})
}

func (srv *HttpService) Call(c *gin.Context) {
	var req JobReq
	fmt.Printf("Get In adapter call\n")
	if err := c.BindJSON(&req); err != nil {
		log.Println("invalid JSON payload",err)
		errorJob(c, http.StatusBadRequest, req.JobId, "Invalid JSON payload")
		return
	}
	fmt.Println("id:",req.Id)
	fmt.Println("JobId:",req.JobId)
	fmt.Println("To:",req.Result.data.To)
	fmt.Println("From:",req.Result.data.From)
	fmt.Println("price",req.Result.data.Result)


	//body, err := ioutil.ReadAll(c.Request.Body)
    //strBody := string(body)
    //fmt.Println(strBody)
	//if err := json.Unmarshal([]byte(strBody), &req); err == nil {
	//	fmt.Println("id:",req.Id)
	//	fmt.Println("JobId:",req.JobId)
	//	fmt.Println("data:",req.Result.data)
	//} else {
	//	fmt.Println(err)
	//}


//
//	println(c.Request.Header)
//	println(string(body))

	//if err := validateRequest(&req); err != nil {
	//	log.Println(err)
	//	errorJob(c, http.StatusBadRequest, req.JobID, err.Error())
	//	return
	//}
	fmt.Printf("Get data From link node %+v\n", req.Id)
	res, err := srv.Handler(Request{
		TokenType: req.Result.data.From,
		Price: uint64(req.Result.data.Result),
	})

	if err != nil {
		log.Println(err)
		errorJob(c, http.StatusInternalServerError, req.JobId, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp{
		JobRunID:   req.JobId,
		StatusCode: http.StatusOK,
		Status:     "success",
		Data:       res,
	})
}
