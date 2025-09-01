// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;



import "@openzeppelin/contracts/access/Ownable.sol";

contract BeggingContract is Ownable {


    mapping(address account=> uint256 amount) private donateMapping;

    uint256 public totalDonations;

event Donation(address indexed account,uint256 amount);

    constructor () Ownable(msg.sender) {

    }


    function getDonation (address account) public view returns(uint256){

        return donateMapping[account];
    }

    function donate () external payable{
        donateMapping[msg.sender] += msg.value;
        totalDonations += msg.value;
        emit Donation(msg.sender,msg.value);
    }

    function withdraw () public payable onlyOwner{
       uint256 balance = address(this).balance;
        payable(owner()).transfer(balance);
    }
}

