package models

import (
	"database/sql"
	"errors"
	"keypass/lib/utils"
)

func InitDb() error {
	utils.SetPragmaWAL()

	_, _err := utils.ExecDb("CREATE TABLE IF NOT EXISTS 'tb_credentials' ('id' INTEGER, 'nama' TEXT, 'username' TEXT DEFAULT '', 'password' TEXT DEFAULT '', PRIMARY KEY('id' AUTOINCREMENT) )")
	if _err != nil {
		return _err
	}

	_presetnama := "keypass"
	_row := utils.QueryRow("SELECT nama FROM tb_credentials ORDER BY id LIMIT 1")
	var _nama string
	_err = _row.Scan(&_nama)
	if _err != nil {
		if _err == sql.ErrNoRows {
			_encnama := utils.Encrypt(_presetnama)
			_, _err = utils.InsertDb("INSERT INTO tb_credentials (nama) VALUES(?)", _encnama)
			if _err != nil {
				return _err
			}
		} else {
			return _err
		}
	} else {
		_decnama := utils.Decrypt(_nama)
		if _decnama != _presetnama {
			return errors.New("password salah")
		}
	}
	return nil
}
