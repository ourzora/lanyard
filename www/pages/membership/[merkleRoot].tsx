import classNames from 'classnames'
import About from 'components/About'
import Button from 'components/Button'
import PageTitle from 'components/PageTitle'
import { FailureContainer, SuccessContainer } from 'components/Container'
import {
  GetServerSideProps,
  GetServerSidePropsContext,
  InferGetServerSidePropsType,
} from 'next'
import { useCallback, useEffect, useRef } from 'react'
import { useVirtualizer } from '@tanstack/react-virtual'
import { useRouter } from 'next/router'
import { useQuery } from 'hooks/useQuery'
import { useWindowSize } from 'hooks/useWindowSize'
import { isAddressLike } from 'utils/address'
import { getMerkleTree } from 'utils/api'
import { isENSLike, resolveEnsDomain } from 'utils/ens'

type Props = {
  merkleRoot: string
  addresses: ReadonlyArray<string>
  lastQuery: string | null
  address: string | null
  error: Error | null
}

export default function MembershipPage({
  address,
  addresses,
  error,
  lastQuery,
  merkleRoot,
}: InferGetServerSidePropsType<typeof getServerSideProps>) {
  const router = useRouter()
  const { query, trimmedQuery, setQuery, isDisabled } = useQuery(lastQuery)
  const addressIndex = addresses.findIndex((leaf) => leaf === address)
  const isFound = addressIndex !== -1
  const { width } = useWindowSize()

  const handleQueryChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      setQuery(e.target.value)
    },
    [setQuery],
  )

  const handleSearchSubmit = useCallback(
    (e: React.FormEvent<HTMLFormElement>) => {
      e.preventDefault()
      router.replace(`/membership/${merkleRoot}?query=${trimmedQuery}`)
    },
    [merkleRoot, router, trimmedQuery],
  )

  return (
    <div className="flex flex-col gap-y-8">
      <PageTitle>Search the allowlist for an address or ENS name</PageTitle>
      <form className="flex flex-col gap-y-8" onSubmit={handleSearchSubmit}>
        <label htmlFor="query" className="sr-only">
          Address or ENS
        </label>
        <input
          className={classNames(
            'w-full p-4 font-mono border-2 rounded-lg',
            'min-h-fit border-neutral-200 focus:outline-none',
          )}
          value={query}
          onChange={handleQueryChange}
          placeholder="Enter an address or ENS name"
          id="query"
        />
        <Button
          className="w-full max-w-[30rem] sm:w-60 h-[66px]"
          disabled={isDisabled}
          label="Search"
        />
      </form>
      <div className="flex flex-col gap-y-4">
        <ResultContainer
          isFound={isFound}
          error={error}
          lastQuery={lastQuery}
        />
        <div>
          Showing {addresses.length} addresses from the{' '}
          <span className="font-mono break-words">{merkleRoot}</span> allowlist
        </div>
        {width !== undefined ? (
          <List
            rows={addresses}
            addressIndex={addressIndex}
            lanes={width < 1000 ? 1 : 2}
          />
        ) : null}
        <About />
      </div>
    </div>
  )
}

const ResultContainer = ({
  isFound,
  error,
  lastQuery,
}: {
  isFound: boolean
  error: Error | null
  lastQuery: string | null
}) => {
  if (lastQuery === null) {
    return null
  }

  if (isFound) {
    return (
      <SuccessContainer>{lastQuery} found in the allowlist</SuccessContainer>
    )
  }

  if (lastQuery) {
    return (
      <FailureContainer>
        {lastQuery} was not found in the allowlist
      </FailureContainer>
    )
  }

  if (error) {
    return <FailureContainer>{error.message}</FailureContainer>
  }

  return null
}

function List({
  addressIndex,
  lanes,
  rows,
}: {
  rows: ReadonlyArray<string>
  addressIndex: number
  lanes: number
}) {
  const parentRef = useRef<HTMLDivElement>(null)
  const isSingleLane = lanes === 1
  const ROW_HEIGHT = 50
  const rowVirtualizer = useVirtualizer({
    count: rows.length,
    getScrollElement: () => parentRef.current,
    estimateSize: () => ROW_HEIGHT,
    paddingStart: 8,
    paddingEnd: 8,
    lanes,
  })

  useEffect(() => {
    rowVirtualizer.measure()
  }, [lanes, rowVirtualizer])

  useEffect(() => {
    if (addressIndex !== -1) {
      rowVirtualizer.scrollToIndex(addressIndex, {
        align: 'center',
        behavior: 'smooth',
      })
    }
  }, [addressIndex, rowVirtualizer])

  return (
    <div
      ref={parentRef}
      className={classNames(
        'px-2 border-2 rounded-lg bg-neutral-50',
        'border-neutral w-full font-mono overflow-auto h-[300px]',
      )}
    >
      <ul
        className="relative grid w-full gap-x-4"
        style={{
          height: rowVirtualizer.getTotalSize(),
          gridTemplateColumns: isSingleLane ? '1fr' : '440px 1fr',
        }}
      >
        {rowVirtualizer.getVirtualItems().map((virtualRow) => {
          const isMatch = virtualRow.index === addressIndex
          return (
            <li
              key={virtualRow.key}
              data-index={virtualRow.index}
              className="absolute top-0"
              style={{
                gridColumn: isSingleLane ? 'auto' : (virtualRow.index % 2) + 1,
                transform: `translateY(${virtualRow.start}px)`,
              }}
            >
              <div
                className={classNames(
                  'px-4 py-2 border-2 rounded-lg w-fit',
                  {
                    'bg-white': !isMatch,
                  },
                  {
                    'border-brand bg-brand-light': isMatch,
                  },
                )}
              >
                {rows[virtualRow.index]}
              </div>
            </li>
          )
        })}
      </ul>
    </div>
  )
}

const getAddress = async (
  query: string | null,
): Promise<{ address: string | null; error: Error | null }> => {
  if (query === null) {
    return { address: null, error: null }
  }

  if (!isENSLike(query)) {
    return { address: query, error: null }
  }

  try {
    const address = await resolveEnsDomain(query)
    return { address, error: null }
  } catch (error) {
    if (error instanceof Error) {
      return { address: null, error }
    }
    return { address: null, error: Error('Error resolving ENS name') }
  }
}

const normalizeQueryParam = (query: string | string[] | undefined) => {
  return typeof query === 'string' ? query : null
}

export const getServerSideProps: GetServerSideProps<Props> = async (
  ctx: GetServerSidePropsContext,
) => {
  const merkleRoot = String(ctx.params?.merkleRoot)
  const lastQuery = normalizeQueryParam(ctx.query.query)
  const { address, error } = await getAddress(lastQuery)

  try {
    const tree = await getMerkleTree(merkleRoot)
    const addresses = tree.unhashedLeaves.filter(isAddressLike)
    return {
      props: {
        merkleRoot,
        addresses,
        lastQuery,
        error,
        address,
      },
    }
  } catch {
    return {
      notFound: true,
    }
  }
}
