const { ethers, deployments, upgrades } = require("hardhat");
const { expect } = require("chai");

describe("Test upgrade", function () {

  it("should create an deploy", async function () {
    const [signer, buyer] = await ethers.getSigners()
    // 部署NFT合约
    const MyNft = await ethers.getContractFactory("MyNFT");
    const myNft = await MyNft.deploy();
    await myNft.waitForDeployment();
    const myNftAddress = await myNft.getAddress();
    console.log("MyNft address:", myNftAddress);
       // 将NFT铸造给测试用户
        myNft.mintNFT(signer.address, "https://my-nft-uri.com/metadata/1");
    // 1.部署业务合约
    await deployments.fixture("deployNFTAuction");
    const nftAuctionProxy = await deployments.get("NFTAuctionProxy");



    // 2.调用 createAuction 方法 创建拍卖

    const nftAuction = await ethers.getContractAt("NFTAuction", nftAuctionProxy.address);
    const nftAuctionAddress = await nftAuction.getAddress();
    console.log("NFTAuction address:", nftAuctionAddress);
    // 设置NFT合约的授权
    await myNft.connect(signer).setApprovalForAll(nftAuctionAddress, true);

    await nftAuction.createAuction(1, myNftAddress, 10, ethers.parseEther("0.01"));
    let auction = await nftAuction.auctions(0);
    console.log("拍卖信息:", auction);

    console.log("升级前", await upgrades.erc1967.getImplementationAddress(nftAuctionProxy.address));

    // 3.升级合约
    await deployments.fixture("upgradeNFTAuction");
    //4.读取合约
    const auction2 = await nftAuction.auctions(0);
    console.log("升级后", await upgrades.erc1967.getImplementationAddress(nftAuctionProxy.address));
    console.log("升级后的合约::", auction2);

    const nftAuctionV2 = await ethers.getContractAt("NFTAuctionV2", nftAuctionProxy.address);
    const result = await nftAuctionV2.testUpgrade();
    console.log("hello::", result);
    expect(auction2.startTime).to.equal(auction.startTime);


  });
});

