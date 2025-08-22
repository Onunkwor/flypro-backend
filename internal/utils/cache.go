package utils

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
)

func MakeCacheKey(filters map[string]interface{}, offset, limit int) string {
	b, _ := json.Marshal(filters)
	hash := fmt.Sprintf("%x", sha1.Sum(b))
	return fmt.Sprintf("expenses:%s:%d:%d", hash, offset, limit)
}
