package controllers

import (
	"bufio"
	"fmt"
	"keypass/lib/models"
	"keypass/lib/stores"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func RunApp() {
	_interrupt := make(chan os.Signal, 1)
	signal.Notify(_interrupt, os.Interrupt)
	fmt.Println("CTRL+C to stop")
	go func() {
		_reader := bufio.NewReader(os.Stdin)
		fmt.Print("Password: ")
		_readpwd, _err := term.ReadPassword(int(syscall.Stdin))
		if _err != nil {
			panic(_err)
		}
		fmt.Println()
		_textread := string(_readpwd)
		_err = models.ProsesKey(_textread)
		if _err != nil {
			panic(_err)
		}
		fmt.Print("DB Location: ")
		_dbtextread, _err := _reader.ReadString('\n')
		if _err != nil {
			panic(_err)
		}
		_dbtextread = strings.Trim(_dbtextread, " \n")
		_panjangdb := len(_dbtextread)
		if _panjangdb == 0 {
			panic("password harus diisi")
		}
		stores.Config.DbPath = _dbtextread
		_err = models.InitDb()
		if _err != nil {
			panic(_err)
		}
		for {
			fmt.Print("[A] Add New Entry, [L] List Data, [^C] Selesai: ")
			_char, _err := _reader.ReadString('\n')
			if _err != nil {
				panic(_err)
			}
			_char = strings.Trim(_char, " \n")
			switch _char {
			case "a", "A":
				AddEntry()
			case "l", "L":
				ListEntry()
			}
		}
	}()
	<-_interrupt
	fmt.Println()
}
