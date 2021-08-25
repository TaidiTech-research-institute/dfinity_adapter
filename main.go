package main

import (
	"dfinity_adapter/adapter"
	"fmt"
)

func main() {
	fmt.Println("Starting Substrate adapter")

	//privkey := os.Getenv("SA_PRIVATE_KEY")
	//txType := os.Getenv("SA_TX_TYPE")
	//endpoint := os.Getenv("SA_ENDPOINT")
	//port := os.Getenv("SA_PORT")

	//adapterClient, err := adapter.NewSubstrateAdapter(privkey, txType, endpoint)
	//if err != nil {
	//	fmt.Println("Failed starting Substrate adapter:", err)
	//	return
	//}

	//adapter.RunWebserver(adapterClient.Handle, port)

	privkey := "0x65b69e7356c2e8c68f1be482b9b3db9c33196d11c988b3db37ca6953adaf10a8"
	pubKey  := "0x049d68bdf6a02aab91f9eb17af2930267007284d8984c90f7bd2a7c54edbee965ce0a53b660b5fc43fe69dc87d2aed1c5eeffe41e7fbc23242bba6685df1143ecb"
	endpoint := "http://101.132.161.57:3000/v1/update"
	localPort := "2333"

	adapterClient, err := adapter.NewdfinityAdaptor(endpoint, privkey,pubKey)
	if err != nil {
		panic(err)
	}
	adapter.RunWebserver(adapterClient.Handle,localPort)
}

