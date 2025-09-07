// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @title Simple Counter - 简单计数器合约（带返回值）
/// @notice 提供增/减/重置/设置计数器的基本操作，并返回最新值
contract Counter {
    // 当前计数值
    uint256 public count;

    // 合约所有者地址
    address public owner;

    // 当计数改变时触发事件
    event CountChanged(address indexed by, uint256 newCount);

    /// @dev 构造函数：部署合约时设置所有者为 msg.sender
    constructor() {
        owner = msg.sender;
        count = 0;
    }

    /// @notice 增加计数器
    /// @param delta 增加多少（必须 > 0）
    /// @return 新的计数值
    function increment(uint256 delta) external returns (uint256) {
        require(delta > 0, "delta must be > 0");
        count += delta;
        emit CountChanged(msg.sender, count);
        return count;
    }

    /// @notice 减少计数器
    /// @param delta 减少多少（必须 > 0 且 <= count）
    /// @return 新的计数值
    function decrement(uint256 delta) external returns (uint256) {
        require(delta > 0, "delta must be > 0");
        require(delta <= count, "delta exceeds count");
        count -= delta;
        emit CountChanged(msg.sender, count);
        return count;
    }

    /// @notice 重置计数器为 0（仅限所有者）
    /// @return 新的计数值（总是 0）
    function reset() external returns (uint256) {
        count = 0;
        emit CountChanged(msg.sender, count);
        return count;
    }

    /// @notice 获取当前计数值（任何人都可以调用）
    /// @return 当前计数值
    function getCount() external view returns (uint256) {
        return count;
    }
}
