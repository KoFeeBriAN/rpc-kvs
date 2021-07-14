package distkvs

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

type StorageAddr string

// this matches the config file format in config/frontend_config.json
type FrontEndConfig struct {
	ClientAPIListenAddr  string
	StorageAPIListenAddr string
	Storage              StorageAddr
	TracerServerAddr     string
	TracerSecret         []byte
}

type FrontEndStorageStarted struct{}

type FrontEndStorageFailed struct{}

type FrontEndPut struct {
	Key   string
	Value string
}

type FrontEndPutResult struct {
	Err bool
}

type FrontEndGet struct {
	Key string
}

type FrontEndGetResult struct {
	Key   string
	Value *string
	Err   bool
}

type FrontEnd struct {
	// state may go here
	rpcClient *rpc.Client
}

type StoragePut struct {
	Key   string
	Value string
	delay int
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
	// //We will handle GET, PUT in herere???
	// c := kvslib.ResultStruct{}
	// log.Printf("%s", *c.Result)

	return errors.New("Frontend Fail")
}

func (f *FrontEnd) HandlePut(args kvslib.KvslibPut, reply *kvslib.ResultStruct) error {
	storageArgs := StoragePut{args.Key, args.Value, 0}
	storageReply := new(string)
	funcCall := f.rpcClient.Go("Storage.StoragePut", storageArgs, &storageReply, nil)
	replyCall := <-funcCall.Done

	if replyCall.Error != nil {
		return errors.New("FE to Strage fail")
	}
	// log.Print(*storageReply)

	ret := kvslib.ResultStruct{}
	ret.Result = storageReply
	*reply = ret

	log.Printf("OpId: %d Put value %s to %s", args.OpId, args.Value, args.Key)
	// log.Print(*reply.Result)
	return nil
}
