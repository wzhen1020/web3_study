// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.0;

import "./NFTAuction.sol";
contract NFTAuctionFactory  {

    // 拍卖合约数组
    NFTAuction[] public auctions;

    event AuctionDeployed(address auction);

    constructor() {
        //  初始化合约
    }

    //  创建新的拍卖合约
    function createAuctionContract() external returns (NFTAuction) {
        NFTAuction auction = new NFTAuction();
        auction.initialize();
        auctions.push(auction);

        emit AuctionDeployed(address(auction));
        return auction;
    }

    function getAuctions() external view returns (NFTAuction[] memory) {
        return auctions;
    }

}