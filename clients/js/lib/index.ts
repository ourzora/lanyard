import fetch from 'isomorphic-fetch'
import {
  CreateTreeRequest,
  CreateTreeResponse,
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
  return resp.json()
}

export const createTree = (
  req: CreateTreeRequest,
): Promise<CreateTreeResponse> => {
  return client('POST', 'tree', req)
}

export const getTree = (merkleRoot: string): Promise<GetTreeResponse> => {
  return client('GET', `tree?root=${encodeURIComponent(merkleRoot)}`)
}

export const getProof = (req: GetProofRequest): Promise<GetProofResponse> => {
  const { merkleRoot, unhashedLeaf } = req
  return client(
    'GET',
    `proof?root=${encodeURIComponent(merkleRoot)}&leaf=${encodeURIComponent(
      unhashedLeaf,
    )}`,
  )
}

export const getRoots = (proof: string[]): Promise<GetRootsResponse> => {
  const _p = encodeURIComponent(proof.join(','))
  return client('GET', `roots?proof=${_p}`)
}

export * from './types'
