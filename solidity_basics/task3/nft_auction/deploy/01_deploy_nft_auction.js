const { upgrades, ethers } = require("hardhat");
const fs = require("fs");
const path = require("path");
module.exports = async function ({ getNamedAccounts, deployments }) {
    const { save } = deployments;
    const { deployer } = await getNamedAccounts();

    console.log("部署者用户地址：", deployer)
    const NFTAuction = await ethers.getContractFactory("NFTAuction");

    // 通过代理部署合约
    const ntfAuctionProxy = await upgrades.deployProxy(NFTAuction, [], { initializer: 'initialize' });

    await ntfAuctionProxy.waitForDeployment();
    const proxyAddress = await ntfAuctionProxy.getAddress();
    console.log("代理合约地址:", proxyAddress);
    const implAddress = await upgrades.erc1967.getImplementationAddress(proxyAddress);
    console.log("目标合约地址:", implAddress);


    const storePath = path.resolve(__dirname, "./json/proxyNFTAuction.json");



    fs.writeFileSync(
        storePath,
        JSON.stringify({
            proxyAddress,
            implAddress,
            abi: NFTAuction.interface.format("json"),
        })
    );

    await save("NFTAuctionProxy", {
        abi: NFTAuction.interface.format("json"),
        address: proxyAddress,
        // args: [],
        // log: true,
    })


    // await deploy("NFTAuction", {
    //     from: deployer,

    //     log: true,
    // });
};

module.exports.tags = ["deployNFTAuction"];