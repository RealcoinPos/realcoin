package core

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/gob"
	"runtime"
	"fmt"
	"os"
	"strconv"
	"strings"
    "path/filepath"
)


func Err(err error, exit int) int {
	if err != nil {
		fmt.Println(err)	
	}
	if exit>=1 {
		os.Exit(exit)
		return 1
	}
	return 0
}

func StrToByte(str string) []byte {
	sb := make([]byte, len(str))
	for k, v := range str {
		sb[k] = byte(v)
	}
	return sb[:]
}

func ByteToStr(bytes []byte) string {
	var str []byte
	for _, v := range bytes {
		if v != 0x0 {
			str = append(str, v)
		}
	}
	return fmt.Sprintf("%s", str)
}


func vaildCubeno(cno int) bool {
	result:=true
	if cno<1 || cno>27 {
		result=false	
	} 
	return result
}

func setHash(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func sethash2(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}


func ArgServer() {
    fmt.Println(os.Args)
}




func fileWrite(path string, object interface{}) error {
	datapath:="rblock"+ string(filepath.Separator) 
	file, err := os.Create(datapath+path)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

func fileRead(filename string, object interface{}) error {
	datapath:="rblock"+ string(filepath.Separator) 
	file, err := os.Open(datapath+filename)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}

func fileCheck(e error) {
	if e != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Println(line, "\t", file, "\n", e)
		os.Exit(1)
	}
}

func blockName(find string) string {
	result:=""
    dirname := "." + string(filepath.Separator) + "rblock"
    d, err := os.Open(dirname)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer d.Close()
    fi, err := d.Readdir(-1)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    for _, fi := range fi {
        if fi.Mode().IsRegular() {
            fstr:=fi.Name()
			fstr=fstr[0:len(find)]
			if fstr==find {
				result=fi.Name()
				return result
			}
			//fmt.Println(fi.Name(), fi.Size(), "bytes", fstr ," ", find)
        }
    }
	return result
}

func blockFinder(find string) bool {
	result:=false
    dirname := "." + string(filepath.Separator) + "rblock"
    d, err := os.Open(dirname)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer d.Close()
    fi, err := d.Readdir(-1)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    for _, fi := range fi {
        if fi.Mode().IsRegular() {
            fstr:=fi.Name()
			fstr=fstr[0:len(find)]
			if fstr==find {
				result=true
				return result
			}
			//fmt.Println(fi.Name(), fi.Size(), "bytes", fstr ," ", find)
        }
    }
	return result
}


func blockRead(index int,object interface{}) error {
	err:=fileRead(blockName(strconv.Itoa(index)+"_"),object)
	return err;
}

func blockFile(block Block) error {
	filename:=blockName(strconv.Itoa(block.Index)+"_")
	err := fileWrite(filename, block)
	return err
}


func CurrentHeight() int {
	result:=0
    dirname := "." + string(filepath.Separator) + "rblock"
    d, err := os.Open(dirname)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer d.Close()
    fi, err := d.Readdir(-1)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
	max:=0
    for _, fi := range fi {
        if fi.Mode().IsRegular() {
			fstr:=strings.Split(fi.Name(),"_")
			fmt.Println(fi.Name(), fi.Size(), "bytes", fstr[0] ," ", fstr[1]," ", fstr[2])
			cstr:=fstr[0]
			cno,_:=strconv.Atoi(cstr)
			if cno>max {
				max=cno
			}
			result=max
        }
    }
	result++
	return result
}

