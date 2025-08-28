// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Solution {
    // 使用 bytes1 作为键类型，因为罗马数字是单个字符
    mapping(bytes1 => uint256) private symbolValues;

    constructor() {
        // 初始化罗马数字映射
        symbolValues['I'] = 1;
        symbolValues['V'] = 5;
        symbolValues['X'] = 10;
        symbolValues['L'] = 50;
        symbolValues['C'] = 100;
        symbolValues['D'] = 500;
        symbolValues['M'] = 1000;
    }

    function romanToInt(string memory str) public view returns (uint256) {
        // 将字符串转换为字节数组以便处理
        bytes memory strBytes = bytes(str);
        uint256 length = strBytes.length;
        
        // 如果字符串为空，返回0
        if (length == 0) {
            return 0;
        }

        uint256 sum = 0;
        uint256 prevValue = 0;

        // 从右向左遍历字符串（罗马数字通常从右向左处理更高效）
        for (uint256 i = length; i > 0; i--) {
            bytes1 currentChar = strBytes[i - 1];
            uint256 currentValue = symbolValues[currentChar];
            
            // 如果当前值小于前一个值，则减去当前值
            if (currentValue < prevValue) {
                sum -= currentValue;
            } else {
                sum += currentValue;
            }
            
            prevValue = currentValue;
        }

        return sum;
    }
   string[]  thousands = ["", "M", "MM", "MMM"];
        string[]  hundreds  =  ["", "C", "CC", "CCC", "CD", "D", "DC", "DCC", "DCCC", "CM"];
        string[]  tens      =  ["", "X", "XX","XXX", "XL", "L", "LX", "LXX", "LXXX", "XC"];
        string[]  ones      =  ["", "I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX"];

    function intToRoman (uint num) public view returns ( string memory ){

        return string(abi.encodePacked(thousands[num / 1000],
                    hundreds[(num % 1000) / 100],
                    tens[(num % 100) / 10],
                    ones[num % 10]));

            }
}