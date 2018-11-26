package tools

import (
	"crypto/sha1"
	"fmt"
	"io"
	"sort"
	"strings"
	"third-gbss/third/wechat/model"
)

func MakeSignature(timestamp, nonce string) string {
	sl := []string{model.Token, timestamp, nonce}
	sort.Strings(sl)
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))
	return fmt.Sprintf("%x", s.Sum(nil))
}