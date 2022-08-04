import { useAsync } from '@react-hook/async'
import classNames from 'classnames'
import Tutorial from 'components/Tutorial'
import { useCallback, useMemo, useState } from 'react'
import { isErrorResponse, createMerkleRoot } from 'utils/api'
import Button from 'components/Button'
import { parseAddressesFromText } from 'utils/addressParsing'

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

export default function CreatePage() {
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

  return (
    <div className="flex flex-col items-start gap-y-2">
      <textarea
        className={classNames(
          'w-full min-h-fit border-2',
          'focus:outline-none',
          'p-4 rounded-lg',
        )}
        style={{ minHeight: '8rem' }}
        value={addressInput}
        onChange={(e) => addressInputSet(e.target.value)}
        placeholder="Paste addresses here separated by commas, spaces, or new lines"
      />

      <div className="flex gap-4 items-center">
        <Button
          onClick={handleSubmit}
          label="Create Merkle Root"
          pending={status === 'loading'}
          disabled={parsedAddressesCount === 0}
        />

        {parsedAddressesCount > 0 && (
          <div>
            {parsedAddressesCount} address
            {parsedAddressesCount === 1 ? '' : 'es'} found
          </div>
        )}
      </div>

      {status === 'success' && (
        <>
          {merkleRoot !== undefined && <div>Merkle root: {merkleRoot}</div>}
          {errorResponse !== undefined && (
            <div>Error: {errorResponse.message}</div>
          )}
        </>
      )}

      {/* spacer */}
      <div className="mt-10" />

      <Tutorial addresses={parsedAddresses} />
    </div>
  )
}
