// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

contract Math {
    // 加法
    function add(uint256 a, uint256 b) external pure returns (uint256) {
        return a + b;
    }

    // 减法
    function sub(uint256 a, uint256 b) external pure returns (uint256) {
        require(b <= a, "Underflow");
        return a - b;
    }
}
