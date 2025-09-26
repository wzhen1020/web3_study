// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

contract MathOptimized1 {
    function add(uint256 a, uint256 b) external pure returns (uint256) {
        return a + b;
    }

    // 使用 unchecked 避免 Solidity 自动溢出检查
    function sub(uint256 a, uint256 b) external pure returns (uint256) {
        unchecked {
            return a - b;
        }
    }
}
