package models

import (
	"keypass/lib/stores"
	"time"

	"github.com/atotto/clipboard"
)

func Menyalin(ygdisalin string) {
	stores.Config.Menyalin = true
	clipboard.WriteAll(ygdisalin)
	stores.Config.Expired = time.Now().Unix() + 60
	_firsttime := true
	for stores.Config.Expired > time.Now().Unix() {
		if _firsttime {
			_firsttime = false
			stores.Config.Menyalin = false
		}
		if stores.Config.Menyalin {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
	if stores.Config.Expired > time.Now().Unix() {
		clipboard.WriteAll("")
	}
}
