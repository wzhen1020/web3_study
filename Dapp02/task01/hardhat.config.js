require("@nomicfoundation/hardhat-toolbox");

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  solidity: "0.8.28",
  networks: {
    sepolia: {
      // Infura URL sepolia 获取地址：https://developer.metamask.io/key/active-endpoints
      url: `https://sepolia.infura.io/v3/${process.env.INFURA_API_KEY}`,
      // 钱包KEY
      accounts: [process.env.WALLET_PRIVATE_KEY]
    }
  }
};
