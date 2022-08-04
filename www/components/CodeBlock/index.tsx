import React, { useCallback, useState } from 'react'
import Highlight, { defaultProps, Language } from 'prism-react-renderer'
import builtinTheme from 'prism-react-renderer/themes/ultramin'
import classNames from 'classnames'
import { copyToClipboard } from 'utils/clipboard'
import loadSolidityLanguageHighlighting from './loadSolidityLanguageHighlighting'

loadSolidityLanguageHighlighting()

const theme = {
  ...builtinTheme,
  plain: {
    backgroundColor: 'transparent',
  },
}

type Props = {
  code: string
  codeForCopy?(): string
  language: Language | 'sol' | 'txt'
  title?: string
}

const CodeBlock = ({ title, code, codeForCopy, language }: Props) => {
  const [justCopied, justCopiedSet] = useState(false)

  const copyCode = useCallback(() => {
    copyToClipboard(codeForCopy?.() ?? code)
    justCopiedSet(true)
    const id = setTimeout(() => justCopiedSet(false), 1000)
    return () => clearTimeout(id)
  }, [code, codeForCopy])

  return (
    <div className="group flex flex-col border-2 rounded-lg overflow-hidden divide-y-2">
      {title !== undefined && (
        <div
          className={classNames(
            'flex justify-between items-center',
            'bg-neutral-50 text-neutral-500 text-lg px-4 py-2',
          )}
        >
          {title}
          <button
            className="transition-colors hover:text-neutral-800"
            onClick={copyCode}
          >
            {justCopied ? 'Copied!' : 'Copy code'}
          </button>
        </div>
      )}
      <Highlight
        {...defaultProps}
        theme={theme}
        code={code}
        language={language as Language}
      >
        {({ className, style, tokens, getLineProps, getTokenProps }) => (
          <pre
            className={classNames('p-4 overflow-x-scroll', className)}
            style={style}
          >
            {tokens.map((line, i) => (
              <div key={i} {...getLineProps({ line, key: i })}>
                {line.map((token, key) => (
                  <span key={key} {...getTokenProps({ token, key })} />
                ))}
              </div>
            ))}
          </pre>
        )}
      </Highlight>
    </div>
  )
}

export default React.memo(CodeBlock)
