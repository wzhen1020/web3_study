require("@nomicfoundation/hardhat-toolbox");
require("dotenv").config();
require("hardhat-deploy");
require('@openzeppelin/hardhat-upgrades');
/** @type import('hardhat/config').HardhatUserConfig */
console.log("INFURA_API_KEY:", process.env.INFURA_API_KEY);
console.log("PRIVATE_KEY:", process.env.PRIVATE_KEY);
module.exports = {
  solidity: "0.8.28",
  networks: {
    sepolia: {
      // Infura URL sepolia 获取地址：https://developer.metamask.io/key/active-endpoints
      url: `https://sepolia.infura.io/v3/${process.env.INFURA_API_KEY}`,
      // 钱包KEY
      accounts: [process.env.PRIVATE_KEY]
    }
  }
};