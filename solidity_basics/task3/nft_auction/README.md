


# NFT 拍卖市场
项目结构:

```lua
nft-auction/
├─ contracts/
│  ├─ MyNFT.sol
│  ├─ MyERC20.sol
│  ├─ NFTAuction.sol
│  ├─ NFTAuctionV2.sol
│  ├─ NFTAuctionFactory.sol
├─ deploy/
│  ├─json/
│  │  ├─proxyNFTAuction.json
│  ├─ 01_deploy_nft_auction.js
│  ├─ 02_upgrades_nft_auction.js
├─ test/
│  ├─ TestNFTAuction.js
│  ├─ TestUpgrade.js
├─ hardhat.config.js
├─ package.json
```
* NFT 合约 + 拍卖合约 + 工厂合约 + 测试 ERC20
* 部署脚本
* 测试脚本
* Hardhat 配置文件
* package.json


运行：
1. npm install  # 安装依赖
2. npx hardhat compile  # 编译合约
3. npx hardhat test  # 运行测试
4. npx hardhat run deploy/01_deploy_nft_auction.js --network localhost  # 部署

## 功能
- NFT 拍卖
- ERC20/ETH 出价，统一转换为 USD
- 使用 Chainlink 喂价
- 工厂模式批量管理拍卖
- 支持合约升级 (UUPS)

