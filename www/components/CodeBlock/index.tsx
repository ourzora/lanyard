import React, { useCallback } from 'react'
import Highlight, { defaultProps, Language } from 'prism-react-renderer'
import builtinTheme from 'prism-react-renderer/themes/ultramin'
import classNames from 'classnames'
import loadSolidityLanguageHighlighting from './loadSolidityLanguageHighlighting'
import CopyCodeButton from './CopyCodeButton'

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
  oneLiner?: boolean
  title?: string
}

const CodeBlock = ({
  title,
  code,
  codeForCopy,
  oneLiner = false,
  language,
}: Props) => {
  const codeForCopyButton = useCallback(
    () => codeForCopy?.() ?? code,
    [codeForCopy, code],
  )

  return (
    <div className="group flex flex-col border-2 rounded-lg overflow-hidden divide-y-2">
      {title !== undefined && (
        <div
          className={classNames(
            'flex justify-between items-center',
            'text-neutral-500 text-lg px-4 py-2 h-14',
          )}
        >
          <span className="font-mono text-base">{title}</span>
          <CopyCodeButton codeForCopy={codeForCopyButton} />
        </div>
      )}

      <div className="flex justify-between items-center min-h-[5rem] p-4 gap-x-4">
        {oneLiner ? (
          <pre className="text-ellipsis overflow-hidden">{code}</pre>
        ) : (
          <Highlight
            {...defaultProps}
            theme={theme}
            code={code}
            language={language as Language}
          >
            {({ className, style, tokens, getLineProps, getTokenProps }) => (
              <pre
                className={classNames('overflow-x-scroll', className)}
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
        )}
        {oneLiner && title === undefined && (
          <CopyCodeButton
            codeForCopy={codeForCopyButton}
            className="flex-shrink-0"
          />
        )}
      </div>
    </div>
  )
}

export default React.memo(CodeBlock)
