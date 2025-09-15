package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

func main() {
	// RPC 与 WS 端点 (使用主网或测试网)
	rpcClient := rpc.New(rpc.DevNet_RPC) // 测试网 rpc.DevNet_RPC
	wsClient, err := ws.Connect(context.Background(), rpc.MainNetBeta_WS)
	ctx := context.Background()
	if err != nil {
		log.Fatalf("连接 WebSocket 失败: %v", err)
	}
	defer wsClient.Close()

	// ============ 1. 查询最新区块数据 ============
	getRecentBlock(rpcClient)

	// ============ 2. 构造原生 SOL 转账交易 ============
	// 注意：实际执行需提供真实私钥并签名交易，这里仅演示构造

	sendTransferAndConfirm(rpcClient)

	// ============ 3. 实时订阅账户交易 ============

	// 监听的钱包地址
	account := solana.MustPublicKeyFromBase58("你要监听的钱包地址")

	// 用 goroutine 并行监听
	go subscribeAccount(ctx, wsClient, account)

	// 保持主线程运行
	select {}
}

// 查询最新区块信息
func getRecentBlock(c *rpc.Client) {
	ctx := context.Background()
	slot, err := c.GetSlot(ctx, rpc.CommitmentFinalized)
	if err != nil {
		log.Fatalf("获取slot失败: %v", err)
	}
	fmt.Println("当前最新slot:", slot)

	// block, err := c.GetBlock(ctx, slot, &rpc.GetBlockOpts{
	// 	TransactionDetails: rpc.TransactionDetailsSignatures,
	// })
	// if err != nil {
	// 	log.Fatalf("获取区块失败: %v", err)
	// }
	// fmt.Printf("区块 %d 包含 %d 笔交易\n", block.Slot, len(block.Transactions))
}

// 构造并发送一笔SOL转账交易（演示）
func waitForConfirmation(ctx context.Context, c *rpc.Client, sig solana.Signature) {
	fmt.Println("等待交易确认中...")
	for {
		status, err := c.GetSignatureStatuses(ctx, true, sig)
		if err != nil {
			log.Printf("查询状态失败: %v", err)
			time.Sleep(time.Second)
			continue
		}
		if len(status.Value) > 0 && status.Value[0] != nil && status.Value[0].ConfirmationStatus == rpc.ConfirmationStatusFinalized {
			fmt.Println("✅ 交易已确认并落块!")
			break
		}
		fmt.Println("尚未确认...")
		time.Sleep(2 * time.Second)
	}
}

func sendTransferAndConfirm(c *rpc.Client) {
	ctx := context.Background()

	// ⚠️ 实际使用时请换成真实钱包（这里随机生成仅用于演示）
	sender := solana.NewWallet()
	receiver := solana.NewWallet()
	fmt.Println("发送方地址:", sender.PublicKey())
	fmt.Println("接收方地址:", receiver.PublicKey())

	// ⚠️ 记得先给 sender 地址在 Devnet 里空投 SOL，否则会余额不足
	// _, err := c.RequestAirdrop(ctx, sender.PublicKey(), 2_000_000_000)
	// if err != nil { log.Fatalf("空投失败: %v", err) }

	// 构造转账指令
	transferIx := system.NewTransferInstruction(
		1_000_000,            // lamports
		sender.PublicKey(),   // from
		receiver.PublicKey(), // to
	).Build()

	// 获取最近区块哈希
	recent, err := c.GetRecentBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		log.Fatalf("获取区块hash失败: %v", err)
	}

	// 构造交易
	tx, err := solana.NewTransaction(
		[]solana.Instruction{transferIx},
		recent.Value.Blockhash,
		solana.TransactionPayer(sender.PublicKey()),
	)
	if err != nil {
		log.Fatalf("创建交易失败: %v", err)
	}

	// 签名交易
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if sender.PublicKey().Equals(key) {
			return &sender.PrivateKey
		}
		return nil
	})
	if err != nil {
		log.Fatalf("签名失败: %v", err)
	}

	// 发送交易
	sig, err := c.SendTransaction(ctx, tx)
	if err != nil {
		log.Fatalf("发送交易失败: %v", err)
	}
	fmt.Println("交易已广播，签名:", sig)

	// 打印 base64 交易（可用于调试）
	serialized, _ := tx.MarshalBinary()
	fmt.Println("交易Base64:", base64.StdEncoding.EncodeToString(serialized))

	// 等待确认落块
	waitForConfirmation(ctx, c, sig)
}

// 实时订阅账户的交易通知（修正版）
func subscribeAccount(ctx context.Context, wsClient *ws.Client, account solana.PublicKey) {
	fmt.Println("✅ 已订阅账户变动，等待消息中...")

	// 发起订阅
	sub, err := wsClient.AccountSubscribe(account, rpc.CommitmentFinalized)
	if err != nil {
		log.Fatalf("订阅账户失败: %v", err)
	}
	defer sub.Unsubscribe()

	for {
		// 传入 context.Context
		res, err := sub.Recv(ctx)
		if err != nil {
			log.Printf("接收消息出错: %v", err)
			continue
		}
		fmt.Printf("🔔 slot=%d 余额=%d lamports\n",
			res.Context.Slot,
			res.Value.Lamports,
		)
	}
}
