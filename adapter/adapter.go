package adapter

import (
	"dfinity_adapter/adapter/util"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"time"
)

type Request struct {
	// Common params:
	TokenType string `json:"token"`
	Price     uint64 `json:"price"`
}

type Result struct {
	data Data `json:"data"`
}

type Data struct {
	To     string  `json:"to"`
	From   string  `json:"from"`
	Result float64 `json:"result"`
}
type Proposal struct {
	Signature []byte `json:"signature"`
	PublicKey []byte `json:"public_key"`
	TokenType string `json:"token_type"`
	Price     uint64 `json:"price"`
	TimeStamp uint64 `json:"time_stamp"`
}

type signatureContent struct {
	TokenType string `json:"token_type"`
	Price     uint64 `json:"price"`
	TimeStamp uint64 `json:"time_stamp"`
}

type dfinityAdaptor struct {
	prikey []byte
	pubkey []byte
	api    *API
}

func NewdfinityAdaptor(url string, port string, priKey, pubKey string) (*dfinityAdaptor, error) {
	strPrikey, _ := hexutil.Decode(priKey)
	strPubkey, _ := hexutil.Decode(pubKey)
	return &dfinityAdaptor{
		prikey: strPrikey,
		pubkey: strPubkey,
		api:    newAPI(url, port),
	}, nil
}

func (adapter *dfinityAdaptor) Handle(req Request) (interface{}, error) {
	if req.TokenType == "" {
		fmt.Errorf("token type can not be empty")
	} else {
		timestamp := uint64(time.Now().Unix())
		signature, err := adapter.sign(signatureContent{
			req.TokenType, req.Price, timestamp,
		})
		if err != nil {
			fmt.Errorf("can not sign, error: %v", err)
			return nil, err
		}
		fmt.Println("Sign message successfully")
		body, err := adapter.api.post(Proposal{
			Signature: signature,
			Price:     req.Price,
			TimeStamp: timestamp,
			TokenType: req.TokenType,
			PublicKey: adapter.pubkey,
		})
		if err != nil {
			fmt.Errorf("Can not post dfinity api: body:%v, error: %v", body, err)
			return body, err
		}
		return body, nil
	}
	return nil, nil
}

func (adapter *dfinityAdaptor) sign(content signatureContent) ([]byte, error) {
	msg, err := json.Marshal(content)
	if err != nil {
		return nil, err
	}
	msgHash := util.Sha256(msg)
	signatiure, err := secp256k1.Sign(msgHash, adapter.prikey)
	if err != nil {
		return nil, err
	}
	fmt.Println("signature:", signatiure)
	return signatiure, nil
}
