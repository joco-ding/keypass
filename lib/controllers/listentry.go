package controllers

import (
	"bufio"
	"fmt"
	"keypass/lib/models"
	"keypass/lib/utils"
	"os"
	"strings"

	"github.com/atotto/clipboard"
)

func ListEntry() {
	_reader := bufio.NewReader(os.Stdin)
	_rows, _err := utils.QueryDb("SELECT id, nama, username FROM tb_credentials WHERE username!='' AND password!=''")
	if _err != nil {
		utils.ErrorLogging("ListEntry1", _err.Error())
		fmt.Println("Gagal")
		return
	}
	for _rows.Next() {
		var _id int64
		var _nama, _username string
		_err = _rows.Scan(&_id, &_nama, &_username)
		if _err != nil {
			utils.ErrorLogging("ListEntry2", _err.Error())
			fmt.Println("Gagal")
			return
		}
		_dnama := utils.Decrypt(_nama)
		_dusername := utils.Decrypt(_username)
		fmt.Printf("%d\t%s\t%s\n", _id, _dnama, _dusername)
	}
	fmt.Print("[S] Selesai, [0-9] Tampilkan ID: ")
	_iddata, _err := _reader.ReadString('\n')
	if _err != nil {
		utils.ErrorLogging("ListEntry3", _err.Error())
		fmt.Println("Gagal")
		return
	}
	_iddata = strings.Trim(_iddata, " \n")
	_panjang := len(_iddata)
	if _panjang == 0 {
		fmt.Println("Gagal")
		return
	}
	switch _iddata {
	case "s", "S":
		fmt.Println("Selesai")
		return
	default:
		_row := utils.QueryRow("SELECT nama, username, password FROM tb_credentials WHERE id=?", _iddata)
		var _nama, _username, _password string
		_err = _row.Scan(&_nama, &_username, &_password)
		if _err != nil {
			utils.ErrorLogging("ListEntry4", _err.Error())
			fmt.Println("Gagal")
			return
		}
		_dnama := utils.Decrypt(_nama)
		_dusername := utils.Decrypt(_username)
		_dassword := utils.Decrypt(_password)
		for {
			fmt.Printf("Nama: %s\nUsername: %s\n[U] Salin Username, [P] Salin Password, [H] Hapus Data, [S] Selesai: ", _dnama, _dusername)
			_pilihan, _err := _reader.ReadString('\n')
			if _err != nil {
				utils.ErrorLogging("ListEntry5", _err.Error())
				fmt.Println("Gagal")
				return
			}
			_pilihan = strings.Trim(_pilihan, " \n")
			_panjang := len(_pilihan)
			if _panjang == 0 {
				fmt.Println("Gagal")
				return
			}
			_selesai := false
			switch _pilihan {
			case "u", "U":
				go models.Menyalin(_dusername)
				_hasilcopy, _err := clipboard.ReadAll()
				if _err != nil {
					utils.ErrorLogging("ListEntry6", _err.Error())
					fmt.Println("Gagal")
					return
				}
				if _hasilcopy == _dusername {
					fmt.Println("USERNAME berhasil disalin")
				}
			case "p", "P":
				go models.Menyalin(_dassword)
				_hasilcopy, _err := clipboard.ReadAll()
				if _err != nil {
					utils.ErrorLogging("ListEntry6", _err.Error())
					fmt.Println("Gagal")
					return
				}
				if _hasilcopy == _dassword {
					fmt.Println("PASSWORD berhasil disalin!")
				}
			case "h", "H":
				_, _err := utils.ExecDb("DELETE FROM tb_credentials WHERE id=?", _iddata)
				if _err != nil {
					utils.ErrorLogging("ListEntry7", _err.Error())
					fmt.Println("Gagal")
					continue
				} else {
					_selesai = true
				}
			default:
				_selesai = true
			}
			if _selesai {
				break
			}
		}
	}
}
