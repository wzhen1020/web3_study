// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";


/**
 * @title MyNFT
 * 允许合约所有者铸造新的NFT
 */
contract MyNFT is ERC721URIStorage {

   uint256 private _tokenIds;

    constructor() ERC721("MyNFT", "MNFT") {}


    function mintNFT(
        address recipient,
        string memory tokenURI
    ) public returns (uint256) {
        _tokenIds++;

        uint256 newItemId = _tokenIds;
        _mint(recipient, newItemId);
        _setTokenURI(newItemId, tokenURI);

        return newItemId;
    }
}
