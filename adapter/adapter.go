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

//type Result struct {
//	Data Data `json:"data"`
//}

type Data struct {
	To     string  `json:"to"`
	From   string  `json:"from"`
	Result float64 `json:"result"`
}
type Proposal struct {
	Canister string `json:"canister"`
	Function string `json:"function"`
	Arg      Args   `json:"args"`
}

type Args struct {
	Signature []byte `json:"signature"` // 64 byte
	PublicKey []byte `json:"pub_key"`   //
	TokenType string `json:"token_type"`
	Price     uint64 `json:"price"`
	Timestamp uint64 `json:"timestamp"`
}

type signatureContent struct {
	TokenType string `json:"token_type"`
	Price     uint64 `json:"price"`
	Timestamp uint64 `json:"timestamp"`
}

type dfinityAdaptor struct {
	prikey []byte
	pubkey []byte
	api    *API
}

func NewdfinityAdaptor(url string, priKey, pubKey string) (*dfinityAdaptor, error) {
	strPrikey, _ := hexutil.Decode(priKey)
	strPubkey, _ := hexutil.Decode(pubKey)
	return &dfinityAdaptor{
		prikey: strPrikey,
		pubkey: strPubkey,
		api:    newAPI(url),
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
			Canister: "ywrdt-7aaaa-aaaah-qaaaa-cai",
			Function: "provide_price",
			Arg: Args{
				Signature: signature,
				Price:     req.Price,
				Timestamp: timestamp,
				TokenType: req.TokenType,
				PublicKey: adapter.pubkey,
			},
		})

		fmt.Println("body:", body)

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
	return signatiure, nil
}
