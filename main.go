package main

import (
	"fmt"
	 "./config"
	 "./core"
	"net/http"
	"flag"
)

func main() {
	Configure:=config.SetConfiguration("./config/realcoin.conf")
	core.Configure=Configure
    serverPort := flag.String("port","5000","Http port number where server will run")
    flag.Parse()
    blockchain := core.GenesisChain()
	nodeID := "POSx32De9e2r3jAp09zC"
    fmt.Println("Starting POS HTTP Server. Listening at port %q", *serverPort)
    http.Handle("/", core.NewHandler(blockchain, nodeID))
    http.ListenAndServe(fmt.Sprintf(":%s", *serverPort), nil)
}


