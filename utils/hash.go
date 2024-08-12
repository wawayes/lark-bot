package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// CalculateHash 计算给定接口的 SHA-256 哈希值
func CalculateHash(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}
