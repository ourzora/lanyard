import classNames from 'classnames'
import { ReactNode } from 'react'

const Container = ({
  children,
  variant,
}: {
  children: ReactNode
  variant: 'success' | 'failure'
}) => (
  <div
    className={classNames(
      'border-2 rounded-lg px-4 py-2 overflow-auto',
      { 'bg-brand-light border-brand': variant === 'success' },
      { 'bg-error-light border-error': variant === 'failure' },
    )}
  >
    {children}
  </div>
)

export const FailureContainer = ({ children }: { children: ReactNode }) => (
  <Container variant="failure">{children}</Container>
)

export const SuccessContainer = ({ children }: { children: ReactNode }) => (
  <Container variant="success">{children}</Container>
)
