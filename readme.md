**Command line utility for Ethereum written in Golang**

Compile:

    go build

Check tx rate:

    utils --connect <ipc, ws, http path to node> --method TestPerformance --from <from address> --to <to address> --value <amount of ether to be sent> --txnumber <number of transactions>

Generate new account:

    utils --connect <ipc, ws, http path to node> --method GenerateAccount

Get private key from keystore file:

    utils --connect <ipc, ws, http path to node> --method GetPrivateFromKeystore --keystore <path to keystore file> --password <password>

Check balance:

    utils --connect <ipc, ws, http path to node> --method GetBalance --address <address of account>

Send tx:

    utils --connect <ipc, ws, http path to node> --method SendTx --from <from address> --to <to address> --value <amount of ether to be sent>

Generate ABI, go package and compile sol to EVM bytecode:

    utils --method Compile --contract <path to .sol contract file>
