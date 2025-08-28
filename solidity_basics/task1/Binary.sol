// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Binary {



    function  search (int[] memory nums, int target) public pure returns (int num) {
        uint len = nums.length;
        if(len == 0){
            return -1;
        }
        uint right = len -1;
        uint left =  0;
        int value = 0;
        uint index = 0;
        while (left <= right) 
        {
            index = (right - left)/2 + left;
           value = nums[index];

            if(right == 0){
                break;
            }
            if(value > target){
                right = index -1;
                continue;
            }

            if(value < target){
                left = index +1;
                continue;
            }

            if(value == target){
                return int(index);
            }
        }
        return -1;
          
    }
}
