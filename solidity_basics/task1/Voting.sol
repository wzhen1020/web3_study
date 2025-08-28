// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;



contract Voting{

mapping(string name => uint nums) private ticketMap;
string[] private candidates;

function vote (string calldata name) public {
    ticketMap[name]+=1;
        candidates.push(name);
}

function getVotes(string memory name)public view returns (uint256){

    return ticketMap[name];
}

function resetVotes() public{

    for (uint256 i; i < candidates.length;i++) 
    {
        ticketMap[candidates[i]] = 0;
    }

}

}




