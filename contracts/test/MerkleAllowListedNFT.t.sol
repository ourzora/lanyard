// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.15;

import "forge-std/Test.sol";

import "../src/MerkleAllowListedNFT.sol";

contract MerkleAllowListedNFTTest is Test {
    function testDeploy() public {
        bytes32 root = 0x0000000000000000000000000000000000000000000000000000000000000001;
        MerkleAllowListedNFT malnft = new MerkleAllowListedNFT(root);
        assertEq(malnft.merkleRoot(), root);
    }

    function testMintNotAllowed() public {
        bytes32 root = 0x0000000000000000000000000000000000000000000000000000000000000001;
        MerkleAllowListedNFT malnft = new MerkleAllowListedNFT(root);

        bytes32[] memory proof = new bytes32[](2);
        proof[
            0
        ] = 0x1000000000000000000000000000000000000000000000000000000000000000;
        proof[
            1
        ] = 0x2000000000000000000000000000000000000000000000000000000000000000;

        vm.expectRevert("not on the allow list");
        malnft.mint(0, proof);
    }

    function testEmptyProofFails() public {
        bytes32 root = 0x0000000000000000000000000000000000000000000000000000000000000001;
        MerkleAllowListedNFT malnft = new MerkleAllowListedNFT(root);

        bytes32[] memory proof = new bytes32[](0);

        vm.expectRevert("not on the allow list");
        malnft.mint(0, proof);
    }

    function testOneAddressEmptyProof() public {
        bytes32 root = 0x1468288056310c82aa4c01a7e12a10f8111a0560e72b700555479031b86c357d;
        MerkleAllowListedNFT malnft = new MerkleAllowListedNFT(root);

        bytes32[] memory emptyProof = new bytes32[](0);

        address allowed = address(1);
        vm.prank(allowed);
        malnft.mint(0, emptyProof);

        address notAllowed = address(2);
        vm.prank(notAllowed);
        vm.expectRevert("not on the allow list");
        malnft.mint(1, emptyProof);
    }

    function testOnlyMintOnce() public {
        bytes32 root = 0x1468288056310c82aa4c01a7e12a10f8111a0560e72b700555479031b86c357d;
        MerkleAllowListedNFT malnft = new MerkleAllowListedNFT(root);

        bytes32[] memory emptyProof = new bytes32[](0);

        address allowed = address(1);
        vm.prank(allowed);
        malnft.mint(0, emptyProof);

        vm.expectRevert("already minted");
        vm.prank(allowed);
        malnft.mint(1, emptyProof);
    }

    function testTwoAddresses() public {
        bytes32 root = 0xf95c14e6953c95195639e8266ab1a6850864d59a829da9f9b13602ee522f672b;
        MerkleAllowListedNFT malnft = new MerkleAllowListedNFT(root);

        bytes32[] memory proof = new bytes32[](1);

        proof[
            0
        ] = 0xd52688a8f926c816ca1e079067caba944f158e764817b83fc43594370ca9cf62;
        vm.prank(address(1));
        malnft.mint(0, proof);

        proof[
            0
        ] = 0x1468288056310c82aa4c01a7e12a10f8111a0560e72b700555479031b86c357d;
        vm.prank(address(2));
        malnft.mint(1, proof);
    }

    function testThreeAddresses() public {
        bytes32 root = 0x344510bd0c324c3912b13373e89df42d1b50450e9764a454b2aa6e2968a4578a;
        MerkleAllowListedNFT malnft = new MerkleAllowListedNFT(root);

        bytes32[] memory firstProof = new bytes32[](2);
        firstProof[
            0
        ] = 0xd52688a8f926c816ca1e079067caba944f158e764817b83fc43594370ca9cf62;
        firstProof[
            1
        ] = 0x5b70e80538acdabd6137353b0f9d8d149f4dba91e8be2e7946e409bfdbe685b9;
        vm.prank(address(1));
        malnft.mint(0, firstProof);

        bytes32[] memory secondProof = new bytes32[](2);
        secondProof[
            0
        ] = 0x1468288056310c82aa4c01a7e12a10f8111a0560e72b700555479031b86c357d;
        secondProof[
            1
        ] = 0x5b70e80538acdabd6137353b0f9d8d149f4dba91e8be2e7946e409bfdbe685b9;
        vm.prank(address(2));
        malnft.mint(1, secondProof);

        bytes32[] memory thirdProof = new bytes32[](1);
        thirdProof[
            0
        ] = 0xf95c14e6953c95195639e8266ab1a6850864d59a829da9f9b13602ee522f672b;
        vm.prank(address(3));
        malnft.mint(2, thirdProof);
    }

    function testFourAddresses() public {
        bytes32 root = 0x5071e19149cc9b870c816e671bc5db717d1d99185c17b082af957a0a93888dd9;
        MerkleAllowListedNFT malnft = new MerkleAllowListedNFT(root);

        bytes32[] memory firstProof = new bytes32[](2);
        firstProof[
            0
        ] = 0xd52688a8f926c816ca1e079067caba944f158e764817b83fc43594370ca9cf62;
        firstProof[
            1
        ] = 0x735c77c52a2b69afcd4e13c0a6ece7e4ccdf2b379d39417e21efe8cd10b5ff1b;
        vm.prank(address(1));
        malnft.mint(0, firstProof);

        bytes32[] memory secondProof = new bytes32[](2);
        secondProof[
            0
        ] = 0x1468288056310c82aa4c01a7e12a10f8111a0560e72b700555479031b86c357d;
        secondProof[
            1
        ] = 0x735c77c52a2b69afcd4e13c0a6ece7e4ccdf2b379d39417e21efe8cd10b5ff1b;
        vm.prank(address(2));
        malnft.mint(1, secondProof);

        bytes32[] memory thirdProof = new bytes32[](2);
        thirdProof[
            0
        ] = 0xa876da518a393dbd067dc72abfa08d475ed6447fca96d92ec3f9e7eba503ca61;
        thirdProof[
            1
        ] = 0xf95c14e6953c95195639e8266ab1a6850864d59a829da9f9b13602ee522f672b;
        vm.prank(address(3));
        malnft.mint(2, thirdProof);

        bytes32[] memory fourthProof = new bytes32[](2);
        fourthProof[
            0
        ] = 0x5b70e80538acdabd6137353b0f9d8d149f4dba91e8be2e7946e409bfdbe685b9;
        fourthProof[
            1
        ] = 0xf95c14e6953c95195639e8266ab1a6850864d59a829da9f9b13602ee522f672b;
        vm.prank(address(4));
        malnft.mint(3, fourthProof);
    }

    function testFiveAddresses() public {
        bytes32 root = 0xa7a6b1cb6d12308ec4818baac3413fafa9e8b52cdcd79252fa9e29c9a2f8aff1;
        MerkleAllowListedNFT malnft = new MerkleAllowListedNFT(root);

        bytes32[] memory firstProof = new bytes32[](3);
        firstProof[
            0
        ] = 0xd52688a8f926c816ca1e079067caba944f158e764817b83fc43594370ca9cf62;
        firstProof[
            1
        ] = 0x735c77c52a2b69afcd4e13c0a6ece7e4ccdf2b379d39417e21efe8cd10b5ff1b;
        firstProof[
            2
        ] = 0x421df1fa259221d02aa4956eb0d35ace318ca24c0a33a64c1af96cf67cf245b6;
        vm.prank(address(1));
        malnft.mint(0, firstProof);

        bytes32[] memory secondProof = new bytes32[](3);
        secondProof[
            0
        ] = 0x1468288056310c82aa4c01a7e12a10f8111a0560e72b700555479031b86c357d;
        secondProof[
            1
        ] = 0x735c77c52a2b69afcd4e13c0a6ece7e4ccdf2b379d39417e21efe8cd10b5ff1b;
        secondProof[
            2
        ] = 0x421df1fa259221d02aa4956eb0d35ace318ca24c0a33a64c1af96cf67cf245b6;
        vm.prank(address(2));
        malnft.mint(1, secondProof);

        bytes32[] memory thirdProof = new bytes32[](3);
        thirdProof[
            0
        ] = 0xa876da518a393dbd067dc72abfa08d475ed6447fca96d92ec3f9e7eba503ca61;
        thirdProof[
            1
        ] = 0xf95c14e6953c95195639e8266ab1a6850864d59a829da9f9b13602ee522f672b;
        thirdProof[
            2
        ] = 0x421df1fa259221d02aa4956eb0d35ace318ca24c0a33a64c1af96cf67cf245b6;
        vm.prank(address(3));
        malnft.mint(2, thirdProof);

        bytes32[] memory fourthProof = new bytes32[](3);
        fourthProof[
            0
        ] = 0x5b70e80538acdabd6137353b0f9d8d149f4dba91e8be2e7946e409bfdbe685b9;
        fourthProof[
            1
        ] = 0xf95c14e6953c95195639e8266ab1a6850864d59a829da9f9b13602ee522f672b;
        fourthProof[
            2
        ] = 0x421df1fa259221d02aa4956eb0d35ace318ca24c0a33a64c1af96cf67cf245b6;
        vm.prank(address(4));
        malnft.mint(3, fourthProof);

        bytes32[] memory fifthProof = new bytes32[](1);
        fifthProof[
            0
        ] = 0x5071e19149cc9b870c816e671bc5db717d1d99185c17b082af957a0a93888dd9;
        vm.prank(address(5));
        malnft.mint(4, fifthProof);
    }
}
