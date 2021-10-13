package helper

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

type Token []byte

func NewToken(user_id int64, app_id int64, debug_id int64) Token {
	var t string
	t = fmt.Sprintf("%d%d%d", user_id, app_id, debug_id)
	return []byte(t)
}

func (this Token) Crawler() Token {
	str := fmt.Sprintf("%s%s", "crawler", string(this))
	return []byte(str)
}

func (this Token) Clean() Token {
	str := fmt.Sprintf("%s%s", "clean", string(this))
	return []byte(str)
}

func (this Token) Pool() Token {
	str := fmt.Sprintf("%s%s", "pool", string(this))
	return []byte(str)
}

func (this Token) ToString() string {
	h := md5.New()
	h.Write(this)
	return hex.EncodeToString(h.Sum(nil))
}
