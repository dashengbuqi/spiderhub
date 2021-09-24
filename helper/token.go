package helper

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
)

type Token []byte

func NewToken(user_id int, app_id primitive.ObjectID, debug_id int) Token {
	var t string
	t = fmt.Sprintf("%s%s%s", strconv.Itoa(user_id), app_id.String(), strconv.Itoa(debug_id))
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
