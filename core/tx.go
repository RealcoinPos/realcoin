package core

import (
	"time"
	"fmt"
	"bytes"
	"encoding/gob"
	"strconv"
	"../wallet"
)


type TxData struct {
	Cno			int		`json:"cno"`
	Id			string	`json:"id"`
	Gtype		string	`json:"gtype"`
	Ctype		string	`json:"ctype"`
	Cash		float64 `json:"cash"`
	Apply		float64 `json:"apply"`
	Fee			float64 `json:"fee"`
	Paydate		string	`json:"paydate"`
	Addr		string	`json:"addr"`	
	Toaddr		string	`json:"toaddr"`
	Hashv		string	`json:"hashv"`
	Hashid		string	`json:"hashid"`
	Hashsum		string	`json:"hashsum"`
	Memo		string	`json:"memo"`
	Ipcheck		string	`json:"ipcheck"`
	Regtime		string	`json:"regtime"`
	Udate		string	`json:"udate"`
	State		string	`json:"state"`
}

type TxPool struct {
	Tdata	[]TxData
}


/*
tx 전송
tx 만들기
*/


func (tx TxData) TxInput(w wallet.Wallet,from string, to string,amount float64) TxData {
	txd:=new(TxData)
	txd.Gtype="transfer"
	txd.Cno=int(time.Now().Unix())
	txd.Addr=from
	txd.Toaddr=to
	txd.Cash=amount
	txd.Hashv=setHash(strconv.Itoa(txd.Cno)+from+to+fmt.Sprintf("%.6f",amount))
	Hashsign,_:=w.Sign(StrToByte(txd.Hashv))
	txd.Hashsum=ByteToStr(Hashsign)
	tx=*txd
	return tx
}

func (tx *TxData) Serialize() []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(tx)
	if err != nil {
		fmt.Println(err)
	}
	return buff.Bytes()
}


func Serialize(object interface{}) []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(object)
	if err != nil {
		fmt.Println(err)
	}
	return buff.Bytes()
}

func Deserialize(data []byte) TxPool {
	var transaction TxPool
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&transaction)
	if err != nil {
		fmt.Println(err)
	}
	return transaction
}

func TxDeserialize(data []byte) TxData {
	var transaction TxData
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&transaction)
	if err != nil {
		fmt.Println(err)
	}
	return transaction
}

