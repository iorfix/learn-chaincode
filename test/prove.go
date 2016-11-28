package main

import (
	"fmt"
	"time"
	"encoding/binary"
)

func main() {

	fmt.Println("Hello, playground")
	id := makeTimestamp()
	fmt.Println(id)
	idByteArr := make([]byte, 4)
	chainuserarray := make([]byte, 0)
	binary.LittleEndian.PutUint32(idByteArr, id)
	fmt.Println("result: ",  idByteArr)
	chainuserarray = append(chainuserarray, idByteArr...)
	fmt.Println("chainuserarray : ",  chainuserarray )

	
	
}

func makeTimestamp() uint32 {
	var now int64 
	now = time.Now().UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))
    return uint32(now)
}