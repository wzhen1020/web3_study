// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title Uniswap V2 Router 接口
 * @dev 只包含流动性相关方法
 */
interface IUniswapV2Router {
    function addLiquidityETH(
        address token, // 要添加流动性的 ERC20 代币地址
        uint amountTokenDesired, // 用户希望添加的代币数量
        uint amountTokenMin, // 添加代币的最小数量（防止滑点，保证不会少于这个值）
        uint amountETHMin, // 添加 ETH 的最小数量（防止滑点）
        address to, // 流动性份额（LP token）的接收者
        uint deadline // 截止时间（timestamp），超过则交易失败，防止挂单被卡住
    )
        external
        payable
        returns (uint amountToken, uint amountETH, uint liquidity);

    function removeLiquidityETH(
        address token, // 要移除流动性的 ERC20 代币地址
        uint liquidity, // 用户希望移除的 LP token 数量
        uint amountTokenMin, // 最少代币数量（防滑点）
        uint amountETHMin, // 最少 ETH 数量（防滑点）
        address to, // 取回代币和 ETH 的接收者
        uint deadline // 截止时间（timestamp）
    ) external returns (uint amountToken, uint amountETH);
    //返回 Router 内部使用的 WETH 合约地址
    function WETH() external pure returns (address);
}

/**
 * @title Uniswap V2 Factory 接口
 * @dev 查询交易对 Pair 地址
 */
interface IUniswapV2Factory {
    function getPair(
        address tokenA,
        address tokenB
    ) external view returns (address pair);
}

