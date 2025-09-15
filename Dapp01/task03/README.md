```mermaid

flowchart TD
  A[用户创建交易（指令+账户）]
  B[序列化 -> 最近 blockhash 设置]
  C[签名（payer + required signers）]
  D[客户端 RPC 提交 sendTransaction / sendRawTransaction]
  E[RPC 节点 mempool (pending)]
  F[Leader 节点 排序与验证]
  G[交易被写入区块（block）]
  H[交易执行：逐个 Instruction 被 BPF 程序处理]
  I[状态更新：账户数据 & lamports 更改写入 Accounts DB]
  J[交易确认确认度：processed -> confirmed -> finalized]
  K[事件/日志 (program logs) 可通过 RPC 或 WS 订阅读取]

  A --> B --> C --> D --> E --> F --> G --> H --> I --> J --> K
```
