package rid

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"hash/fnv"
	"os"
)

// Salt 生成一个随机的64位无符号整数作为盐值
//
// 返回值:
//
//	uint64: 生成的盐值
func Salt() uint64 {
	hasher := fnv.New64a()
	hasher.Write(ReadMachineID())
	hashValue := hasher.Sum64()
	return hashValue
}

// ReadMachineID 读取机器ID
//
// 返回值为一个3字节的切片，表示机器的唯一标识。
//
// 首先尝试从平台读取机器ID，如果读取失败或读取到的ID为空，则尝试获取系统主机名作为机器ID。
// 如果获取主机名成功且主机名不为空，则使用SHA256哈希算法对主机名进行哈希处理，并将哈希值的前3个字节作为机器ID。
// 如果获取主机名失败或主机名为空，则生成一个随机的3字节切片作为机器ID。
//
// 如果在生成随机机器ID时发生错误，则会导致程序崩溃并输出错误信息。
func ReadMachineID() []byte {
	id := make([]byte, 3)
	machineID, err := readPlatformMachineID()
	if err != nil || len(machineID) == 0 {
		machineID, err = os.Hostname()
	}

	if err == nil && len(machineID) != 0 {
		hasher := sha256.New()
		hasher.Write([]byte(machineID))
		copy(id, hasher.Sum(nil))
	} else {
		if _, randErr := rand.Reader.Read(id); randErr != nil {
			panic(fmt.Errorf("id: cannot get hostname nor generate a random number: %w; %w", err, randErr))
		}
	}
	return id
}

// readPlatformMachineID 读取系统机器ID
// 首先尝试从 /etc/machine-id 文件中读取机器ID，如果文件不存在或内容为空，
// 则尝试从 /sys/class/dmi/id/product_uuid 文件中读取机器ID。
//
// 参数：
//
//	无
//
// 返回值：
//
//	返回的字符串为机器ID
//	如果读取过程中发生错误，返回错误信息
//	如果读取成功但文件内容为空，返回空字符串和 nil
func readPlatformMachineID() (string, error) {
	data, err := os.ReadFile("/etc/machine-id")
	if err != nil || len(data) == 0 {
		data, err = os.ReadFile("/sys/class/dmi/id/product_uuid")
	}
	return string(data), err
}
