export interface CreateTreeRequest {
  unhashedLeaves: string[]
  leafTypeDescriptor?: string[]
  packedEncoding?: boolean
}

export interface CreateTreeResponse {
  merkleRoot: string
}

export type GetTreeResponse = {
  unhashedLeaves: string[]
  leafCount: number
  leafTypeDescriptor: string[] | null
  packedEncoding: boolean | null
} | null

export type GetProofByLeaf = {
  merkleRoot: string
  unhashedLeaf: string
} | null

export type GetProofByAddress = {
  merkleRoot: string
  address: string
} | null

export type GetProofRequest = GetProofByLeaf | GetProofByAddress

export type GetProofResponse = {
  proof: string[]
  unhashedLeaf: string
} | null

export interface GetRootsResponse {
  roots: string[]
}
