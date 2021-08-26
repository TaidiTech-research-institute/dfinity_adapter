package adapter

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
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
	Result Data `json:"data"`
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
		errorJob(c, http.StatusBadRequest, req.Id, "Invalid JSON payload")
		return
	}
	fmt.Println("id:",req.Id)
	fmt.Println("To:",req.Result.To)
	fmt.Println("From:",req.Result.From)
	fmt.Println("price",req.Result.Result)
	fmt.Println("time:",time.Now())

	go func() {
			srv.Handler(Request{
			TokenType: req.Result.From,
			Price: uint64(req.Result.Result*10000),
		})
	}()

	//res, err := srv.Handler(Request{
	//	TokenType: req.Result.From,
	//	Price: uint64(req.Result.Result*10000),
	//})
	//
	//if err != nil {
	//	log.Println(err)
	//	errorJob(c, http.StatusInternalServerError, req.Id, err.Error())
	//	return
	//}
	c.JSON(http.StatusOK, resp{
		JobRunID:   req.Id,
		StatusCode: http.StatusOK,
		Status:     "success",
		Data:       req,
	})
	fmt.Println("Deal with the response successfully!")
}
