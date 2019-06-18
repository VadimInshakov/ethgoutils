package contract

//go:generate ./generate.sh CUSTOM

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os/exec"
	"regexp"
	"strconv"

	contr "mycontract"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func Compile(contractPath string) {

	r, _ := regexp.Compile("([/\\w-/]+/)([\\w-]+\\.)")
	arr := r.FindStringSubmatch(contractPath)

	contractName := arr[2]
	contractName = string(contractName[:len(contractName)-1])

	// create abi
	abi := exec.Command("solc", "--abi", contractPath, "-o", "build")
	output, err := abi.CombinedOutput()
	if err != nil {
		log.Fatal(fmt.Sprint(err) + ": " + string(output))
		return
	}
	fmt.Println(string(output))
	abi.Wait()

	// create go package
	gopack := exec.Command("abigen", fmt.Sprintf("--abi=./build/%s.abi", contractName), fmt.Sprintf("--pkg=%s", contractName), fmt.Sprintf("--out=./build/%s.go", contractName))
	output, err = gopack.CombinedOutput()
	if err != nil {
		log.Fatal(fmt.Sprint(err) + ": " + string(output))
		return
	}
	fmt.Println(string(output))
	gopack.Wait()

	// create bin
	bin := exec.Command("solc", "--bin", contractPath, "-o", "build")
	output, err = bin.CombinedOutput()
	if err != nil {
		log.Fatal(fmt.Sprint(err) + ": " + string(output))
		return
	}
	fmt.Println(string(output))
	bin.Wait()

	// create go package with deploy method
	godeploy := exec.Command("abigen", fmt.Sprintf("--bin=./build/%s.bin", contractName), fmt.Sprintf("--abi=./build/%s.abi", contractName), fmt.Sprintf("--pkg=%s", contractName), fmt.Sprintf("--out=./build/%s.go", contractName))
	output, err = godeploy.CombinedOutput()
	if err != nil {
		log.Fatal(fmt.Sprint(err) + ": " + string(output))
		return
	}
	fmt.Println(string(output))
	godeploy.Wait()
}

func TokenInfo(client *ethclient.Client, contractaddr string) {

	instance, err := contr.NewTOKENNAME(common.HexToAddress(contractaddr), client)
	if err != nil {
		log.Fatalf("Failed to instantiate contract: %v", err)
	}
	name, err := instance.Name(nil)
	if err != nil {
		log.Fatalf("Failed to retrieve token name: %v", err)
	}
	TotalSupply, err := instance.TotalSupply(nil)
	if err != nil {
		log.Fatal(err)
	}
	Symbol, err := instance.Symbol(nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nToken name: %s \nSymbol: %s \nTotalSupply: %s \n", name, Symbol, TotalSupply)
}

func TransferToken(client *ethclient.Client, priv string, contractaddr string, to string, value int64, gaslimit uint64, gasprice int64) {

	privateKey, err := crypto.HexToECDSA(priv)
	if err != nil {
		log.Fatal(err)
	}

	// auth := bind.NewKeyedTransactor(privateKey)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	// auth.Nonce = big.NewInt(int64(nonce))
	// auth.Value = big.NewInt(0)
	// auth.GasLimit = uint64(gaslimit)
	// auth.GasPrice = big.NewInt(gasprice)

	// instance, err := contr.NewTOKENNAME(common.HexToAddress(contractaddr), client)
	// if err != nil {
	// 	log.Fatalf("Failed to instantiate contract: %v", err)
	// }
	// tx, err := instance.Transfer(auth, common.HexToAddress(to), big.NewInt(value))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("tx sent: %s\n", tx.Hash().Hex())

	toAddress := common.HexToAddress(to)
	tokenAddress := common.HexToAddress(contractaddr)

	hash := crypto.Keccak256Hash([]byte("transfer(address,uint256)"))
	methodID := hash[:4]

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	amount := new(big.Int)
	amount.SetString(strconv.FormatInt(value, 10), 10)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	// gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
	//     To:   &toAddress,
	//     Data: data,
	// })
	// if err != nil {
	//     log.Fatal(err)
	// }

	tx := types.NewTransaction(nonce, tokenAddress, big.NewInt(0), gaslimit, big.NewInt(gasprice), data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s\n", signedTx.Hash().Hex())
	ListenTx(client, signedTx.Hash().Hex())
}

func GetTokenAmount(client *ethclient.Client, priv string, contractaddr string, address string) {

	privateKey, err := crypto.HexToECDSA(priv)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	instance, err := contr.NewTOKENNAME(common.HexToAddress(contractaddr), client)
	if err != nil {
		log.Fatalf("Failed to instantiate contract: %v", err)
	}

	var opts *bind.CallOpts = &bind.CallOpts{From: fromAddress, Pending: true}
	balance, err := instance.BalanceOf(opts, common.HexToAddress(address))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Balance: ", balance)
}

func ListenTx(client *ethclient.Client, txForCheck string) {

	headers := make(chan *types.Header)

	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Fatal(err)
			}
			for _, tx := range block.Transactions() {
				if txForCheck == tx.Hash().Hex() {
					fmt.Printf("\n\nHash: %x \n", tx.Hash().Hex())
					fmt.Printf("Value: %s \n", tx.Value().String())
					fmt.Printf("Gas: %d \n", tx.Gas())
					fmt.Printf("Gas price: %d \n", tx.GasPrice().Uint64())
					fmt.Printf("Nonce: %d \n", tx.Nonce())
					fmt.Printf("Data: %x \n", tx.Data())
					fmt.Printf("To: %x \n", tx.To().Hex())
				}
				return
			}
		}
	}

}
