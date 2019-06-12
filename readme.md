**Command line utility for Ethereum written in Golang**

Compile:

    go build

Check tx rate:

    utils --connect <ipc, ws, http path to node> --method TestPerformance --from <from address> --to <to address> --value <amount of ether to be sent> --txnumber <number of transactions>

Generate new account:

    utils --connect <ipc, ws, http path to node> --method GenerateAccount

Check balance:

    utils --connect <ipc, ws, http path to node> --method GetBalance --address <address of account>

Send tx:

    utils --connect <ipc, ws, http path to node> --method SendTx --from <from address> --to <to address> --value <amount of ether to be sent>
