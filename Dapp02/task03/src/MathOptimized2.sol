// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

contract MathOptimized2 {
    function add(uint256 a, uint256 b) external pure returns (uint256) {
        unchecked { return a + b; }
    }

    function sub(uint256 a, uint256 b) external pure returns (uint256) {
        unchecked { return a - b; }
    }
}
