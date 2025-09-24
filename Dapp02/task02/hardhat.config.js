require("@nomicfoundation/hardhat-toolbox");

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  solidity: "0.8.28",
      localhost: {
      url: "http://127.0.0.1:8545"
    },
    sepolia: {
        url: `https://sepolia.infura.io/v3/INFURA_PROJECT_ID`,
        accounts: `PRIVATE_KEY`
      }
};
