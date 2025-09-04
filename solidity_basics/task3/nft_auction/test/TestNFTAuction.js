const { expect } = require("chai");
const { ethers, upgrades } = require("hardhat");


describe("Test MyNFT", async function () {

    it("MyNft RUN", async function () {
        const [signer, buyer] = await ethers.getSigners()
        // 部署NFT合约
        const MyNft = await ethers.getContractFactory("MyNFT");
        const myNft = await MyNft.deploy();
        await myNft.waitForDeployment();
        const myNftAddress = await myNft.getAddress();
        console.log("MyNft address:", myNftAddress);

        // 部署 NFTAuction 合约
        const NFTAuction = await ethers.getContractFactory("NFTAuction");
        const nftAuction = await NFTAuction.deploy();
        await nftAuction.waitForDeployment();
        const nftAuctionAddress = await nftAuction.getAddress();
        console.log("NFTAuction address:", nftAuctionAddress);
        // 将NFT铸造给测试用户
        myNft.mintNFT(signer.address, "https://my-nft-uri.com/metadata/1");

        // 部署测试ERC20合约
        const MyERC20 = await ethers.getContractFactory("MyERC20");
        const myERC20 = await MyERC20.deploy();
        await myERC20.waitForDeployment();
        const UsdcAddress = await myERC20.getAddress();
        // 将USDC转账给买家
        let tx = await myERC20.connect(signer).transfer(buyer, ethers.parseEther("1000"))
        await tx.wait()
        // console.log("USDC transferred to buyer:", buyer);

        // 部署ETH价格喂价合约
        const MockPriceFeed = await ethers.getContractFactory("MockPriceFeed");
        const mockPriceFeedDeploy = await MockPriceFeed.deploy(ethers.parseEther("10000"))
        const priceFeedEth = await mockPriceFeedDeploy.waitForDeployment()
        const priceFeedEthAddress = await priceFeedEth.getAddress()
        console.log("ethFeed: ", priceFeedEthAddress)
        // 部署USDC价格喂价合约
        const priceFeedUSDCDeploy = await MockPriceFeed.deploy(ethers.parseEther("1"))
        const priceFeedUSDC = await priceFeedUSDCDeploy.waitForDeployment()
        const priceFeedUSDCAddress = await priceFeedUSDC.getAddress()
        console.log("usdcFeed: ", await priceFeedUSDCAddress)
        // 构建价格喂价映射
        const token2Usd = [{
            token: ethers.ZeroAddress,
            priceFeed: priceFeedEthAddress
        }, {
            token: UsdcAddress,
            priceFeed: priceFeedUSDCAddress
        }]
        // 设置价格喂价合约到拍卖合约映射 
        for (let i = 0; i < token2Usd.length; i++) {
            const { token, priceFeed } = token2Usd[i];
            await nftAuction.setTokenPriceFeed(token, priceFeed);
        }

        console.log("设置价格喂价合约完成");
        // 设置NFT合约的授权
        await myNft.connect(signer).setApprovalForAll(nftAuctionAddress, true);

        // 创建拍卖
        await nftAuction.createAuction(1, myNftAddress, 10, ethers.parseEther("0.01"));
        let auction = await nftAuction.auctions(0);
        console.log("拍卖信息:", auction);

        // const ZeroAddress = ethers.ZeroAddress
        //ETH 出价
        tx = await nftAuction.connect(buyer).priceBid(0, ethers.parseEther("0.02"), ethers.ZeroAddress, {
            value: ethers.parseEther("0.2")
        });
        await tx.wait();
        auction = await nftAuction.auctions(0);

        console.log("ETH出价后拍卖信息:", auction);

        // USDC参与竞价
        tx = await myERC20.connect(buyer).approve(nftAuctionAddress, ethers.MaxUint256)
        await tx.wait()
        tx = await nftAuction.connect(buyer).priceBid(0, ethers.parseEther("0.03"), UsdcAddress);
        await tx.wait()
        auction = await nftAuction.auctions(0);
        console.log("USDC出价后拍卖信息:", auction);

        // 4. 结束拍卖
        // 等待 10 s
        await new Promise((resolve) => setTimeout(resolve, 10 * 1000));

        await nftAuction.connect(signer).endAuction(0);

        // 验证结果
        const auctionResult = await nftAuction.auctions(0);
        console.log("结束拍卖后读取拍卖成功：：", auctionResult);
        expect(auctionResult.highestBidder).to.equal(buyer.address);
        expect(auctionResult.highestBid).to.equal(ethers.parseEther("0.03"));

        // 验证 NFT 所有权  
        const owner = await myNft.ownerOf(1);
        console.log("owner::", owner);
        expect(owner).to.equal(buyer.address);
    });
})
//  1000000000000000000
//  2000000000000000000