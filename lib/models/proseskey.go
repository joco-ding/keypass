package models

import (
	"encoding/hex"
	"errors"
	"keypass/lib/stores"
	"strings"
)

func ProsesKey(teks string) error {
	teks = strings.Trim(teks, " \n")
	_panjangpwd := len(teks)
	if _panjangpwd == 0 {
		return errors.New("password harus diisi")
	}
	if _panjangpwd < 32 {
		_sisa := 32 - _panjangpwd
		teks = teks + strings.Repeat("0", _sisa)
	} else if _panjangpwd > 32 {
		_runes := []rune(teks)
		teks = string(_runes[0:32])
	}
	bytes := []byte(teks)
	tocopy := hex.EncodeToString(bytes)
	stores.Config.KeyString = tocopy
	return nil
}
