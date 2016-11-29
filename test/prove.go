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
	
	id = makeTimestamp() + 250
	binary.LittleEndian.PutUint32(idByteArr, id)
	chainuserarray = append(chainuserarray, idByteArr...)
	
	fmt.Println("chainuserarray : ",  chainuserarray )

	numElems := len(chainuserarray)/4
	fmt.Println("len: ",  numElems )
	chainuserint := make([]uint32, numElems)
	fmt.Println("init:" , chainuserint)
	for i := 0; i < numElems; i++ {
		elemByte := chainuserarray[i*4:i*4+4]
		fmt.Println("elemByte : ",  elemByte )	
		val := binary.LittleEndian.Uint32(elemByte)
		chainuserint[i] = val
	}
	fmt.Println(chainuserint)

	
}

func makeTimestamp() uint32 {
	var now int64 
	now = time.Now().UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))
    return uint32(now)
}