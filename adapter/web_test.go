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
	fmt.Println(req)
	pubkey,_:= hexutil.Decode(req.Arg.PublicKey)
	signature,_ := hexutil.Decode(req.Arg.Signature)
	content := signatureContent{
		TokenType: req.Arg.TokenType,
		Timestamp: req.Arg.Timestamp,
		Price:     req.Arg.Price,
	}

	msg, _ := json.Marshal(content)
	msgHash := util.Sha256(msg)

	verified := secp256k1.VerifySignature(pubkey, msgHash, signature[:64])
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
	seckey,_ := hexutil.Decode("0x65b69e7356c2e8c68f1be482b9b3db9c33196d11c988b3db37ca6953adaf10a8")
	pubkey,_  := hexutil.Decode( "0x049d68bdf6a02aab91f9eb17af2930267007284d8984c90f7bd2a7c54edbee965ce0a53b660b5fc43fe69dc87d2aed1c5eeffe41e7fbc23242bba6685df1143ecb")

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

	fmt.Println("signature", signature)
	fmt.Println("public key:", pubkey)
	//pub,err := secp256k1.RecoverPubkey(msgH, signature)
	//if err != nil{
	//	fmt.Println("can not recover public key From signature,error:",err)
	//}
	//fmt.Printf("The pubkey is :%v \n The pub   is :%v \n",hexutil.Encode(pubkey),hexutil.Encode(pub))
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
