import React, { useState, useCallback } from 'react'
import classNames from 'classnames'
import { copyToClipboard } from 'utils/clipboard'

type Props = {
  codeForCopy(): string
  className?: string
}

function CopyCodeButton({ codeForCopy, className }: Props) {
  const [justCopied, justCopiedSet] = useState(false)

  const copyCode = useCallback(() => {
    copyToClipboard(codeForCopy())
    justCopiedSet(true)
    const id = setTimeout(() => justCopiedSet(false), 1000)
    return () => clearTimeout(id)
  }, [codeForCopy])

  return (
    <button
      className={classNames(
        'bg-neutral-100 hover:bg-neutral-200 border-2 border-neutral-200 rounded-md',
        'transition-colors',
        'px-2 py-0.5 w-[6.6rem] h-9',
        'font-bold text-black text-sm',
        className,
      )}
      onClick={copyCode}
    >
      {justCopied ? 'Copied!' : 'Copy code'}
    </button>
  )
}

export default React.memo(CopyCodeButton)
