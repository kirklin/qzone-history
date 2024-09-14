package utils

import (
	"regexp"
	"strconv"
)

// GeneratePtqrToken 根据提供的 qrsig 生成 ptqr_token
// GeneratePtqrToken generates a ptqr_token based on the provided qrsig.
func GeneratePtqrToken(qrsig string) string {
	e := 0
	for i := 0; i < len(qrsig); i++ {
		e += (e << 5) + int(qrsig[i])
	}
	return strconv.Itoa(2147483647 & e)
}

// ExtractUin 从cookies中提取QQ号
func ExtractUin(cookies map[string]string) string {
	uin, exists := cookies["uin"]
	if !exists {
		return ""
	}
	re := regexp.MustCompile(`^o0*`)
	return re.ReplaceAllString(uin, "")
}

// GenerateGTK 根据提供的 skey 生成 GTK (g_tk)
// GenerateGTK generates a GTK (g_tk) based on the provided skey.
func GenerateGTK(skey string) string {
	h := 5381
	for i := 0; i < len(skey); i++ {
		h += (h << 5) + int(skey[i])
	}
	return strconv.Itoa(h & 2147483647)
}
