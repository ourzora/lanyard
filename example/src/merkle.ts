import { utils } from 'ethers'
import { MerkleTree } from 'merkletreejs'

const addresses: string[] = [
  '0x0000000000000000000000000000000000000001',
  '0x0000000000000000000000000000000000000002',
  '0x0000000000000000000000000000000000000003',
  '0x0000000000000000000000000000000000000004',
  '0x0000000000000000000000000000000000000005',
]

const tree = new MerkleTree(addresses.map(utils.keccak256), utils.keccak256, {
  sortPairs: true,
})

console.log('the Merkle root is:', tree.getRoot().toString('hex'))

export function getMerkleRoot() {
  return tree.getRoot().toString('hex')
}

export function getMerkleProof(address: string) {
  const hashedAddress = utils.keccak256(address)
  return tree.getHexProof(hashedAddress)
}
