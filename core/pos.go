package core

import (
	"fmt"
	"os"
	"net"
	"net/http"
	"io"
	"io/ioutil"
	"strings"
	"strconv"
	"encoding/json"
	"../config"
)

const protocol = "tcp"
type PosNode struct {
	Ip		string
	Port	string
	Etc		string
}

type POS struct {
	Index		int
	Hash		string
	Type		int
	Address		string
	Amount		int
	State		int
}

type PosWallet struct {
	Address		string
	Amount		int
}

type Handler struct {
	Chain	   []Block
	nodeId     string
}

type Response struct {
	value      interface{}
	statusCode int
	err        error
}

var nodes []PosNode
var cnodes []string
var Configure config.Configuration

func IpCheck() []string {
	host, err := os.Hostname()
	if err != nil {
		return nil
	}
	addrs, err := net.LookupHost(host)
	if err != nil {
		return nil
	}
	addrs=append(addrs,host)
	return addrs
}

func netError(err error) {
	if err != nil && err != io.EOF {
		fmt.Println("Network Error : ", err)
	}
}

func NewHandler(Chain []Block, nodeID string) http.Handler {
	h := Handler{Chain,nodeID}
	mux := http.NewServeMux()
	mux.HandleFunc("/height", PosResponse(h.GetHeight))
	mux.HandleFunc("/node", PosResponse(h.RegisterNode))
	mux.HandleFunc("/transaction", PosResponse(h.AddTransaction))
	mux.HandleFunc("/chain", PosResponse(h.GetBlockchain))
	mux.HandleFunc("/genesis", PosResponse(h.Genesis))
	return mux
}

func PosCheck(mode string) string {
	arr:=IpCheck()
	reader :=strings.NewReader("cmode="+mode+"&mac="+arr[0]+"&ip="+arr[1]+"&hostname="+arr[2]+"&netname="+Configure.Netname+"&netset="+Configure.Netset+"&netport="+strconv.Itoa(Configure.Netport)+"&blockperiod="+strconv.Itoa(Configure.Blockperiod)+"&maxnumber="+strconv.Itoa(Configure.Maxnumber))
	request, _ := http.NewRequest("POST", "http://"+Configure.Masterserver+"/_bchain/rcoin/"+mode, reader)
	request.Header.Add("content-type", "application/x-www-form-urlencoded")
	request.Header.Add("cache-control", "no-cache")
	client := &http.Client{}
	res, _ := client.Do(request)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return string(body)
}

func PosResponse(h func(io.Writer, *http.Request) Response) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := h(w, r)
		msg := resp.value
		if resp.err != nil {
			msg = resp.err.Error()
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.statusCode)
		if err := json.NewEncoder(w).Encode(msg); err != nil {
			fmt.Printf("could not encode Response to output: %v", err)
		}
		//PosCheck("Http")
	}
}

func (h *Handler) AddTransaction(w io.Writer, r *http.Request) Response {
	if r.Method != http.MethodPost {
		return Response{
			nil,
			http.StatusMethodNotAllowed,
			fmt.Errorf("method %s not allowd", r.Method),
		}
	}
	fmt.Printf("Transaction data to the Chain...\n")
	var tx TxData
	err := json.NewDecoder(r.Body).Decode(&tx)
	index := CurrentHeight()
	resp := map[string]string{
		"message": fmt.Sprintf("Transaction will be added to Block %d", index),
	}
	status := http.StatusCreated
	if err != nil {
		status = http.StatusInternalServerError
		fmt.Printf("There was an error when trying to add a transaction %v\n", err)
	}
	return Response{resp, status, err}
}

func (h *Handler) RegisterNode(w io.Writer, r *http.Request) Response {
	if r.Method != http.MethodPost {
		return Response{
			nil,
			http.StatusMethodNotAllowed,
			fmt.Errorf("method %s not allowd", r.Method),
		}
	}
	fmt.Println("Adding node to the Chain")

	var body map[string][]string
	err := json.NewDecoder(r.Body).Decode(&body)
	for _, node := range body["nodes"] {
		cnodes=append(cnodes,node)
	}
	resp := map[string]interface{}{
		"message": "New nodes have been added",
		"nodes":   cnodes,
	}
	status := http.StatusCreated
	if err != nil {
		status = http.StatusInternalServerError
		err = fmt.Errorf("fail to register nodes")
		fmt.Printf("there was an error when trying to register a new node %v\n", err)
	}
	return Response{resp, status, err}
}


func (h *Handler) GetHeight(w io.Writer, r *http.Request) Response {
	if r.Method != http.MethodGet {
		return Response{
			nil,
			http.StatusMethodNotAllowed,
			fmt.Errorf("method %s not allowd", r.Method),
		}
	}
	fmt.Println("GetInfo requested")
	resp := map[string]interface{}{"length": len(h.Chain)}
	return Response{resp, http.StatusOK, nil}
}

func (h *Handler) GetBlockchain(w io.Writer, r *http.Request) Response {
	fmt.Println("Blockchain requested")
	resp := map[string]interface{}{"blockchain": h.Chain, "height": len(h.Chain)}
	return Response{resp, http.StatusOK, nil}
}

func (h *Handler) Genesis(w io.Writer, r *http.Request) Response {
	fmt.Println("Genesis requested")
	c:=CurrentHeight()
	if c>1 {
		h.Chain=append(h.Chain,genesisBlock("Genesis data"))
	}
	resp := map[string]interface{}{"blockchain": h.Chain, "height": len(h.Chain), "genesis":"genesis file"}
	return Response{resp, http.StatusOK, nil}
}


