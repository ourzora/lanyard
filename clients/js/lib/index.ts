import fetch from 'isomorphic-fetch'
import {
  CreateTreeRequest,
  CreateTreeResponse,
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

export const CreateTree = (
  req: CreateTreeRequest,
): Promise<CreateTreeResponse> => {
  return client('POST', 'tree', req)
}

export const GetTree = (merkleRoot: string): Promise<GetTreeResponse> => {
  return client('GET', `tree?root=${encodeURIComponent(merkleRoot)}`)
}

export const GetProof = (
  root: string,
  unhashedLeaf: string,
): Promise<GetProofResponse> => {
  return client(
    'GET',
    `proof?root=${encodeURIComponent(root)}&leaf=${encodeURIComponent(
      unhashedLeaf,
    )}`,
  )
}

export const GetRoots = (proof: string): Promise<GetRootsResponse> => {
  return client('GET', `roots?proof=${encodeURIComponent(proof)}`)
}

export * from './types'
