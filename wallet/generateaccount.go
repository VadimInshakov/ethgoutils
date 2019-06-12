package wallet

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
)

const choiceCreateAccount = 1

// ChoiceReader - Read user choice
func ChoiceReader(choice chan uint8) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		cmd, err := strconv.Atoi(scanner.Text())

		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		if err := scanner.Err(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		choice <- uint8(cmd)
		return
	}
}

func Generate(keystoreDir string, pwd string) (common.Address, error) {
	ks := keystore.NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)
	acc, err := ks.NewAccount(pwd)
	if err != nil {
		return common.Address{}, err
	}

	return acc.Address, nil
}

func PassRead(pass chan string) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("enter passphrase")
	for scanner.Scan() {
		passphrase := scanner.Text()

		if err := scanner.Err(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		pass <- passphrase
		return
	}
}

func GenerateAccount() {
	fmt.Println("What action would you like to perform?")
	fmt.Println("1. Create account")
	choice := make(chan uint8)
	go ChoiceReader(choice)
	defer func() { close(choice) }()
	for {
		switch <-choice {
		case 1:
			fmt.Println("Ok let's create account")

			pass := make(chan string)
			go PassRead(pass)
			for {
				switch passphrase, ok := <-pass; ok {
				case true:
					acc, err := Generate("/home/admin/.ethereum/keystore", passphrase)
					if err != nil {
						panic(err)
					}
					fmt.Printf("Account %s created!", acc)
					return
				}
			}

		}
	}

}
