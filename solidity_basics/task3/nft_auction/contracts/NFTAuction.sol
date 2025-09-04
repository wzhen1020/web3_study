// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {AggregatorV3Interface} from "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";
import "hardhat/console.sol";
// upgradeable 版本的库
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
// import "@chainlink/contracts/src/v0.8/interfaces/AggregatorV3Interface.sol";
// import "@openzeppelin/contracts/access/Ownable.sol";

contract NFTAuction is Initializable, UUPSUpgradeable, OwnableUpgradeable {
    // Auction contract code goes here

    /**
     * @dev 拍卖结构体
     */
    struct Auction {
        uint256 tokenId;
        // 卖家
        address seller;
        // NFT地址
        address nftContract;
        // 拍卖持续时间
        uint256 duration;
        // 起始价格
        uint256 startPrice;
        // 拍卖开始时间
        uint256 startTime;
        // 拍卖结束时间
        uint256 endTime;
        // 是否已结束
        bool ended;
        // 最高出价者
        address highestBidder;
        // 最高价
        uint256 highestBid;
        // 支付代币地址 0x0000000000000000000000000000000000000000 表示eth
        address highestTokenAddress;
    }

    // 合约映射
    mapping(uint256 => Auction) public auctions;

    // 管理员
    address public admin;
    // 下一个拍卖的Id
    uint256 public nextAuctionId;

    mapping(address => AggregatorV3Interface) public priceFeeds;

    event AuctionCreated(
        uint256 indexed nextAuctionId,
        address indexed seller,
        address nftContract,
        uint256 tokenId,
        uint256 startPrice,
        uint256 duration
    );

    /// 事件：出价
    event BidPlaced(
        uint256 indexed auctionId,
        address indexed bidder,
        uint256 amount,
        address tokenAddress
    );

    /// 事件：拍卖结束
    event AuctionEnded(
        uint256 indexed auctionId,
        address winner,
        uint256 amount
    );

    function initialize() public initializer {
        __Ownable_init(msg.sender);
        __UUPSUpgradeable_init();
        // owner 已由 OwnableUpgradeable 设置为 msg.sender
        nextAuctionId = 0;
    }

    //   constructor() initializer {}
    /**
     * @dev 设置代币价格预言机
     * @param _tokenAddress 代币地址
     * @param _priceFeed 预言机合约地址
     */
    function setTokenPriceFeed(
        address _tokenAddress,
        address _priceFeed
    ) external {
        priceFeeds[_tokenAddress] = AggregatorV3Interface(_priceFeed);
    }

    /**
     * @dev 预言机获取ERC20价格或ETH价格
     * @param _tokenAddress 代币地址
     * @return 最新价格
     */
    function getChainlinkDataFeedLatestAnswer(
        address _tokenAddress
    ) public view returns (uint256) {
        AggregatorV3Interface priceFeed = priceFeeds[_tokenAddress];

        (
            ,
            /* uint80 roundId */ int256 answer,
            ,
            ,

        ) = /*uint256 startedAt*/ /*uint256 updatedAt*/ /*uint80 answeredInRound*/ priceFeed
                .latestRoundData();

        return uint256(answer);
    }

    /**
     * @dev 创建拍卖
     * @param _tokenId NFT的ID
     * @param _nftContract NFT合约地址
     * @param _duration 拍卖持续时间
     * @param _startPrice 起始价格
     */
    function createAuction(
        uint256 _tokenId,
        address _nftContract,
        uint256 _duration,
        uint256 _startPrice
    ) external {
        // 只有管理员才能创建拍卖
        // require(msg.sender == admin, "Only admin can create auction");
        require(_duration > 0, "Invalid duration");
        require(_startPrice > 0, "Invalid start price");
        // require(ended,"");

        // 转移NFT到合约
        IERC721(_nftContract).transferFrom(msg.sender, address(this), _tokenId);

        // 将拍卖信息存储在映射中
        auctions[nextAuctionId] = Auction({
            tokenId: _tokenId,
            seller: msg.sender,
            nftContract: _nftContract,
            duration: _duration,
            startPrice: _startPrice,
            startTime: block.timestamp,
            endTime: block.timestamp + _duration,
            ended: false,
            highestBidder: address(0),
            highestBid: 0,
            highestTokenAddress: address(0)
        });
        // id自增
        nextAuctionId++;
        emit AuctionCreated(
            nextAuctionId,
            msg.sender,
            _nftContract,
            _tokenId,
            _startPrice,
            _duration
        );
    }

    // 买家出价
    function priceBid(
        uint256 _auctionId,
        uint256 _amount,
        address _tokenAddress
    ) external payable {
        Auction storage auction = auctions[_auctionId];
              console.log("_amount",_amount);
              console.log("auction.startPrice:", auction.startPrice);
        require(block.timestamp < auction.endTime, "Auction ended");
        require(auction.seller != address(0), "Auction not exists");
  

        require(_amount >= auction.startPrice, "Bid too low");
        require(_amount > auction.highestBid, "There already is a higher bid");

        uint payValue;
        if (_tokenAddress == address(0)) {
            // ETH折算

            payValue =
                _amount *
                uint(getChainlinkDataFeedLatestAnswer(_tokenAddress));
        } else {
            // ERC20代币折算
            payValue =
                _amount *
                uint(getChainlinkDataFeedLatestAnswer(address(0)));
        }


        // 折合起始价
        uint startPriceValue = auction.startPrice *
            uint(getChainlinkDataFeedLatestAnswer(auction.highestTokenAddress));

        // 折合最高价
        uint highestBidValue = auction.highestBid *
            uint(getChainlinkDataFeedLatestAnswer(auction.highestTokenAddress));

        // 当前出价大于起始价和最高价
        require(
            payValue >= startPriceValue && payValue > highestBidValue,
            "Bid must be higher than the current highest bid"
        );

        // auction.highestBid = payValue;
        // 转移 ERC20 到合约
        if (_tokenAddress != address(0)) {
            IERC20(_tokenAddress).transferFrom(
                msg.sender,
                address(this),
                _amount
            );
        }

        // 退还前最高价
        if (auction.highestBid > 0) {
            if (auction.highestTokenAddress == address(0)) {
                // auction.highestTokenAddress = _tokenAddress;
                payable(auction.highestBidder).transfer(auction.highestBid);
            } else {
                // 退回之前的ERC20
                IERC20(auction.highestTokenAddress).transfer(
                    auction.highestBidder,
                    auction.highestBid
                );
            }
        }
        // 更新最高出价和最高出价者
        auction.highestTokenAddress = _tokenAddress;
        auction.highestBid = _amount;
        auction.highestBidder = msg.sender;
        emit BidPlaced(_auctionId, msg.sender, _amount, _tokenAddress);
    }

    /**
     * @dev 结束拍卖
     * @param _auctionId 拍卖ID
     */
    function endAuction(uint256 _auctionId) external {
        Auction storage auction = auctions[_auctionId];
        require(block.timestamp >= auction.endTime, "Auction not ended yet");
        require(!auction.ended, "Already ended");

        auction.ended = true;

        if (auction.highestBidder != address(0)) {
            // 转 NFT 给赢家
            IERC721(auction.nftContract).transferFrom(
                address(this),
                auction.highestBidder,
                auction.tokenId
            );

            // 把资金转给卖家
            if (auction.highestTokenAddress == address(0)) {
                payable(auction.seller).transfer(auction.highestBid);
            } else {
                IERC20(auction.highestTokenAddress).transfer(
                    auction.seller,
                    auction.highestBid
                );
            }

            emit AuctionEnded(
                _auctionId,
                auction.highestBidder,
                auction.highestBid
            );
        } else {
            // 无人出价，NFT 退还给卖家
            IERC721(auction.nftContract).transferFrom(
                address(this),
                auction.seller,
                auction.tokenId
            );
        }
    }

      // UUPS 升级权限控制：仅合约 owner 可以升级实现
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}

}
