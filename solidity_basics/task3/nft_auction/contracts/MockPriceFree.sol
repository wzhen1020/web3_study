// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.0;

import {AggregatorV3Interface} from "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";

contract MockPriceFeed is AggregatorV3Interface {
    // Mock implementation of the AggregatorV3Interface


// Latest price
    int256 private price;

    constructor(int256 _answer) {

        price = _answer;
    }

/**
 * @dev Returns the number of decimals the price feed uses. 
 */
    function decimals() external pure override returns (uint8) {
        return 8;
    }

    /// @dev Returns a description of the price feed.
    function description() external pure override returns (string memory) {
        return "Mock Price Feed";
    }

    /// @dev Returns the version of the price feed.
    function version() external pure override returns (uint256) {
        return 1;
    }

    /**
     * @dev Returns the round data for a specific round ID.
     */
    function getRoundData(
        uint80 _roundId
    )
        external
        view
        override
        returns (
            uint80 roundId,
            int256 answer,
            uint256 startedAt,
            uint256 updatedAt,
            uint80 answeredInRound
        )
    {
        return (_roundId, price, 0, 0, 0);
    }

    /// @dev Returns the latest round data.
    function latestRoundData()
        external
        view
        override
        returns (
            uint80 roundId,
            int256 answer,
            uint256 startedAt,
            uint256 updatedAt,
            uint80 answeredInRound
        )
    {
        return (0, price, 0, 0, 0);
    }
}
