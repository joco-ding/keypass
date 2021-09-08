package controllers

import (
	"bufio"
	"fmt"
	"keypass/lib/utils"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func AddEntry() {
	_reader := bufio.NewReader(os.Stdin)
	fmt.Print("Nama: ")
	_namatext, _err := _reader.ReadString('\n')
	if _err != nil {
		utils.ErrorLogging("AddEntry1", _err.Error())
		fmt.Println("Gagal")
		return
	}
	_namatext = strings.Trim(_namatext, " \n")
	if len(_namatext) == 0 {
		fmt.Println("Nama harus diisi")
		return
	}
	fmt.Print("Username: ")
	_usertext, _err := _reader.ReadString('\n')
	if _err != nil {
		utils.ErrorLogging("AddEntry2", _err.Error())
		fmt.Println("Gagal")
		return
	}
	_usertext = strings.Trim(_usertext, " \n")
	if len(_usertext) == 0 {
		fmt.Println("Username harus diisi")
		return
	}
	fmt.Print("[A] Auto Generate, [M] Manual: ")
	_pilpwd, _err := _reader.ReadString('\n')
	if _err != nil {
		utils.ErrorLogging("AddEntry2a", _err.Error())
		fmt.Println("Gagal")
		return
	}
	_pilpwd = strings.Trim(_pilpwd, " \n")
	if len(_pilpwd) == 0 {
		fmt.Println("Pilihan harus diisi")
		return
	}
	_spwd := ""
	switch _pilpwd {
	case "a", "A":
		_spwd = utils.GenPass3()
		fmt.Printf("Password: %s\n", _spwd)
	default:
		fmt.Print("Password: ")
		_pwdtext, _err := term.ReadPassword(int(syscall.Stdin))
		if _err != nil {
			utils.ErrorLogging("AddEntry3", _err.Error())
			fmt.Println("Gagal")
			return
		}
		fmt.Println()
		fmt.Print("Ulangi Password: ")
		_repwdtext, _err := term.ReadPassword(int(syscall.Stdin))
		if _err != nil {
			utils.ErrorLogging("AddEntry4", _err.Error())
			fmt.Println("Gagal")
			return
		}
		_spwd = strings.Trim(string(_pwdtext), " \n")
		if len(_spwd) == 0 {
			fmt.Println("Pasword harus diisi")
			return
		}
		_srepwd := strings.Trim(string(_repwdtext), " \n")
		if _spwd != _srepwd {
			fmt.Println("Kedua password tidak sama!")
			return
		}
	}
	_enama := utils.Encrypt(_namatext)
	_euser := utils.Encrypt(_usertext)
	_epwd := utils.Encrypt(_spwd)
	_, _err = utils.InsertDb("INSERT INTO tb_credentials (nama, username, password) VALUES(?, ?, ?)", _enama, _euser, _epwd)
	if _err != nil {
		utils.ErrorLogging("AddEntry5", _err.Error())
		fmt.Println("Gagal")
		return
	}
	fmt.Printf("Berhasil menambahkan data %s %s\n", _namatext, _usertext)

}
