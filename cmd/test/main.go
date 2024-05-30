package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/olafszymanski/int-sdk/integration/pb"
	"google.golang.org/protobuf/proto"
)

func main() {
	// t := time.Now()
	f, err := os.Open("test_data/output/PRE_MATCH_BASKETBALL1.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var events []*pb.Event
	dec := json.NewDecoder(f)
	if err := dec.Decode(&events); err != nil {
		log.Fatal(err)
	}
	// t1 := time.Now()
	wg := &sync.WaitGroup{}
	wg.Add(len(events))
	lock := &sync.Mutex{}
	hashMap := make(map[string][]byte)
	for _, e := range events {
		go func(e *pb.Event) {
			defer wg.Done()

			b, _ := proto.Marshal(e)

			lock.Lock()
			fmt.Println(*e.ExternalId)
			hashMap[*e.ExternalId] = b
			lock.Unlock()
		}(e)
	}
	wg.Wait()
	// fmt.Println(time.Since(t1), len(events), len(hashMap), size.Of(hashMap))
}
