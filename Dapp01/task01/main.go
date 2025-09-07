package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"task01/counter"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 测试网 api_key
const API_KEY = "https://sepolia.infura.io/v3/<KEY>"
const FROM__PRIVATE_KEY = ""
const TO_ADDRESS = "0x0000000000000000000000000000000000000000"

func main() {
	// queryBlock(5671744)

	//
	//
	transferETH(FROM__PRIVATE_KEY, TO_ADDRESS, 1000000000000000)
}

/*
使用 Sepolia 测试网络实现基础的区块链交互，包括查询区块和发送交易。
 具体任务
环境搭建
安装必要的开发工具，如 Go 语言环境、 go-ethereum 库。
注册 Infura 账户，获取 Sepolia 测试网络的 API Key。
查询区块
编写 Go 代码，使用 ethclient 连接到 Sepolia 测试网络。
实现查询指定区块号的区块信息，包括区块的哈希、时间戳、交易数量等。
输出查询结果到控制台。

*/

func queryBlock(blockNumber int64) {

	// 创建一个 Ethereum 客户端
	client, err := ethclient.Dial(API_KEY)

	if err != nil {
		log.Fatal(err)
	}
	// 获得完整区块
	block, err := client.BlockByNumber(context.Background(), big.NewInt(blockNumber))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(block.Number().Uint64())            // 5671744
	fmt.Println("时间戳:", block.Time())               // 1712798400
	fmt.Println(block.Difficulty().Uint64())        // 0
	fmt.Println("区块的哈希:", block.Hash().Hex())       // 0xae713dea1419ac72b928ebe6ba9915cd4fc1ef125a606f90f5e783c47cb1a4b5
	fmt.Println("交易数量:", len(block.Transactions())) // 70

}

/*发送交易
准备一个 Sepolia 测试网络的以太坊账户，并获取其私钥。
编写 Go 代码，使用 ethclient 连接到 Sepolia 测试网络。
构造一笔简单的以太币转账交易，指定发送方、接收方和转账金额。
对交易进行签名，并将签名后的交易发送到网络。
输出交易的哈希值。*/

func transferETH(fromPrivateKey string, toAddress string, amount int64) {
	client, err := ethclient.Dial(API_KEY)

	if err != nil {
		log.Fatal(err)
	}

	// 通过私钥创建一个新的交易
	privateKey, err := crypto.HexToECDSA(fromPrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}

	// 从私钥获取公钥
	publicKey := privateKey.Public()
	// 类型断言
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)

	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	// 从公钥获取地址
	fromAdderss := crypto.PubkeyToAddress(*publicKeyECDSA)
	// 获取nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAdderss)

	if err != nil {
		log.Fatal(err)
	}

	// 交易数量
	// 转账的Gas限制
	gasLimit := uint64(21000)
	//  SuggestGasPrice  获取建议的 gas 价格 根据'x'个先前块来获得平均燃气价格。
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// 收款地址

	var data []byte
	// 创建一笔新的交易 NewTransaction 创建一个新的交易
	tx := types.NewTransaction(nonce, common.HexToAddress(toAddress), big.NewInt(amount), gasLimit, gasPrice, data)
	// 获取链 ID
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// 使用私钥对交易进行签名
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	// 发送交易 SendTransaction 来将已签名的事务广播到整个网络。
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	// 打印交易哈希 signedTx.Hash().Hex()
	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
}

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
func DeployContractByGO() {
	client, err := ethclient.Dial(API_KEY)
	if err != nil {
		log.Fatal(err)
	}

	// privateKey, err := crypto.GenerateKey()
	// privateKeyBytes := crypto.FromECDSA(privateKey)
	// privateKeyHex := hex.EncodeToString(privateKeyBytes)
	// fmt.Println("Private Key:", privateKeyHex)
	privateKey, err := crypto.HexToECDSA(FROM__PRIVATE_KEY)
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
	client, err := ethclient.Dial(API_KEY)
	if err != nil {
		log.Fatal(err)
	}
	counterContract, err := counter.NewCounter(common.HexToAddress(contractAddr), client)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA(FROM__PRIVATE_KEY)
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
