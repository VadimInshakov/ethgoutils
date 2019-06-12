**Command line utility for Ethereum written in Golang**

Check tx rate:

    go run perftest.go --connect <ipc, ws, http path to node> --method TestPerformance --from <from address> --to <to address> --value <amount of ether to be sent> --txnumber <number of transactions>

Generate new account:

    go run perftest.go --connect <ipc, ws, http path to node> --method GenerateAccount

Check balance:

    go run perftest.go --connect <ipc, ws, http path to node> --method GetBalance --address <address of account>

Send tx:

    go run perftest.go --connect <ipc, ws, http path to node> --method SendTx --from <from address> --to <to address> --value <amount of ether to be sent>
