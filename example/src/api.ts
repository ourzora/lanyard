import { utils } from 'ethers'
import { MerkleTree } from 'merkletreejs'

const baseUrl = 'http://localhost:8080'

const encode = utils.defaultAbiCoder.encode.bind(utils.defaultAbiCoder)
const encodePacked = utils.solidityPack

const makeMerkleTree = (leafData: string[]) =>
  new MerkleTree(leafData.map(utils.keccak256), utils.keccak256, {
    sortPairs: true,
  })

const healthRes = await fetch(`${baseUrl}/health`, {
  method: 'GET',
  headers: {
    'Content-Type': 'application/json',
    'Accept-Encoding': 'gzip',
  },
})

const version = await healthRes.text()
console.log('api version', version)

const unhashedLeaves = [
  '0x0000000000000000000000000000000000000001',
  '0x0000000000000000000000000000000000000002',
  '0x0000000000000000000000000000000000000003',
  '0x0000000000000000000000000000000000000004',
  '0x0000000000000000000000000000000000000005',
]

// create a merkle tree

const createTreeRes = await fetch(`${baseUrl}/api/v1/tree`, {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Accept-Encoding': 'gzip',
  },
  body: JSON.stringify({ unhashedLeaves }),
})

const createdTree: { merkleRoot: string } = await createTreeRes.json()
console.log('merkle root', createdTree.merkleRoot)
console.log('local merkle root', makeMerkleTree(unhashedLeaves).getHexRoot())

// get a tree from a root

const getTreeRes = await fetch(
  `${baseUrl}/api/v1/tree?root=${createdTree.merkleRoot}`,
  {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      'Accept-Encoding': 'gzip',
    },
  },
)

const tree: {
  unhashedLeaves: string[]
  leafCount: number
} = await getTreeRes.json()

console.log('leaf count', tree.leafCount)

// get proof for a leaf

const proofRes = await fetch(
  `${baseUrl}/api/v1/proof?root=${createdTree.merkleRoot}&unhashedLeaf=${unhashedLeaves[0]}`,
  {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      'Accept-Encoding': 'gzip',
    },
  },
)

const proof: {
  proof: string[]
} = await proofRes.json()

console.log('proof', proof.proof)
console.log(
  'local proof',
  makeMerkleTree(unhashedLeaves).getHexProof(
    utils.keccak256(unhashedLeaves[0]),
  ),
)

// non-address leaf data

const num2Addr = (num: number) =>
  utils.hexlify(utils.zeroPad(utils.hexlify(num), 20))

const leafData = []

for (let i = 1; i <= 5; i++) {
  leafData.push([num2Addr(i), 2, utils.parseEther('0.01').toString()])
}

const encodedTreeRes = await fetch(`${baseUrl}/api/v1/tree`, {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Accept-Encoding': 'gzip',
  },
  body: JSON.stringify({
    unhashedLeaves: leafData.map((leafData) =>
      encode(['address', 'uint256', 'uint256'], leafData),
    ),
    leafTypeDescriptor: ['address', 'uint256', 'uint256'],
    packedEncoding: false,
  }),
})

const encodedTree: { merkleRoot: string } = await encodedTreeRes.json()

console.log('encoded tree', encodedTree.merkleRoot)
console.log(
  'local encoded merkle root',
  makeMerkleTree(
    leafData.map((leafData) =>
      encode(['address', 'uint256', 'uint256'], leafData),
    ),
  ).getHexRoot(),
)

const encodedPackedTreeRes = await fetch(`${baseUrl}/api/v1/tree`, {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Accept-Encoding': 'gzip',
  },
  body: JSON.stringify({
    unhashedLeaves: leafData.map((leafData) =>
      encodePacked(['address', 'uint256', 'uint256'], leafData),
    ),
    leafTypeDescriptor: ['address', 'uint256', 'uint256'],
    packedEncoding: true,
  }),
})

const encodedPackedTree: { merkleRoot: string } =
  await encodedPackedTreeRes.json()

console.log('encoded packed tree', encodedPackedTree.merkleRoot)
console.log(
  'local encoded packed merkle root',
  makeMerkleTree(
    leafData.map((leafData) =>
      encodePacked(['address', 'uint256', 'uint256'], leafData),
    ),
  ).getHexRoot(),
)
