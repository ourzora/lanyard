import fetch from 'isomorphic-fetch'
import {
  CreateTreeRequest,
  CreateTreeResponse,
  GetProofByAddress,
  GetProofByLeaf,
  GetProofRequest,
  GetProofResponse,
  GetRootsResponse,
  GetTreeResponse,
} from './types'
const baseUrl = 'https://lanyard.org/api/v1/'

const client = async (method: 'GET' | 'POST', path: string, data?: any) => {
  const opts: {
    method: string
    body?: string
    headers?: any
  } = {
    method: method,
  }

  if (method !== 'GET') {
    opts.body = JSON.stringify(data)
    opts.headers = {
      'Content-Type': 'application/json',
    }
  }

  const resp = await fetch(baseUrl + path, opts)
  if (resp.status > 299) {
    if (method === 'GET') {
      return null
    }
    throw new Error(
      `Request failed with status ${resp.status}: ${await resp.json()}`,
    )
  }

  return resp.json()
}

export const createTree = (
  req: CreateTreeRequest,
): Promise<CreateTreeResponse> => client('POST', 'tree', req)

export const getTree = (merkleRoot: string): Promise<GetTreeResponse> =>
  client('GET', `tree?root=${encodeURIComponent(merkleRoot)}`)

export const getProof = (req: GetProofRequest): Promise<GetProofResponse> => {
  let url = `proof?root=${encodeURIComponent(req.merkleRoot)}`

  if ((req as GetProofByAddress).address) {
    url += `&address=${encodeURIComponent((req as GetProofByAddress).address)}`
  } else {
    url += `&unhashedLeaf=${encodeURIComponent(
      (req as GetProofByLeaf).unhashedLeaf,
    )}`
  }

  return client('GET', url)
}

export const getRoots = (proof: string[]): Promise<GetRootsResponse> =>
  client('GET', `roots?proof=${encodeURIComponent(proof.join(','))}`)

export * from './types'
