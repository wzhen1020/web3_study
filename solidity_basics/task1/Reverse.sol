// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Reverse {


    function reverseStr(string calldata str) public pure returns(string memory) {

        uint length = bytes(str).length; 
        bytes memory reversed = new bytes(length); 
        for (uint i = 0; i < length; i++) {
            reversed[i] = bytes(str)[length - 1 - i]; 
        }
        return string(reversed); 

        
    }

}