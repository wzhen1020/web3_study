package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"task02/counter"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

/*
使用 abigen 工具自动生成 Go 绑定代码，用于与 Sepolia 测试网络上的智能合约进行交互。

	 具体任务
		编写智能合约
		使用 Solidity 编写一个简单的智能合约，例如一个计数器合约。
		编译智能合约，生成 ABI 和字节码文件。
		使用 abigen 生成 Go 绑定代码
		安装 abigen 工具。
		使用 abigen 工具根据 ABI 和字节码文件生成 Go 绑定代码。
		使用生成的 Go 绑定代码与合约交互
		编写 Go 代码，使用生成的 Go 绑定代码连接到 Sepolia 测试网络上的智能合约。
		调用合约的方法，例如增加计数器的值。
		输出调用结果。
*/
func main() {
	runCounter()
}

func DeployContractByGO() {
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/2275926dc49f42ec894bcbe253fb047a")
	if err != nil {
		log.Fatal(err)
	}

	// privateKey, err := crypto.GenerateKey()
	// privateKeyBytes := crypto.FromECDSA(privateKey)
	// privateKeyHex := hex.EncodeToString(privateKeyBytes)
	// fmt.Println("Private Key:", privateKeyHex)
	privateKey, err := crypto.HexToECDSA("a0c5fb31bf26b2ab6673b218e80b2e8e7c3410337a9dc86405416167559fd25c")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	chainId, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	address, tx, instance, err := counter.DeployCounter(auth, client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(address.Hex())
	fmt.Println(tx.Hash().Hex())

	_ = instance

	// 0x95feb2011faaF4eFa7Deff827a33a3780B3eEA20

	//0x18714dd9d38908f206524dcfa4600e25ed359935e3f0120876cba4bcfe9dd6c3
}

const (
	contractAddr = "0x95feb2011faaF4eFa7Deff827a33a3780B3eEA20"
)

func runCounter() {
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/2275926dc49f42ec894bcbe253fb047a")
	if err != nil {
		log.Fatal(err)
	}
	counterContract, err := counter.NewCounter(common.HexToAddress(contractAddr), client)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA("a0c5fb31bf26b2ab6673b218e80b2e8e7c3410337a9dc86405416167559fd25c")
	if err != nil {
		log.Fatal(err)
	}

	var value = big.NewInt(1)

	// 链ID
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	opt, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}
	tx, err := counterContract.Increment(opt, value)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("tx hash:", tx.Hash().Hex())

	fmt.Println("is value saving in contract equals to origin value:", tx.Data())
}
