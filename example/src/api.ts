import { utils } from 'ethers'
import { MerkleTree } from 'merkletreejs'

const baseUrl = process.env.API_URL ?? 'http://localhost:8080'

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
): Promise<{ proof: string[]; unhashedLeaf: string | null }> => {
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

/** When possible, use `unhashedLeaf` instead of `address` */
const getProofForIndexedAddress = async (
  merkleRoot: string,
  address: string,
): Promise<{ proof: string[]; unhashedLeaf: string | null }> => {
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

  return await proofRes.json()
}

const getRootFromProof = async (proof: string[]): Promise<string> => {
  const rootRes = await fetch(`${baseUrl}/api/v1/root?proof=${proof.join(',')}`)
  const resp: { root: string } = await rootRes.json()
  return resp.root
}

const encode = utils.defaultAbiCoder.encode.bind(utils.defaultAbiCoder)
const encodePacked = utils.solidityPack

const makeMerkleTree = (leafData: string[]) =>
  new MerkleTree(leafData.map(utils.keccak256), utils.keccak256, {
    sortPairs: true,
  })

const checkRootEquality = (remote: string, local: string) => {
  if (remote !== local) {
    throw new Error(`Remote root ${remote} does not match local root ${local}`)
  }
}

const strArrEqual = (a: string[], b: string[]) =>
  a.length === b.length && a.every((v, i) => v === b[i])

const checkProofEquality = (remote: string[], local: string[]) => {
  if (!strArrEqual(remote, local)) {
    throw new Error(
      `Remote proof ${remote} does not match local proof ${local}`,
    )
  }
}

// basic merkle tree

const unhashedLeaves = [
  '0x0000000000000000000000000000000000000001',
  '0x0000000000000000000000000000000000000002',
  '0x0000000000000000000000000000000000000003',
  '0x0000000000000000000000000000000000000004',
  '0x0000000000000000000000000000000000000005',
]

const { merkleRoot: basicMerkleRoot } = await createTree(unhashedLeaves)

console.log('basic merkle root', basicMerkleRoot)
checkRootEquality(basicMerkleRoot, makeMerkleTree(unhashedLeaves).getHexRoot())

const basicTree = await getTree(basicMerkleRoot)
console.log('basic leaf count', basicTree.leafCount)

const { proof: basicProof } = await getProofForUnhashedLeaf(
  basicMerkleRoot,
  unhashedLeaves[0],
)

console.log('proof', basicProof)
checkProofEquality(
  basicProof,
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
checkRootEquality(
  encodedMerkleRoot,
  makeMerkleTree(encodedLeafData).getHexRoot(),
)

const { proof: encodedProof } = await getProofForUnhashedLeaf(
  encodedMerkleRoot,
  encodedLeafData[0],
)

console.log('encoded proof', encodedProof)
checkProofEquality(
  encodedProof,
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
checkRootEquality(
  encodedPackedMerkleRoot,
  makeMerkleTree(encodedPackedLeafData).getHexRoot(),
)

const { proof: encodedPackedProof } = await getProofForUnhashedLeaf(
  encodedPackedMerkleRoot,
  encodedPackedLeafData[0],
)

console.log('encoded packed proof', encodedPackedProof)
checkProofEquality(
  encodedPackedProof,
  makeMerkleTree(encodedPackedLeafData).getHexProof(
    utils.keccak256(encodedPackedLeafData[0]),
  ),
)

const { proof: encodedPackedProofByAddress } = await getProofForIndexedAddress(
  encodedPackedMerkleRoot,
  num2Addr(1),
)
console.log(
  'encoded packed proof by indexed address',
  encodedPackedProofByAddress,
)
checkProofEquality(
  encodedPackedProofByAddress,
  makeMerkleTree(encodedPackedLeafData).getHexProof(
    utils.keccak256(encodedPackedLeafData[0]),
  ),
)

const root = await getRootFromProof(encodedPackedProof)
console.log('root from proof', root)
checkRootEquality(root, encodedPackedMerkleRoot)
