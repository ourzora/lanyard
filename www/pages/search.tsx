import classNames from 'classnames'
import About from 'components/About'
import Button from 'components/Button'
import PageTitle from 'components/PageTitle'
import { useRouter } from 'next/router'
import { useQuery } from 'hooks/useQuery'

export default function SearchPage() {
  const router = useRouter()
  const { query, setQuery, isDisabled } = useQuery()

  const handleSubmit = () => {
    router.push(`/membership/${query}`)
  }

  return (
    <div className="flex flex-col">
      <PageTitle>Search for an existing allowlist</PageTitle>
      <div className="flex flex-col items-start gap-y-8">
        <div className="relative w-full">
          <label htmlFor="merkle-root" className="sr-only">
            Merkle Root
          </label>
          <input
            className={classNames(
              'w-full p-4 font-mono border-2 rounded-lg',
              'min-h-fit border-neutral-200 focus:outline-none',
            )}
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Enter the allowlist Merkle root"
            id="merkle-root"
          />
        </div>
        <Button
          disabled={isDisabled}
          onClick={handleSubmit}
          label="Search"
          className="w-full max-w-[30rem] sm:w-60 h-[66px]"
        />
      </div>
      <About />
    </div>
  )
}
