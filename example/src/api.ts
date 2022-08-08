import { utils } from 'ethers'
import { MerkleTree } from 'merkletreejs'

const baseUrl = 'http://localhost:8080'

const createTree = async (
  unhashedLeaves: string[],
  leafTypeDescriptor?: string[],
  packedEncoding?: boolean,
): Promise<{ merkleRoot: string }> => {
  const encodedTreeRes = await fetch(`${baseUrl}/api/v1/tree`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Accept-Encoding': 'gzip',
    },
    body: JSON.stringify({
      unhashedLeaves,
      leafTypeDescriptor,
      packedEncoding,
    }),
  })

  return await encodedTreeRes.json()
}

const getTree = async (
  merkleRoot: string,
): Promise<{
  unhashedLeaves: string[]
  leafCount: number
  leafTypeDescriptor: string[] | null
  packedEncoding: boolean | null
}> => {
  const getTreeRes = await fetch(`${baseUrl}/api/v1/tree?root=${merkleRoot}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      'Accept-Encoding': 'gzip',
    },
  })

  return await getTreeRes.json()
}

const getProofForUnhashedLeaf = async (
  merkleRoot: string,
  unhashedLeaf: string,
): Promise<{ proof: string[] }> => {
  const proofRes = await fetch(
    `${baseUrl}/api/v1/proof?root=${merkleRoot}&unhashedLeaf=${unhashedLeaf}`,
    {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Accept-Encoding': 'gzip',
      },
    },
  )

  return await proofRes.json()
}

/** Using this endpoint is discouraged. When possible, pass `unhashedLeaf` instead */
const getProofForAddress = async (
  merkleRoot: string,
  address: string,
): Promise<{ proof: string[] }> => {
  const proofRes = await fetch(
    `${baseUrl}/api/v1/proof?root=${merkleRoot}&address=${address}`,
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
  return proof
}

const encode = utils.defaultAbiCoder.encode.bind(utils.defaultAbiCoder)
const encodePacked = utils.solidityPack

const makeMerkleTree = (leafData: string[]) =>
  new MerkleTree(leafData.map(utils.keccak256), utils.keccak256, {
    sortPairs: true,
  })

// health check

const healthRes = await fetch(`${baseUrl}/health`, {
  method: 'GET',
  headers: {
    'Content-Type': 'application/json',
    'Accept-Encoding': 'gzip',
  },
})

const version = await healthRes.text()
console.log('api version', version)

// basic merkle tree

const unhashedLeaves = [
  '0x0000000000000000000000000000000000000001',
  '0x0000000000000000000000000000000000000002',
  '0x0000000000000000000000000000000000000003',
  '0x0000000000000000000000000000000000000004',
  '0x0000000000000000000000000000000000000005',
]

const { merkleRoot: basicMerkleRoot } = await createTree(unhashedLeaves)

console.log('merkle root', basicMerkleRoot)
console.log('local merkle root', makeMerkleTree(unhashedLeaves).getHexRoot())

const { leafCount: basicLeafCount } = await getTree(basicMerkleRoot)
console.log('leaf count', basicLeafCount)

const { proof: basicProof } = await getProofForUnhashedLeaf(
  basicMerkleRoot,
  unhashedLeaves[0],
)

console.log('proof', basicProof)
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

// encoded data

const encodedLeafData = leafData.map((leafData) =>
  encode(['address', 'uint256', 'uint256'], leafData),
)

const { merkleRoot: encodedMerkleRoot } = await createTree(
  encodedLeafData,
  ['address', 'uint256', 'uint256'],
  false,
)

console.log('encoded tree', encodedMerkleRoot)
console.log(
  'local encoded merkle root',
  makeMerkleTree(encodedLeafData).getHexRoot(),
)

const { proof: encodedProof } = await getProofForUnhashedLeaf(
  encodedMerkleRoot,
  encodedLeafData[0],
)

console.log('encoded proof', encodedProof)
console.log(
  'local encoded proof',
  makeMerkleTree(encodedLeafData).getHexProof(
    utils.keccak256(encodedLeafData[0]),
  ),
)

// packed data

const encodedPackedLeafData = leafData.map((leafData) =>
  encodePacked(['address', 'uint256', 'uint256'], leafData),
)

const { merkleRoot: encodedPackedMerkleRoot } = await createTree(
  encodedPackedLeafData,
  ['address', 'uint256', 'uint256'],
  true,
)

console.log('encoded packed tree', encodedPackedMerkleRoot)
console.log(
  'local encoded packed merkle root',
  makeMerkleTree(encodedPackedLeafData).getHexRoot(),
)

const { proof: encodedPackedProof } = await getProofForUnhashedLeaf(
  encodedPackedMerkleRoot,
  encodedPackedLeafData[0],
)

console.log('encoded packed proof', encodedPackedProof)
console.log(
  'local encoded packed proof',
  makeMerkleTree(encodedPackedLeafData).getHexProof(
    utils.keccak256(encodedPackedLeafData[0]),
  ),
)

const { proof: encodedPackedProofByAddress } = await getProofForAddress(
  encodedPackedMerkleRoot,
  num2Addr(1),
)
console.log('encoded packed proof by address', encodedPackedProofByAddress)
