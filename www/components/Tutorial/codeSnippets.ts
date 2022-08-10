export const installDependenciesCode = `
npm install merkletreejs ethers
`.trim()

const addressPlaceholderComments = [
  'your addresses will be filled in here',
  'when you click the "Copy code" button',
]
  .map((comment) => `  // ${comment}`)
  .join('\n')

export const merkleSetupCode = (addresses: string[]) =>
  `
import { MerkleTree } from 'merkletreejs';
import { utils } from 'ethers';

const addresses: string[] = [
${
  addresses.length > 0
    ? addresses.map((a) => `  '${a}',`).join('\n')
    : addressPlaceholderComments
}
];

const tree = new MerkleTree(
  addresses.map(utils.keccak256),
  utils.keccak256,
  { sortPairs: true },
);

console.log('the Merkle root is:', tree.getRoot().toString('hex'));

export function getMerkleRoot() {
  return tree.getRoot().toString('hex');
}

export function getMerkleProof(address: string) {
  const hashedAddress = utils.keccak256(address);
  return tree.getHexProof(hashedAddress);
}
`.trim()

export const passMerkleProofCode = `
import { getMerkleProof } from './merkle.ts';

// right before minting, get the Merkle proof for the current wallet
// const walletAddress = ...
const merkleProof = getMerkleProof(walletAddress);

// pass this to your contract
await myContract.mintAllowList(merkleProof);
`.trim()

export const ourLibraryCode = `
import { getMerkleProof } from 'merklefoolib';

// Get your Merkle root from the merkle.foo website and paste here
const merkleRoot = '0x0123456789abcdef0123456789abcdef01234567';

const walletAddress = '0x1000000000000000000000000000000000000000';
const proof = await getMerkleProof(merkleRoot, walletAddress);
`.trim()

export const nftMerkleProofCode = `
import {MerkleProof} from "openzeppelin/utils/cryptography/MerkleProof.sol";

contract NFTContract is ERC721 {
  bytes32 public merkleRoot;

  constructor(bytes32 _merkleRoot) {
    merkleRoot = _merkleRoot;
  }

  // Check the Merkle proof using this function
  function allowListed(address _wallet, bytes32[] calldata _proof)
      public
      view
      returns (bool)
    {
      return
          MerkleProof.verify(
              _proof,
              merkleRoot,
              keccak256(abi.encodePacked(_wallet))
          );
    }
  
  function mintAllowList(uint256 _tokenId, bytes32[] _proof) external {
    require(allowListed(msg.sender, _proof), "You are not on the allow list");
    _mint(msg.sender, _tokenId);
  }
}
`.trim()
