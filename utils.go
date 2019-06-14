package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"sync"
	"time"

	"context"

	"wallet"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	connectPtr := flag.String("connect", "", "ws, http or ipc path node connection")
	methodPtr := flag.String("method", "", "which operation to execute")
	fromPtr := flag.String("from", "", "address from which the funds are sent")
	toPtr := flag.String("to", "", "address to which funds are sent")
	valuePtr := flag.Int("value", 0, "amount of ether")
	txnumberPtr := flag.Int("txnumber", 0, "number of transactions to execute")
	addressPtr := flag.String("address", "", "address to check balance")
	keystorePathPtr := flag.String("keystore", "", "path to keystore file")
	passwordPrt := flag.String("password", "", "password")

	flag.Parse()

	// convert ether to wei
	multiplier := big.NewInt(1)
	multiplier.SetString("1000000000000000000", 10)
	base := big.NewInt(int64(*valuePtr))
	value := base.Mul(base, multiplier)

	if *connectPtr == "" {
		fmt.Println(`
Please use --connect flag to establish a connection with node
		
Example:
  --connect /home/ubuntu/geth.ipc
		`)
	}
	client, err := ethclient.Dial(*connectPtr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to node")

	if len(os.Args) == 1 {
		log.Fatalf(`

Please choose method:  
  Methods:  
	--TestPerformance [--from x --to y --value 0 --txnumber 0]       - benchmark 
	--GenerateAccount 											     - create account
	--GetBalance [--address] 									     - check balance
	--SendTx [--from x --to y --value 0]						     - send tx 
	--GetPrivateFromKeystore [--keystore path/file --password x]     - get private key from keystore file
  Example: 
	utils --connect /home/ubuntu/store/geth.ipc --method TestPerformance --from 0x0123 --to 0x3210 --value 1 --txnumber 1`)
	}

	switch *methodPtr {

	case "TestPerformance":
		if *fromPtr == "" || *toPtr == "" || value == big.NewInt(0) || *txnumberPtr == 0 {
			log.Fatal("Please specify flags --from, --to, --value, --txnumber")
		}
		// testPerformance(client, "0xB853344f9387304e169B0F0fCB21fEc4AA403375", "0x2Ffd141BbFF6fD973f025E68785c0f9A5759082C", 100000000000000000, 200)
		TestPerformance(client, *fromPtr, *toPtr, value, *txnumberPtr)

	case "GenerateAccount":
		GenerateAccount()

	case "GetBalance":
		if *addressPtr == "" {
			log.Fatal("Please specify flag --address")
		}
		GetBalance(client, *addressPtr)

	case "SendTx":
		if *fromPtr == "" || *toPtr == "" || value == big.NewInt(0) {
			log.Fatal("Please specify flags --from, --to, --value, --txnumber")
		}
		SendTx(client, *fromPtr, *toPtr, value)

	case "GetPrivateFromKeystore":
		if *keystorePathPtr == "" {
			log.Fatal("Please specify flag --keystore")
		}
		key := wallet.GetPrivateFromKeystore(*keystorePathPtr, *passwordPrt)
		fmt.Println("Private key: ", key)
	}
}

func GenerateAccount() {
	wallet.GenerateAccount()
}

func GetBalance(client *ethclient.Client, addr string) {
	address := common.HexToAddress(addr)
	fmt.Printf("%s \n", address.Hex())

	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Current balance: ", balance)

	pendingBalance, err := client.PendingBalanceAt(context.Background(), address)
	fmt.Println("Pending balance: ", pendingBalance)
}

func SendTx(client *ethclient.Client, from string, to string, val *big.Int) {

	payload := []byte(fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_sendTransaction","params":[{"from":"%s", "to": "%s", "value": "0x%x"}],"id":1}`, from, to, val))

	resp, err := http.Post("http://localhost:8504", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Fatal(err)
	}

	type Jsonresp struct {
		Result string
	}

	body := Jsonresp{}

	respbody, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(respbody, &body); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx hash: %s \n", body.Result)

}

func SendTxTest(wg1 *sync.WaitGroup, wg2 *sync.WaitGroup, txch chan string, client *ethclient.Client, from string, to string, val *big.Int) {

	payload := []byte(fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_sendTransaction","params":[{"from":"%s", "to": "%s", "value": "0x%x"}],"id":1}`, from, to, val))

	resp, err := http.Post("http://localhost:8504", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Fatal(err)
	}

	type Jsonresp struct {
		Result string
	}

	body := Jsonresp{}

	respbody, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(respbody, &body); err != nil {
		log.Fatal(err)
	}

	txch <- body.Result
	wg1.Done()
	wg2.Add(1)

}

func listen(wg *sync.WaitGroup, txch chan string, quit chan bool, client *ethclient.Client) {

	// read channel with hashes of sended txs
	var txs []string

LOOP:
	for {
		select {
		case tx := <-txch:
			txs = append(txs, tx)
		case <-quit:
			break LOOP
		}
	}

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
				for _, txitem := range txs {
					if txitem == tx.Hash().Hex() {
						wg.Done()
						fmt.Println("------------------------------------------------")
						fmt.Printf("\n\nHash: %x \n", tx.Hash().Hex())
						fmt.Printf("Value: %s \n", tx.Value().String())
						fmt.Printf("Gas: %d \n", tx.Gas())
						fmt.Printf("Gas price: %d \n", tx.GasPrice().Uint64())
						fmt.Printf("Nonce: %d \n", tx.Nonce())
						fmt.Printf("Data: %x \n", tx.Data())
						fmt.Printf("To: %x \n", tx.To().Hex())
					}
				}
			}
		}
	}
}

func TestPerformance(client *ethclient.Client, from, to string, value *big.Int, num int) {

	start := time.Now()

	var wg1, wg2 sync.WaitGroup
	txch := make(chan string, 1000)
	quit := make(chan bool)

	go listen(&wg2, txch, quit, client)

	for i := 0; i < num; i++ {
		wg1.Add(1)
		go SendTxTest(&wg1, &wg2, txch, client, from, to, value)
	}
	wg1.Wait()
	quit <- true
	wg2.Wait()

	close(txch)
	close(quit)

	fmt.Println("Confirmation time: ", time.Since(start))
}
