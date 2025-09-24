// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts-upgradeable/token/ERC20/IERC20Upgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/security/ReentrancyGuardUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC20/utils/SafeERC20Upgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";

/**
 * @title
 * @author
 * @notice
 */
contract TaskStake is
    Initializable,
    UUPSUpgradeable,
    ReentrancyGuardUpgradeable,
    AccessControlUpgradeable
{
    using SafeERC20Upgradeable for IERC20Upgradeable;

    // 放大比例，用于 accRewardPerShare 的精度计算
    uint256 private constant ACC_PRECISION = 1e18;

    // ---------- 角色定义 ----------
    bytes32 public constant ADMIN_ROLE    = keccak256("ADMIN_ROLE");
    bytes32 public constant PAUSER_ROLE   = keccak256("PAUSER_ROLE");
    bytes32 public constant UPGRADER_ROLE = keccak256("UPGRADER_ROLE");
    // DEFAULT_ADMIN_ROLE 由 AccessControl 自带，用于管理其它角色 & 管理合约敏感设置

    // 质押奖励 token（用于发放奖励） 自定义ERC20
    IERC20Upgradeable public taskToken;

    // 质押 开始区块（区块号）
    uint256 public startBlock;

    // 质押结束区块（区块号）
    uint256 public endBlock;

    // 质押池总权重 （用于按权重分配每区块奖励）
    uint256 public totalPoolWeight;

    // 每个区块 奖励多少 （各质押池按照权重分）
    uint256 public eachBlockRewardTokens;

    // 第一个池的 pid = 0 表示以太坊原生币（ETH）质押
    uint256 public ETH_PID = 0;

    // 质押池结构体
    struct Pool {
        // 质押代币地址（address(0) 表示 ETH）
        address stTokenAddress;
        // 质押池权重
        uint256 poolWeight;
        // 最后一次计算奖励的区块号
        uint256 lastRewardBlock;
        // 每份质押代币累计的奖励（放大 ACC_PRECISION）
        uint256 accTaskTokenPerShare;
        // 质押池中总质押代币数量 （ETH 为 wei）
        uint256 stTokenAmount;
        // 最小质押代币数量
        uint256 minSTTokenAmount;
        // 解除质押的锁定区块数
        uint256 unStakeLockedBlocks;
    }

    // 用户解除质押请求结构体
    struct UnStakeRequest {
        // 解质押代币数量
        uint256 amount;
        // 到达此区块后可提取质押代币
        uint256 unlockBlock;
    }

    struct User {
        // 用户在池中的质押数量
        uint256 amount;
        // 已记账的奖励债务（用于计算 pending）
        uint256 rewardDebt;
        // 等待领取的 token数量 （还未领取的 taskToken）
        uint256 pending;
        // 解质押请求列表，每个请求包含解质押数量和解锁区块。
        UnStakeRequest[] requests;
    }

    // 质押池列表
    Pool[] public pools;
    // 用户质押信息映射 poolId => user => info
    mapping(uint256 => mapping(address => User)) public users;


    /* ========== 暂停/恢复 操作接口（由 PAUSER_ROLE 管理） ========== */

    // 存储暂停状态（新增变量，默认在 initialize 中设置为 false）
    bool public depositPaused;
    bool public unstakePaused;
    bool public claimPaused;
    bool public withdrawUnstakedPaused;

    // 事件
    event AddPool(
        uint256 indexed pid,
        address stTokenAddress,
        uint256 poolWeight
    );
    event Deposit(address indexed user, uint256 indexed pid, uint256 amount);
    event DepositETH(address indexed user, uint256 indexed pid, uint256 amount);
    event UnstakeRequested(
        address indexed user,
        uint256 indexed pid,
        uint256 amount,
        uint256 unlockBlock
    );
    event Claim(address indexed user, uint256 indexed pid, uint256 amount);
    event EmergencyWithdraw(
        address indexed user,
        uint256 indexed pid,
        uint256 amount
    );
    event UpdateOwner(address indexed oldOwner, address indexed newOwner);
    event UpdateEachBlockReward(uint256 oldReward, uint256 newReward);
    event UpdateEndBlock(uint256 oldEnd, uint256 newEnd);

    // 暂停/恢复事件（每类独立）
    event PauseDeposit(bool paused);
    event PauseUnstake(bool paused);
    event PauseClaim(bool paused);
    event PauseWithdrawUnstaked(bool paused);

    /* ========== 修饰器 ========== */

    modifier onlyOwner() {
        require(
            hasRole(DEFAULT_ADMIN_ROLE, msg.sender),
            "TaskStake: caller is not admin"
        );
        _;
    }

    // 操作级暂停检查修饰器
    modifier whenDepositNotPaused() {
        require(!depositPaused, "TaskStake: deposit is paused");
        _;
    }

    modifier whenUnstakeNotPaused() {
        require(!unstakePaused, "TaskStake: unstake is paused");
        _;
    }

    modifier whenClaimNotPaused() {
        require(!claimPaused, "TaskStake: claim is paused");
        _;
    }

    modifier whenWithdrawUnstakedNotPaused() {
        require(
            !withdrawUnstakedPaused,
            "TaskStake: withdrawUnstaked is paused"
        );
        _;
    }

    /* ========== 初始化 ========== */

    /**
     *
     * @param _taskToken 奖励token
     * @param _startBlock 质押开始区块
     * @param _endBlock 质押结束区块
     * @param _eachBlockRewardTokens 每个区块奖励多少 token
     */
    function initialize(
        IERC20Upgradeable _taskToken,
        uint256 _startBlock,
        uint256 _endBlock,
        uint256 _eachBlockRewardTokens
    ) public initializer {
        __UUPSUpgradeable_init();
        __ReentrancyGuard_init();
        __AccessControl_init();

        require(address(_taskToken) != address(0), "taskToken is zero address");
        require(_startBlock < _endBlock, "startBlock must be < endBlock");

        taskToken = _taskToken;
        startBlock = _startBlock;
        endBlock = _endBlock;
        eachBlockRewardTokens = _eachBlockRewardTokens;

        // 授予初始角色：管理员、升级角色、暂停角色
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
              _grantRole(ADMIN_ROLE, msg.sender);
        _grantRole(UPGRADER_ROLE, msg.sender);
        _grantRole(PAUSER_ROLE, msg.sender);

        // 初始化：默认不暂停任何操作
        depositPaused = false;
        unstakePaused = false;
        claimPaused = false;
        withdrawUnstakedPaused = false;
    }

    /* ========== 管理函数 ========== */

    /**
     *  添加质押池
     * @param _stTokenAddress 质押token地址
     * @param _poolWeight 质押池权重
     * @param _minSTTokenAmount 最小质押数
     * @param _unStakeLockedBlocks 解除质押的锁定区块数
     *
     * 仅管理员（DEFAULT_ADMIN_ROLE）可调用
     */
    function addPool(
        address _stTokenAddress,
        uint256 _poolWeight,
        uint256 _minSTTokenAmount,
        uint256 _unStakeLockedBlocks
    ) public onlyRole(ADMIN_ROLE) {
        require(_unStakeLockedBlocks > 0, "invalid withdraw locked blocks");
        require(block.number < endBlock, "already ended");
        // 第一个质押池只能是 native token
        if (pools.length == 0) {
            // 第一个池要求为 ETH
            require(
                _stTokenAddress == address(0),
                "first pool must be ETH (address(0))"
            );
        } else {
            require(
                _stTokenAddress != address(0),
                "staking token must be ERC20"
            );
        }

        uint256 lastRewardBlock = block.number > startBlock
            ? block.number
            : startBlock;
        pools.push(
            Pool({
                stTokenAddress: _stTokenAddress,
                poolWeight: _poolWeight,
                lastRewardBlock: lastRewardBlock,
                accTaskTokenPerShare: 0,
                stTokenAmount: 0,
                minSTTokenAmount: _minSTTokenAmount,
                unStakeLockedBlocks: _unStakeLockedBlocks
            })
        );
        totalPoolWeight = totalPoolWeight + _poolWeight;
        emit AddPool(pools.length - 1, _stTokenAddress, _poolWeight);
    }

    /**
     * 更新每区块奖励
     *
     */
    function setEachBlockReward(
        uint256 _eachBlockRewardTokens
    ) external onlyRole(ADMIN_ROLE) {
        emit UpdateEachBlockReward(
            eachBlockRewardTokens,
            _eachBlockRewardTokens
        );
        eachBlockRewardTokens = _eachBlockRewardTokens;
    }

    /**
     * 更新结束区块
     */
    function setEndBlock(uint256 _endBlock) external onlyOwner {
        require(_endBlock > startBlock, "endBlock must be > startBlock");
        emit UpdateEndBlock(endBlock, _endBlock);
        endBlock = _endBlock;
    }



    /**
     * 控制 deposit（ERC20/ETH 存入）是否暂停
     * 仅 PAUSER_ROLE 可调用
     */
    function setDepositPaused(bool _paused) external onlyRole(PAUSER_ROLE) {
        depositPaused = _paused;
        emit PauseDeposit(_paused);
    }

    /**
     * 控制 unstake（生成解除请求）是否暂停
     * 仅 PAUSER_ROLE 可调用
     */
    function setUnstakePaused(bool _paused) external onlyRole(PAUSER_ROLE) {
        unstakePaused = _paused;
        emit PauseUnstake(_paused);
    }

    /**
     * 控制 claim（领取奖励）是否暂停
     * 仅 PAUSER_ROLE 可调用
     */
    function setClaimPaused(bool _paused) external onlyRole(PAUSER_ROLE) {
        claimPaused = _paused;
        emit PauseClaim(_paused);
    }

    /**
     * 控制 withdrawUnstaked（提取已解锁 unstake 请求）是否暂停
     * 仅 PAUSER_ROLE 可调用
     */
    function setWithdrawUnstakedPaused(bool _paused) external onlyRole(PAUSER_ROLE) {
        withdrawUnstakedPaused = _paused;
        emit PauseWithdrawUnstaked(_paused);
    }

    /* ========== 用户操作 ========== */

    /**
     * 质押存入 ERC20 质押代币
     * - 先更新池奖励，然后结算用户 pending
     * - 要求 _pid != ETH_PID
     * @param _pid 质押池ID
     * @param _amount 质押代币数量
     */
    function deposit(
        uint256 _pid,
        uint256 _amount
    ) public whenDepositNotPaused {
        require(_pid < pools.length, "invalid pid");
        require(
            _pid != ETH_PID,
            "deposit not support ETH via deposit(); use depositETH"
        );
        Pool storage pool = pools[_pid];
        require(_amount >= pool.minSTTokenAmount, "deposit amount too small");
        require(block.number >= startBlock, "not started");
        require(block.number <= endBlock, "already ended");

        // 先更新当前质押池的奖励信息
        _updatePoolReward(_pid);
        //  查询该质押池中 当前用户的质押信息
        User storage user = users[_pid][msg.sender];
        // 结算用户的 pending（基于新的 accTaskTokenPerShare）
        if (user.amount > 0) {
            uint256 accumulated = (user.amount * pool.accTaskTokenPerShare) /
                ACC_PRECISION;
            if (accumulated > user.rewardDebt) {
                uint256 pendingNow = accumulated - user.rewardDebt;
                user.pending = user.pending + pendingNow;
            }
        }

        // 将质押代币转入合约
        if (_amount > 0) {
            IERC20Upgradeable(pool.stTokenAddress).safeTransferFrom(
                msg.sender,
                address(this),
                _amount
            );
            user.amount = user.amount + _amount;
            pool.stTokenAmount = pool.stTokenAmount + _amount;
        }

        // 更新 rewardDebt 为当前已记账值
        user.rewardDebt =
            (user.amount * pool.accTaskTokenPerShare) /
            ACC_PRECISION;

        emit Deposit(msg.sender, _pid, _amount);
    }

    /**
     * 质押存入 ETH（第一个池）
     * @param _pid 质押池ID
     */
    function depositETH(
        uint256 _pid
    ) external payable nonReentrant whenDepositNotPaused {
        require(_pid == ETH_PID, "depositETH only for ETH pool");
        require(_pid < pools.length, "invalid pid");
        Pool storage pool = pools[_pid];
        require(pool.stTokenAddress == address(0), "pool not ETH");
        require(msg.value >= pool.minSTTokenAmount, "deposit amount too small");
        require(block.number >= startBlock, "not started");
        require(block.number <= endBlock, "already ended");

        uint256 _amount = msg.value;

        // 更新池奖励
        _updatePoolReward(_pid);

        User storage user = users[_pid][msg.sender];

        if (user.amount > 0) {
            uint256 accumulated = (user.amount * pool.accTaskTokenPerShare) /
                ACC_PRECISION;
            if (accumulated > user.rewardDebt) {
                uint256 pendingNow = accumulated - user.rewardDebt;
                user.pending = user.pending + pendingNow;
            }
        }

        if (_amount > 0) {
            user.amount = user.amount + _amount;
            pool.stTokenAmount = pool.stTokenAmount + _amount;
        }

        user.rewardDebt =
            (user.amount * pool.accTaskTokenPerShare) /
            ACC_PRECISION;

        emit DepositETH(msg.sender, _pid, _amount);
    }

    /**
     * 领取奖励
     * - 将 user.pending 发放并清零
     * - 先更新池奖励再计算
     * @param _pid 质押池Id
     */
    function claim(uint256 _pid) public whenClaimNotPaused {
        require(_pid < pools.length, "invalid pid");
        Pool storage pool = pools[_pid];
        User storage user = users[_pid][msg.sender];
        // 更新质押池奖励信息
        _updatePoolReward(_pid);
        // 该质押池中应得奖励
        uint256 accumulated = (user.amount * pool.accTaskTokenPerShare) /
            ACC_PRECISION;
        // 应得奖励 大于 已记账的奖励
        if (accumulated > user.rewardDebt) {
            uint256 pendingNow = accumulated - user.rewardDebt;
            user.pending = user.pending + pendingNow;
        }

        uint256 toSend = user.pending;
        if (toSend == 0) {
            // 无奖励可领，仍需更新 rewardDebt（避免重复计入）
            user.rewardDebt = accumulated;
            return;
        }

        // 清空 pending 并更新 rewardDebt
        user.pending = 0;
        user.rewardDebt = accumulated;

        // 安全转账（如果合约余额不足则转尽量多）
        _safeTaskTokenTransfer(msg.sender, toSend);

        emit Claim(msg.sender, _pid, toSend);
    }

    /**
     * 请求解除质押（产生一个解锁请求）
     * - 将数量从用户与池中立即扣除（不可立刻提现），并产生一个 unlockBlock
     * - 用户可在 unlockBlock 之后调用提现函数（此合约目前只生成请求；你可以扩展增加领取解除的函数）
     * @param _pid 质押池Id
     * @param _amount 解除数量
     */
    function unstake(
        uint256 _pid,
        uint256 _amount
    ) public whenUnstakeNotPaused {
        require(_pid < pools.length, "invalid pid");
        Pool storage pool = pools[_pid];
        User storage user = users[_pid][msg.sender];

        require(_amount > 0, "amount must > 0");
        require(user.amount >= _amount, "not enough staked");
        // 更新池奖励并结算 pending
        _updatePoolReward(_pid);
        // 计算待领取金额
        uint256 accumulated = (user.amount * pool.accTaskTokenPerShare) /
            ACC_PRECISION;
        // 累计待领取金额
        if (accumulated > user.rewardDebt) {
            uint256 pendingNow = accumulated - user.rewardDebt;
            user.pending = user.pending + pendingNow;
        }
        // 扣减用户和池中的质押数量
        user.amount = user.amount - _amount;
        pool.stTokenAmount = pool.stTokenAmount - _amount;

        // 生成解除请求（unlockBlock = 当前区块 + unStakeLockedBlocks）
        uint256 unlockBlock = block.number + pool.unStakeLockedBlocks;
        user.requests.push(
            UnStakeRequest({amount: _amount, unlockBlock: unlockBlock})
        );

        // 更新 rewardDebt
        user.rewardDebt =
            (user.amount * pool.accTaskTokenPerShare) /
            ACC_PRECISION;

        emit UnstakeRequested(msg.sender, _pid, _amount, unlockBlock);
    }

    /**
     * 提取已经解锁的 unstake 请求中的代币（用户提取其之前请求解除的代币）
     * - 如果池是 ETH，则转 ETH；如果是 ERC20 则转 ERC20
     * - 支持按索引逐条提取（从 0 开始）
     * @param _pid 质押池ID
     * @param _requestIndex 请求下标
     */
    function withdrawUnstaked(
        uint256 _pid,
        uint256 _requestIndex
    ) external nonReentrant whenWithdrawUnstakedNotPaused {
        require(_pid < pools.length, "invalid pid");
        User storage user = users[_pid][msg.sender];
        require(_requestIndex < user.requests.length, "invalid request index");

        UnStakeRequest memory req = user.requests[_requestIndex];
        require(req.amount > 0, "already withdrawn or zero");
        require(block.number >= req.unlockBlock, "not unlocked yet");

        // 清空该请求（以免重复提取）——把最后一个替换到该位置然后 pop，减少 gas（但保留顺序可能改变）
        uint256 amount = req.amount;
        // 删除请求：用最后一个覆盖然后 pop
        uint256 lastIndex = user.requests.length - 1;
        if (_requestIndex != lastIndex) {
            user.requests[_requestIndex] = user.requests[lastIndex];
        }
        user.requests.pop();

        Pool storage pool = pools[_pid];

        // 代币返回给用户
        if (pool.stTokenAddress == address(0)) {
            // ETH
            (bool sent, ) = msg.sender.call{value: amount}("");
            require(sent, "ETH transfer failed");
        } else {
            IERC20Upgradeable(pool.stTokenAddress).safeTransfer(
                msg.sender,
                amount
            );
        }
    }

    /* ========== 池相关内部逻辑 ========== */

    /**
     *内部：更新某个池的奖励数据
     *  计算从 lastRewardBlock 到当前块之间产生的奖励（注意 endBlock 上限）
     * 按池权重分配，然后更新 accTaskTokenPerShare
     * @param _pid  质押池Id
     */
    function _updatePoolReward(uint256 _pid) internal {
        require(_pid < pools.length, "invalid pid");
        Pool storage pool = pools[_pid];
        uint256 currentBlock = block.number;
        if (currentBlock <= pool.lastRewardBlock) {
            return;
        }
        uint256 stTotal = pool.stTokenAmount;
        if (stTotal == 0 || pool.poolWeight == 0 || totalPoolWeight == 0) {
            // 仍然要更新 lastRewardBlock（否则后续会重复计算）
            pool.lastRewardBlock = currentBlock > endBlock
                ? endBlock
                : currentBlock;
            return;
        }

        // 计算可计入奖励的区块差（不超过 endBlock）
        uint256 toBlock = currentBlock;
        if (toBlock > endBlock) {
            toBlock = endBlock;
        }
        if (toBlock <= pool.lastRewardBlock) {
            pool.lastRewardBlock = toBlock;
            return;
        }

        // (block.number - pool.lastRewardBlock) * eachBlockRewardTokens / totalPoolWeight * pool.poolWeight / stTokenAmount
        // 当前区块减去 最后奖励区块 得到中间挖出区块
        // 挖出区块 乘以 每个区块奖励 得到 总奖励代币数
        uint256 multiplier = toBlock - pool.lastRewardBlock; // 区块数
        uint256 totalReward = multiplier * eachBlockRewardTokens;
        // 总奖励代币数  乘以 质押池权重 除以 总质押池权重 得到质押池可分配到的奖励代币数
        // 防止除 0
        uint256 poolReward = (totalReward * pool.poolWeight) / totalPoolWeight;

        // 奖励代币数 除以 质押池质押代币数 得到每个质押代币 可以得到多少奖励
        // acc per share 增量 = poolReward * ACC_PRECISION / stTotal
        uint256 accIncrease = (poolReward * ACC_PRECISION) / stTotal;

        // 更新 accTaskTokenPerShare
        pool.accTaskTokenPerShare = pool.accTaskTokenPerShare + accIncrease;

        // 更新 lastRewardBlock 到 toBlock（可能是 endBlock）
        pool.lastRewardBlock = toBlock;
    }

    /* ========== 视图函数 ========== */

    function poolLength() external view returns (uint256) {
        return pools.length;
    }

    /**
     * 查询某用户在某池的 pending（包括尚未写入 user.pending 的部分）
     * @param _pid 质押池ID
     * @param _user 用户地址
     */
    function pendingReward(
        uint256 _pid,
        address _user
    ) external view returns (uint256) {
        require(_pid < pools.length, "invalid pid");
        Pool storage pool = pools[_pid];
        User storage user = users[_pid][_user];

        uint256 accPerShare = pool.accTaskTokenPerShare;
        uint256 stTotal = pool.stTokenAmount;

        uint256 currentBlock = block.number;
        if (
            currentBlock > pool.lastRewardBlock &&
            stTotal != 0 &&
            pool.poolWeight != 0 &&
            totalPoolWeight != 0
        ) {
            uint256 toBlock = currentBlock > endBlock ? endBlock : currentBlock;
            if (toBlock > pool.lastRewardBlock) {
                uint256 multiplier = toBlock - pool.lastRewardBlock;
                uint256 totalReward = multiplier * eachBlockRewardTokens;
                uint256 poolReward = (totalReward * pool.poolWeight) /
                    totalPoolWeight;
                uint256 accIncrease = (poolReward * ACC_PRECISION) / stTotal;
                accPerShare = accPerShare + accIncrease;
            }
        }

        uint256 accumulated = (user.amount * accPerShare) / ACC_PRECISION;
        if (accumulated <= user.rewardDebt) {
            return user.pending;
        } else {
            return user.pending + (accumulated - user.rewardDebt);
        }
    }

    /**
     *  安全转账 taskToken（如果余额不足则转尽量多）
     * @param _to 转账方
     * @param _amount 数量
     */
    function _safeTaskTokenTransfer(address _to, uint256 _amount) internal {
        uint256 bal = taskToken.balanceOf(address(this));
        uint256 sendAmount = _amount;
        if (bal < _amount) {
            sendAmount = bal;
        }
        if (sendAmount > 0) {
            IERC20Upgradeable(address(taskToken)).safeTransfer(_to, sendAmount);
        }
    }

    /* ========== 升级控制 ========== */

    /**
     * UUPS 升级授权：只有拥有 UPGRADER_ROLE 的账号可以进行升级
     */
    function _authorizeUpgrade(address newImplementation)
        internal
        override
        onlyRole(UPGRADER_ROLE)
    {}

    /**
     * 重写 supportsInterface，以兼容 AccessControl
     */
    function supportsInterface(
        bytes4 interfaceId
    ) public view virtual override(AccessControlUpgradeable) returns (bool) {
        return super.supportsInterface(interfaceId);
    }

    // 接收 ETH（使合约能持有 ETH）
    // receive() external payable {}

    // fallback() external payable {}
}