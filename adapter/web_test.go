package adapter

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"dfinity_adapter/adapter/util"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestMockDfinityServer(t *testing.T) {
	router := gin.Default()
	router.POST("/", dfinityServerCall)

	router.Run(":2334")
}

func dfinityServerCall(c *gin.Context) {
	var req Proposal
	if err := c.BindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, dfinityResp{
			StatusCode: http.StatusInternalServerError,
			Status:     "false",
			Data:       req,
			Error:      "Can not decode the json data",
		})
		return
	}

	pubkey := req.Arg.PublicKey
	content := signatureContent{
		TokenType: req.Arg.TokenType,
		Timestamp: req.Arg.Timestamp,
		Price:     req.Arg.Price,
	}

	msg, _ := json.Marshal(content)
	msgHash := util.Sha256(msg)

	verified := secp256k1.VerifySignature(pubkey, msgHash, req.Arg.Signature[:64])
	if !verified {
		log.Println("Can not verify the signature.")
		c.JSON(http.StatusInternalServerError, dfinityResp{
			StatusCode: http.StatusInternalServerError,
			Status:     "false",
			Data:       req,
			Error:      "Can not verify the signature",
		})
		return
	}
	c.JSON(http.StatusOK, dfinityResp{
		StatusCode: http.StatusOK,
		Status:     "success",
		Data:       req,
	})
	fmt.Println("Dfinity api server verify the Result successfully!")
}

type dfinityResp struct {
	StatusCode int
	Status     string
	Error      string
	Data       Proposal
}

//func TestLinkNode(t *testing.T) {
//	pro := JobReq{
//		JobID: "link",
//		Data: Request{
//			TokenType: "ETH",
//			Price:     10000,
//		},
//	}
//
//	client := &http.Client{}
//	bytesData, _ := json.Marshal(pro)
//	req, _ := http.NewRequest("POST", "http://localhost:2333/", bytes.NewReader(bytesData))
//	resp, err := client.Do(req)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	//body, err := ioutil.ReadAll(resp.Body)
//	//if err != nil {
//	//	fmt.Println("error:", err)
//	//}
//	//fmt.Println(string(body))
//	if resp.StatusCode != http.StatusOK {
//		//fmt.Println("body:", string(body))
//		//fmt.Println("error:", err)
//		fmt.Println("link fail To send price")
//		return
//	}
//	fmt.Println(resp)
//	defer resp.Body.Close() //一定要关闭resp.Body
//	//fmt.Println("successful: body:", string(body))
//}

func TestSignature(t *testing.T) {
	pubkey, seckey := generateKeyPair()

	content := signatureContent{
		TokenType: "ETH",
		Price:     10000,
		Timestamp: 0,
	}
	msg, _ := json.Marshal(content)
	fmt.Println("msg:", hexutil.Encode(msg))
	msgH := util.Sha256(msg)
	fmt.Println("msg hash:", hexutil.Encode(msgH))
	signature, _ := secp256k1.Sign(msgH, seckey)

	fmt.Println("signature", hexutil.Encode(signature))
	fmt.Println("public key:", hexutil.Encode(pubkey))
	//pub,err := secp256k1.RecoverPubkey(msgH, signature)
	//if err != nil{
	//	fmt.Println("can not recover public key From signature,error:",err)
	//}
	//fmt.Printf("The pubkey is :%v \n The pub   is :%v \n",hexutil.Encode(pubkey),hexutil.Encode(pub))
	fmt.Printf("private key len: %v\n public key len: %v\n sig len: %v\n", len(seckey), len(pubkey), len(signature))
	verified := secp256k1.VerifySignature(pubkey, msgH, signature[:64])
	fmt.Println("past the verify?", verified)
}

func generateKeyPair() (pubkey, privkey []byte) {
	key, err := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	pubkey = elliptic.Marshal(secp256k1.S256(), key.X, key.Y)
	privkey = make([]byte, 32)
	blob := key.D.Bytes()
	copy(privkey[32-len(blob):], blob)

	return pubkey, privkey
}

func TestTime(t *testing.T) {
	fmt.Println(time.Now().Unix())
	fmt.Println(uint64(time.Now().Unix()))
}

func TestMockAdapter(t *testing.T) {
	router := gin.Default()
	router.POST("", mockAdaptorCall)
	router.Run(":2333")
}

func mockAdaptorCall(c *gin.Context) {
	var req testRequest
	//resp := c.Request
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(string(body))
	//json.Unmarshal([]byte(string(body)),req)
	//fmt.Println(req)
	if err := c.BindJSON(&req); err != nil {
		log.Println("invalid JSON payload", err)
		errorJob(c, http.StatusBadRequest, req.JobId, "Invalid JSON payload")
		return
	}
	fmt.Println(req)
}

type testRequest struct {
	JobId  string `json:"jobId"`
	Id     string `json:"id"`
	Result testResult `json:"result"`
}

type testResult struct {
	Data testData `json:"data"`
}

type testData struct {
	To     string  `json:"to"`
	From   string  `json:"from"`
	Result float64 `json:"result"`
}
