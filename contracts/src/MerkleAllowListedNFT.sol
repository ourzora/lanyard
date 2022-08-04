// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.15;

import "solmate/tokens/ERC721.sol";
import "openzeppelin/utils/cryptography/MerkleProof.sol";

contract MerkleAllowListedNFT is ERC721 {
    bytes32 public merkleRoot;

    constructor(bytes32 _merkleRoot) ERC721("MerkleAllowListedNFT", "MALNFT") {
        merkleRoot = _merkleRoot;
    }

    mapping(address => bool) allowListAddressMinted;

    function mint(uint256 tokenId, bytes32[] calldata proof) external {
        require(allowListed(msg.sender, proof), "not on the allow list");
        require(allowListAddressMinted[msg.sender] == false, "already minted");

        allowListAddressMinted[msg.sender] = true;
        _mint(msg.sender, tokenId);
    }

    function allowListed(address _address, bytes32[] calldata _proof)
        public
        view
        returns (bool)
    {
        return
            MerkleProof.verify(
                _proof,
                merkleRoot,
                keccak256(abi.encodePacked(_address))
            );
    }

    function tokenURI(uint256) public pure override returns (string memory) {
        revert();
    }
}
