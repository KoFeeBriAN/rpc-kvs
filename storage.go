package distkvs

import (
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"

	"github.com/DistributedClocks/tracing"
)

type StorageConfig struct {
	StorageID        string
	StorageAdd       StorageAddr
	ListenAddr       string
	FrontEndAddr     string
	DiskPath         string
	TracerServerAddr string
	TracerSecret     []byte
}

type StorageLoadSuccess struct {
	State map[string]string
}

type StoragePut struct {
	Key   string
	Value string
	delay int
}

type StorageSaveData struct {
	Key   string
	Value string
}

type StorageGet struct {
	Key string
}

type StorageGetResult struct {
	Key   string
	Value *string
}

type Storage struct {
	// state may go here
}

var database = make(map[string]string)

func (*Storage) Start(frontEndAddr string, storageAddr string, diskPath string, strace *tracing.Tracer) error {

	result := new(Storage)
	rpc.Register(result)

	rpc.HandleHTTP()

	if _, err := os.Stat("mem"); os.IsNotExist(err) {
		_, err := os.Create("mem")
		if err != nil {
			panic(err)
		}
	}
	data := make([]byte, 100000)
	file, err := os.Open("mem")
	if err != nil {
		log.Println(err)
	}
	count, err := file.Read(data)
	if err != nil {
		if err == io.EOF {
			log.Println(err)
		} else {
			log.Fatal(err)
		}
	}

	var key string
	var val string
	i := 0
	s := 0

	for i < count {
		if s == 1 {
			if data[i] == '\n' {
				s = 0
				database[key] = val
				key = ""
				val = ""
				i++
			} else {
				val = val + string(data[i])
			}
		}
		if s == 0 {
			if data[i] == ';' {
				s++
			} else {
				key = key + string(data[i])
			}
		}
		i++
	}

	if s != 0 {
		database[key] = val
	}

	database["k99"] = "delay"
	// log.Print(database["key1"])
	PrintDB()

	l, e := net.Listen("tcp", frontEndAddr)
	if e != nil {
		log.Fatal("listen error:", e)
	}

	err = http.Serve(l, nil)
	if err != nil {
		log.Fatal("listen error:", err)
	}

	return errors.New("not implemented")
}

func (*Storage) StorageGet(args StorageGet, reply *StorageGetResult) error {
	ret := StorageGetResult{}
	s := ""
	_, prs := database[args.Key]
	if !prs {
		s = "key not found"
	} else {
		s = database[args.Key]
	}

	// s := args.Key + "Value"
	ret.Value = &s
	*reply = ret
	log.Printf("Get value %s from %s", *reply.Value, args.Key)
	// log.Print(database["k0"])
	// PrintDB()
	return nil
}

func (*Storage) StoragePut(args StoragePut, reply *string) error {
	log.Printf("Put value %s to %s", args.Value, args.Key)
	s := "success"

	if args.Key == "k99" {
		log.Print("delay for 5 second")
		time.Sleep(5 * time.Second)
	} else {
		database[args.Key] = args.Value

		if _, err := os.Stat("mem"); os.IsNotExist(err) {
			_, err := os.Create("mem")
			if err != nil {
				panic(err)
			}
		}
		file, err := os.OpenFile("mem", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
			s = "cannot open file"
		}
		defer file.Close()
		if _, err := file.WriteString(args.Key + ";" + args.Value + "\n"); err != nil {
			log.Fatal(err)
			s = "cannot write"
		}
	}

	// sPtr := new(string)
	// sPtr = &s
	*reply = s

	// log.Print(database[args.Key])
	// PrintDB()
	return nil
}

func PrintDB() {
	for index, element := range database {
		log.Println(index, "=>", element)
	}
}
