package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"

	distkvs "example.org/cpsc416/a5"
	"example.org/cpsc416/a5/kvslib"
)

type StorageAddr string

type FrontEndConfig struct {
	ClientAPIListenAddr 	string
	StorageAPIListenAddr 	string
	Storage 							StorageAddr
}

type FrontEndStorageFailed struct{}

type FrontEndStorageStarted struct{}

type FrontEndGet struct {
	Key string
}

type FrontEndGetResult struct {
	Key   string
	Value *string
	Err   bool
}

type StorageGet struct {
	Key string
}

type StorageGetResult struct {
	Key   string
	Value *string
}

type FrontEnd struct {
	// state may go here
	rpcClient *rpc.Client
}

func main() {
	var config distkvs.FrontEndConfig
	err := distkvs.ReadJSONConfig("config/frontend_config.json", &config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(config)

	frontend := FrontEnd{}
	err = frontend.Start(config.ClientAPIListenAddr, config.StorageAPIListenAddr, 0)

	if err != nil {
		log.Fatal(err)
	}
}

func (f *FrontEnd) Start(clientAPIListenAddr string, storageAPIListenAddr string, storageTimeout uint8) error {
	// result := new(FrontEndGetResult)
	result := new(FrontEnd)
	rpc.Register(result)

	// result2 := new(FrontEndPutResult)
	// rpc.Register(result2)

	rpc.HandleHTTP()

	rpcClient, rpcErr := rpc.DialHTTP("tcp", storageAPIListenAddr)
	if rpcErr != nil {
		return errors.New("Cannot established connection with Storage server.")
	}

	result.rpcClient = rpcClient

	l, e := net.Listen("tcp", clientAPIListenAddr)
	if e != nil {
		log.Fatal("listen error:", e)
	}

	err := http.Serve(l, nil)
	if err != nil {
		log.Fatal("listen error:", err)
	}

	return errors.New("Frontend Fail")
}

func (f *FrontEnd) HandleGet(args kvslib.KvslibGet, reply *kvslib.ResultStruct) error {
	storageArgs := StorageGet{args.Key}
	storageReply := new(StorageGetResult)
	funcCall := f.rpcClient.Go("Storage.StorageGet", storageArgs, &storageReply, nil)
	replyCall := <-funcCall.Done

	// log.Print(replyCall.Error)
	// log.Print(*storageReply.Value)

	if replyCall.Error != nil {
		return errors.New("FE to Strage fail")
	}

	ret := kvslib.ResultStruct{}
	ret.Result = storageReply.Value
	*reply = ret

	log.Printf("OpId: %d Get value %s from %s", args.OpId, *reply.Result, args.Key)

	return nil
}
