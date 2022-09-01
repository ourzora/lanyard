import { useAsync } from '@react-hook/async'
import classNames from 'classnames'
import { useCallback, useEffect, useMemo, useState } from 'react'
import { isErrorResponse, createMerkleRoot } from 'utils/api'
import Button from 'components/Button'
import { parseAddressesFromText, prepareAddresses } from 'utils/addressParsing'
import { useRouter } from 'next/router'
import { randomBytes } from 'crypto'
import { resolveEnsDomains } from 'utils/ens'

const useCreateMerkleRoot = () => {
  const [ensMap, setEnsMap] = useState<Record<string, string>>({})

  const [{ value, status, error: reqError }, create] = useAsync(
    async (addressesOrENSNames: string[]) => {
      let prepared = prepareAddresses(addressesOrENSNames, ensMap)

      if (prepared.unresolvedEnsNames.length > 0) {
        const ensAddresses = await resolveEnsDomains(
          prepared.unresolvedEnsNames,
        )

        setEnsMap((prev) => ({
          ...prev,
          ...ensAddresses,
        }))

        prepared = prepareAddresses(addressesOrENSNames, {
          ...ensMap,
          ...ensAddresses,
        })

        if (prepared.unresolvedEnsNames.length > 0) {
          throw new Error(`Could not resolve all ENS names`)
        }
      }

      if (prepared.addresses.length !== prepared.dedupedAddresses.length) {
        return (
          await Promise.all([
            createMerkleRoot(prepared.dedupedAddresses),
            createMerkleRoot(prepared.addresses),
          ])
        )[0]
      }

      return await createMerkleRoot(prepared.dedupedAddresses)
    },
  )

  const merkleRoot = useMemo(() => {
    if (value === undefined) return undefined
    if (isErrorResponse(value)) return undefined
    return value.merkleRoot
  }, [value])

  const error = useMemo(() => {
    if (isErrorResponse(value)) return value
    if (reqError !== undefined)
      return { error: true, message: reqError.message }
    return undefined
  }, [value, reqError])

  return { merkleRoot, error, status, create, parsedEnsNames: ensMap }
}

const randomAddress = () => `0x${randomBytes(20).toString('hex')}`

export default function CreateRoot() {
  const {
    merkleRoot,
    error: errorResponse,
    status,
    create,
    parsedEnsNames,
  } = useCreateMerkleRoot()
  const [addressInput, addressInputSet] = useState('')

  const handleSubmit = useCallback(() => {
    if (addressInput.trim().length === 0) {
      return
    }

    const addresses = parseAddressesFromText(addressInput)
    create(addresses)
  }, [addressInput, create])

  const parsedAddresses = useMemo(
    () => parseAddressesFromText(addressInput),
    [addressInput],
  )

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

  const handleLoadExample = useCallback(() => {
    const addresses = Array(Math.floor(Math.random() * 16) + 3)
      .fill(0)
      .map(() => randomAddress())

    addressInputSet(addresses.join('\n'))
  }, [])

  const handleRemoveInvalidENSNames = useCallback(() => {
    const addresses = parsedAddresses
      .filter((address) => {
        return (
          !address.includes('.') ||
          parsedEnsNames[address.toLowerCase()] !== undefined
        )
      })
      .join('\n')
    addressInputSet(addresses)
  }, [parsedAddresses, parsedEnsNames])

  const showRemoveInvalidENSNames = useMemo(
    // show the button if the error message contains `Could not resolve all ENS names`
    () =>
      errorResponse?.message?.includes('Could not resolve all ENS names') ??
      false,
    [errorResponse?.message],
  )

  const buttonPending = status === 'loading' || merkleRoot !== undefined

  return (
    <div className="flex flex-col items-start gap-y-8">
      <div className="relative w-full">
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
          placeholder="Paste addresses or ENS names here, separated by commas, spaces or new lines"
        />
      </div>

      <div className="flex flex-col sm:flex-row w-full gap-x-4 gap-y-2 items-center">
        <Button
          onClick={handleSubmit}
          label="Generate Merkle root"
          pending={buttonPending}
          disabled={parsedAddressesCount === 0}
          className="w-full max-w-[30rem] h-[66px]"
        />

        <Button
          disabled={parsedAddressesCount > 0}
          onClick={handleLoadExample}
          label="Load example"
          className="w-full max-w-[30rem] sm:w-60 h-[66px]"
        />

        {parsedAddressesCount > 0 && (
          <div>
            {parsedAddressesCount} address
            {parsedAddressesCount === 1 ? '' : 'es'} found
          </div>
        )}
      </div>

      {errorResponse !== undefined && (
        <div className="text-center sm:text-left w-full">
          Error: {errorResponse.message}
          {showRemoveInvalidENSNames && (
            <>
              {' '}
              <button
                className="underline"
                onClick={handleRemoveInvalidENSNames}
                type="button"
              >
                Remove invalid ENS names
              </button>
            </>
          )}
        </div>
      )}
    </div>
  )
}