contract MemeToken is ERC20, Ownable, ReentrancyGuard {
    // 手续费比例 单位为万分比  100 = 1%
    // uint256 public tradeTax = 100;

    // 收税者的地址
    address public taxCollector;

    //免税白名单
    mapping(address => bool) public isFeeExempt;

    // 销毁比例（万分比），例如 200 = 2%
    uint256 public burnRate = 200;

    // 已销毁总数
    uint256 public burnCount = 0;

    uint256 public maxTxAmount; // 单笔交易最大额度
    uint256 public maxDailyTxCount; // 每个地址每日最大交易次数

    mapping(address => uint256) public dailyTxCount; // 每日交易计数
    mapping(address => uint256) public lastTxDay; // 上一次交易日（天数）

    // Uniswap Router 地址
    IUniswapV2Router public uniswapRouter;
    // Uniswap Factory
    IUniswapV2Factory public factory;
    // 交易对 LP 代币地址
    address public pair;

    // 构造函数
    constructor(
        string memory name,
        string memory symbol,
        uint256 totalSupply,
        address routerAddress,
        address factoryAddress
    ) ERC20(name, symbol) Ownable(msg.sender) {
        // 部署者初始持有全部代币
        _mint(msg.sender, totalSupply * 10 ** decimals());

        // 部署者、收税者免税
        taxCollector = msg.sender;
        isFeeExempt[msg.sender] = true;

        // 初始化 Uniswap Router 和 factory
        uniswapRouter = IUniswapV2Router(routerAddress);
        factory = IUniswapV2Factory(factoryAddress);
        // 合约自身免税，用于流动性操作
        isFeeExempt[address(this)] = true;

        // 尝试自动查询 TOKEN/WETH Pair
        address existingPair = factory.getPair(
            address(this),
            uniswapRouter.WETH()
        );
        if (existingPair != address(0)) {
            pair = existingPair; // 如果池子存在，自动初始化
        }

        // 设置默认交易限制
        // 默认总数的1%
        maxTxAmount = (totalSupply * 10 ** decimals()) / 100;
        // 每日最大交易次数
        maxDailyTxCount = 5;
    }

    // 转账
    function transfer(
        address to,
        uint256 value
    ) public override returns (bool) {
        address sender = msg.sender;
        require(to != address(0), "Invalid address");
        _applyTxLimits(sender, value);
        uint256 taxAmount = 0;
        uint256 burnAmount = (value * burnRate) / 10000;
        // 剩余金额是否足够
        require(
            balanceOf(sender) >= value + burnAmount,
            "Insufficient balance"
        );
        // 转出和接收双方都不在免税名单中
        if (!isFeeExempt[sender] && !isFeeExempt[to]) {
            //计算交易税
            uint256 tradeTax = getDynamicTaxRate(value);
            // 交易税
            taxAmount = (value * tradeTax) / 10000;
            // 销毁额

            // 余额是否足够交易税及转账
            require(
                balanceOf(sender) >= taxAmount + value + burnAmount,
                "Insufficient balance"
            );
            // 转交易税
            _transfer(sender, taxCollector, taxAmount);
        }
        // 销毁
        _burn(sender, burnAmount);
        burnCount += burnAmount;
        // 转帐
        _transfer(sender, to, value);


        return true;
    }

    // 授权转账
    function transferFrom(
        address from,
        address to,
        uint256 value
    ) public override returns (bool) {
        require(from != address(0), "ERC20: zero sender");
        require(to != address(0), "ERC20: zero recipient");
        require(value > 0, "ERC20: zero amount");
        _applyTxLimits(from, value);
        uint256 taxAmount = 0;
        uint256 burnAmount = (value * burnRate) / 10000;

        address spender = _msgSender();
        // 转出和接收双方都不在免税名单中
        if (!isFeeExempt[from] && !isFeeExempt[to]) {
            //计算交易税
            uint256 tradeTax = getDynamicTaxRate(value);
            taxAmount = (value * tradeTax) / 10000;
            // 余额是否足够交易税、转账、销毁额
            require(
                balanceOf(from) >= taxAmount + value + burnAmount,
                "Insufficient balance"
            );
            _spendAllowance(spender, from, value);
            _transfer(from, taxCollector, taxAmount);
        } else {
            _spendAllowance(spender, from, value);
        }
        // 销毁
        _burn(from, burnAmount);
        burnCount += burnAmount;
        // 转账
        _transfer(from, to, value);

        return true;
    }

    // ------------------------------
    //交易税功能
    //-------------------------------

    // // 设置交易税
    // function setTaxFee(uint256 tax) external onlyOwner {
    //     // 交易税最高20%
    //     require(tax <= 2000, "tax too high");
    //     tradeTax = tax;
    // }

    // 收税人
    function setTaxCollector(address newCollector) external onlyOwner {
        require(newCollector != address(0), "zero address");
        taxCollector = newCollector;
    }

    // 设置免税
    function setFeeEExempt(address account, bool exempt) external onlyOwner {
        isFeeExempt[account] = exempt;
    }

    /**
     * @dev 自定义动态税率逻辑（单位: 万分比）
     * 例如:
     * - 交易额 >= 1000 代币 => 2% 税
     * - 交易额 >= 100 代币  => 1% 税
     * - 其他                => 0.5% 税
     */
    function getDynamicTaxRate(uint256 amount) public view returns (uint256) {
        if (amount >= 1000 * 10 ** decimals()) {
            return 200; // 2%
        } else if (amount >= 100 * 10 ** decimals()) {
            return 100; // 1%
        } else {
            return 50; // 0.5%
        }
    }
    // ------------------------------
    //销毁功能
    //-------------------------------

    // 设置销毁比例（最高 5%）
    function setBurnRate(uint256 rate) external onlyOwner {
        require(rate <= 500, "burn too high");
        burnRate = rate;
    }
    // ------------------------------
    //最大交易限制功能
    //-------------------------------

    // 设置单笔交易最大额度
    function setMaxTxAmount(uint256 newTaxTxAmount) external onlyOwner {
        maxTxAmount = newTaxTxAmount;
    }
    // 设置 每个地址每日最大交易次数
    function setMaxDailyTxCount(uint256 newMaxDailyTxCount) external onlyOwner {
        maxDailyTxCount = newMaxDailyTxCount;
    }

    function _applyTxLimits(address sender, uint256 amount) internal {
        // 免税名单不限制
        if (isFeeExempt[sender]) return;

        // 检查单笔交易额度
        require(
            amount < maxTxAmount,
            "Exceeding the limit for a single transaction"
        );

        uint256 today = block.timestamp / 1 days;
        // 如果是新的一天，重置计数
        if (lastTxDay[sender] < today) {
            dailyTxCount[sender] = 0;
            lastTxDay[sender] = today;
        }

        // 检查每日交易次数
        require(
            dailyTxCount[sender] < maxDailyTxCount,
            "Exceeding the daily transaction limit"
        );
        // 增加今日交易次数
        dailyTxCount[sender] += 1;
    }

    // ------------------------------
    //流动性池功能
    //-------------------------------
    // 设置 LP Token 的 pair 地址
    function setPair(address _pair) external onlyOwner {
        require(_pair != address(0), "Invalid pair");
        pair = _pair;
    }

    /**
     * @notice 向 Uniswap 添加 ETH + Token 流动性
     * @param tokenAmount 用户希望提供的代币数量
     * @param amountETHMin 添加 ETH 的最小数量（防止滑点）
     * @param deadline 截止时间 秒级时间戳
     * @dev 用户需先 approve 合约足够 token
     */
    function addLiquidity(
        uint256 tokenAmount,
        uint256 amountETHMin,
        uint256 deadline
    ) external payable nonReentrant {
        address sender = msg.sender;
        require(
            tokenAmount > 0 && msg.value > 0 && amountETHMin > 0,
            "Invalid quantity"
        );
        require(
            deadline >= block.timestamp + 300,
            "Ensure that the transaction remains valid for at least 5 minutes in the future."
        );
        // 临时将用户添加到免税白名单，避免交易税
        bool prevExempt = isFeeExempt[sender];
        isFeeExempt[sender] = true;

        // 将用户代币转入合约
        _transfer(sender, address(this), tokenAmount);
        // 授权Router 使用合约内代币
        _approve(address(this), address(uniswapRouter), tokenAmount);

        // 调用 Uniswap Router 添加流动性
        uniswapRouter.addLiquidityETH{value: msg.value}(
            address(this), // 代币地址
            tokenAmount, // 代币数量
            0, // 最少代币，防滑点用
            amountETHMin, // 最少 ETH
            sender, // LP token 接收者
            deadline // 截止时间 block.timestamp + 300
        );
        // 恢复用户原有免税状态
        isFeeExempt[sender] = prevExempt;
    }

    /**
     * @notice 移除 Uniswap 流动性（ETH + Token）
     * @param liquidity 用户希望移除的 LP token 数量
     * @param amountETHMin 添加 ETH 的最小数量（防止滑点）
     * @param deadline 截止时间 秒级时间戳
     * @dev 用户需先 approve 合约足够 LP token
     */
    function removeLiquidity(
        uint256 liquidity,
        uint256 amountETHMin,
        uint256 deadline
    ) external nonReentrant {
        require(pair != address(0), "LP pair not set");
        require(liquidity > 0, "Invalid quantity");
        require(
            deadline >= block.timestamp + 300,
            "Ensure that the transaction remains valid for at least 5 minutes in the future."
        );
        address sender = msg.sender;
        // 临时免税
        bool prevExempt = isFeeExempt[sender];
        isFeeExempt[sender] = true;

        // 将 LP token 转入合约
        IERC20(pair).transferFrom(sender, address(this), liquidity);
        IERC20(pair).approve(address(uniswapRouter), liquidity);

        // 调用 Router 移除流动性
        uniswapRouter.removeLiquidityETH(
            address(this), // 代币地址
            liquidity, // LP token 数量
            0, // 最少代币
            amountETHMin, // 最少 ETH
            sender, // 接收者
            deadline
        );

        // 恢复免税状态
        isFeeExempt[sender] = prevExempt;
    }
}
