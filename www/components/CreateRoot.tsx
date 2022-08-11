import { useAsync } from '@react-hook/async'
import classNames from 'classnames'
import { useCallback, useEffect, useMemo, useState } from 'react'
import { isErrorResponse, createMerkleRoot } from 'utils/api'
import Button from 'components/Button'
import { parseAddressesFromText } from 'utils/addressParsing'
import { useRouter } from 'next/router'

const useCreateMerkleRoot = () => {
  const [{ value, status }, create] = useAsync(
    async (addresses: string[]) => await createMerkleRoot(addresses),
  )

  const merkleRoot = useMemo(() => {
    if (value === undefined) return undefined
    if (isErrorResponse(value)) return undefined
    return value.merkleRoot
  }, [value])

  const error = useMemo(() => {
    if (isErrorResponse(value)) return value
    return undefined
  }, [value])

  return { merkleRoot, error, status, create }
}

export default function CreateRoot() {
  const {
    merkleRoot,
    error: errorResponse,
    status,
    create,
  } = useCreateMerkleRoot()
  const [addressInput, addressInputSet] = useState('')

  const handleSubmit = useCallback(() => {
    if (addressInput.trim().length === 0) {
      return
    }

    const addresses = parseAddressesFromText(addressInput)
    create(addresses)
  }, [addressInput, create])

  const parsedAddresses = useMemo(() => {
    if (addressInput.trim().length === 0) {
      return []
    }

    return parseAddressesFromText(addressInput)
  }, [addressInput])

  const parsedAddressesCount = useMemo(
    () => parsedAddresses.length,
    [parsedAddresses],
  )

  const router = useRouter()

  useEffect(() => {
    if (merkleRoot !== undefined) {
      router.push(`/tree/${merkleRoot}`)
    }
  }, [merkleRoot, router])

  return (
    <div className="flex flex-col items-start gap-y-8">
      <textarea
        className={classNames(
          'w-full min-h-fit border-2 border-neutral-200 resize-none',
          'focus:outline-none',
          'p-4 rounded-lg',
          'font-mono',
          'h-44',
        )}
        value={addressInput}
        onChange={(e) => addressInputSet(e.target.value)}
        placeholder="Paste addresses here, separated by commas, spaces or new lines"
      />

      <div className="flex flex-col sm:flex-row w-full gap-x-4 gap-y-2 items-center">
        <Button
          onClick={handleSubmit}
          label="Generate Merkle root"
          pending={status === 'loading'}
          disabled={parsedAddressesCount === 0}
          className="w-full max-w-[30rem] h-[66px]"
        />

        {parsedAddressesCount > 0 && (
          <div>
            {parsedAddressesCount} address
            {parsedAddressesCount === 1 ? '' : 'es'} found
          </div>
        )}
      </div>

      {status === 'success' && errorResponse !== undefined && (
        <div className="text-center w-full">Error: {errorResponse.message}</div>
      )}
    </div>
  )
}
