// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Merge {

function mergeArray(uint[] memory nums1,uint[] memory nums2) public pure returns (uint[] memory ){

 uint numsLength1 = nums1.length;
  uint numsLength2 = nums2.length;

   uint newLength = numsLength1+numsLength2;
 uint[] memory newNums = new uint[](newLength);

uint numsIndex1 = 0;
uint numsIndex2 = 0;
  for (uint i = 0; i < newLength; i++) 
  {
    
    if(numsIndex1 >= numsLength1){
         newNums[i] = nums2[numsIndex2];
         numsIndex2++;
         continue;
    }

        if(numsIndex2 >= numsLength2){
          newNums[i] = nums1[numsIndex1];
        numsIndex1++;
          continue;
    }

    if(nums1[numsIndex1] < nums2[numsIndex2]){

        newNums[i] = nums1[numsIndex1];
        numsIndex1++;
    }else{
        newNums[i] = nums2[numsIndex2];
        numsIndex2++;
    }


  }

  return newNums;
}

}