package dbconn

import (
	"crypto/sha1"
	"encoding/hex"
)

func getPassHex(salt string, passwd string) string {
	o := sha1.New()
	o.Write([]byte(salt))
	o.Write([]byte(passwd))
	return hex.EncodeToString(o.Sum(nil))
}
