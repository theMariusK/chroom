package packet

import (
	"crypto/sha256"
	"net"
)

type packet struct {
	data_length []byte
	data        []byte
	hash        []byte
}

func createSize(id int) []byte {
	arr := make([]byte, 2)

	k := 1
	for n := id % 10; id > 0; {
		arr[k] = byte(n)
		k--
		id = id / 10
		n = id % 10
	}

	return arr
}

func ParsePacket(p []byte) (int, string, []byte) {
	var length int
	var msg string

	length = int(p[0])*10 + int(p[1])

	for i := 2; i < length+2; i++ {
		msg += string(p[i])
	}

	return length, msg, p[length+2:]
}

func InitPacket(msg []byte) packet {
	length := createSize(len(msg))
	p := packet{data_length: length, data: msg, hash: GenerateChecksum(msg)}
	return p
}

func GenerateChecksum(msg []byte) []byte {
	hash := sha256.New()
	hash.Write(msg)
	h := hash.Sum(nil)
	return h
}

func SendPacket(conn net.Conn, p packet) {
	buffer := append(p.data_length, p.data...)
	buffer = append(buffer, p.hash...)
	conn.Write(buffer)
}

func CompareChecksum(ch1 []byte, ch2 []byte) bool {
	for i := 0; i < len(ch1); i++ {
		if ch1[i] != ch2[i] {
			return false
		}
	}

	return true
}
