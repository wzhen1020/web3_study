// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

import "forge-std/Test.sol";
import "../src/Math.sol";
import "../src/MathOptimized1.sol";
import "../src/MathOptimized2.sol";

contract MathTest is Test {
    // Math math;
    // MathOptimized1 math;
    MathOptimized2 math;

    function setUp() public {
        // math = new Math();
        // math = new MathOptimized1();
        math = new MathOptimized2();
    }

    function testAdd() public {
        uint256 a = 10;
        uint256 b = 20;

        uint256 gasBefore = gasleft();
        uint256 result = math.add(a, b);
        uint256 gasAfter = gasleft();

        emit log_named_uint("Add Gas Used", gasBefore - gasAfter);
        assertEq(result, 30);
    }

    function testSub() public {
        uint256 a = 50;
        uint256 b = 20;

        uint256 gasBefore = gasleft();
        uint256 result = math.sub(a, b);
        uint256 gasAfter = gasleft();

        emit log_named_uint("Sub Gas Used", gasBefore - gasAfter);
        assertEq(result, 30);
    }
}
