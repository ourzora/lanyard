import { isObject, isString } from 'assertate'

export type ErrorResponse = {
  error: true
  message: string
}

export function isErrorResponse(value: unknown): value is ErrorResponse {
  return isObject(value) && value.error === true && isString(value.message)
}

export type CreateMerkleResponse = {
  merkleRoot: string
  abiSig: string
}

export type TreeResponse = {
  unhashedLeaves: string[]
  leafCount: number
  leafTypeDescriptor: string | null
  packedEncoding: boolean | null
}
