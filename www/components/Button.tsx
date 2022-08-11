import classNames from 'classnames'
import React from 'react'

type Props = {
  onClick: () => void
  label: string
  disabled?: boolean
  pending?: boolean
  className?: string
}

function Button({
  onClick,
  label,
  disabled = false,
  pending = false,
  className,
}: Props) {
  return (
    <button
      onClick={onClick}
      className={classNames(
        'flex items-center justify-center',
        'bg-neutral-800 py-2 px-4 rounded-lg text-white font-semibold',
        'group group-aria aria-disabled:pointer-events-none aria-busy:pointer-events-none aria-disabled:opacity-50 aria-busy:opacity-50',
        className,
      )}
      disabled={disabled || pending}
      aria-disabled={disabled}
      aria-busy={pending}
      aria-label={label}
    >
      <span className="grid">
        <span
          className={classNames(
            'row-span-full col-span-full flex gap-2 items-center justify-center text-center group-aria-busy:opacity-0',
          )}
        >
          {label}
        </span>
        <span
          className="row-span-full col-span-full hidden group-aria-busy:flex items-center justify-center"
          aria-hidden
        >
          <Spinner />
        </span>
      </span>
    </button>
  )
}

export default React.memo(Button)

const Spinner = () => (
  <svg
    className="animate-spin h-5 w-5 text-white"
    xmlns="http://www.w3.org/2000/svg"
    fill="none"
    viewBox="0 0 24 24"
  >
    <circle
      className="opacity-25"
      cx="12"
      cy="12"
      r="10"
      stroke="currentColor"
      strokeWidth="4"
    ></circle>
    <path
      className="opacity-75"
      fill="currentColor"
      d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
    ></path>
  </svg>
)
