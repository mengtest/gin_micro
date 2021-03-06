/*
 * @Author: yhlyl
 * @Date: 2019-07-08 23:19:27
 * @LastEditTime: 2019-11-04 21:28:03
 * @LastEditors: yhlyl
 * @Description:
 * @FilePath: /gin_micro/util/uuid/uuid.go
 * @Github: https://github.com/android-coco/gin_micro
 */
/*
 * @Author: yhlyl
 * @Date: 2019-07-08 23:19:27
 * @LastEditTime: 2019-11-04 16:56:30
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /gin_micro/util/uuid/uuid.go
 */
package uuid

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

var uuidCounter uint32 = 0

var machineId = readMachineId()

type UUID string

func readMachineId() []byte {
	var sum [3]byte
	id := sum[:]
	hostname, err1 := os.Hostname()
	if err1 != nil {
		_, err2 := io.ReadFull(rand.Reader, id)
		if err2 != nil {
			panic(fmt.Errorf("cannot get hostname: %v; %v", err1, err2))
		}
		return id
	}
	hw := md5.New()
	hw.Write([]byte(hostname))
	copy(id, hw.Sum(nil))
	//fmt.Println("readMachineId:" + string(id))
	return id
}

// NewUUID returns a new unique UUID.
// 4byte 时间，
// 3byte 机器ID
// 2byte pid
// 3byte 自增ID
func NewUUID() UUID {
	var b [12]byte
	// Timestamp, 4 bytes, big endian
	binary.BigEndian.PutUint32(b[:], uint32(time.Now().Unix()))
	// fmt.Println(b)
	// Machine, first 3 bytes of md5(hostname)
	b[4] = machineId[0]
	b[5] = machineId[1]
	b[6] = machineId[2]
	// Pid, 2 bytes, specs don't specify endianness, but we use big endian.
	pid := os.Getpid()
	b[7] = byte(pid >> 8)
	b[8] = byte(pid)
	// Increment, 3 bytes, big endian
	i := atomic.AddUint32(&uuidCounter, 1)
	b[9] = byte(i >> 16)
	b[10] = byte(i >> 8)
	b[11] = byte(i)
	return UUID(b[:])
}

// Hex returns a hex representation of the UUID.
// 返回16进制对应的字符串
func (id UUID) Hex() string {
	m := md5.Sum([]byte(id))
	return hex.EncodeToString(m[:])[8:24]
}

// Hex returns a hex representation of the UUID.
// 返回大写16进制对应的字符串
func (id UUID) HexToUpper() string {
	m := md5.Sum([]byte(id))
	return strings.ToUpper(hex.EncodeToString(m[:])[8:24])
}

// Hex32
func (id UUID) Hex32() string {
	m := md5.Sum([]byte(id))
	return hex.EncodeToString(m[:])
}
