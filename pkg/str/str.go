package str

import (
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"strconv"
	"time"
)

func RandString(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

func UidGenerate() string {
	return uuid.New().String()
}

func GenerateRandomID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func StringConvertInt64(value string) int64 {
	num, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		// 错误处理
		fmt.Println("Error:", err)
	}
	return num
}
