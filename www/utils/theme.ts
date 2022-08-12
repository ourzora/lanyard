import classNames from 'classnames'

export const background = classNames('bg-white')
export const textPrimary = classNames('text-black')
export const textSecondary = classNames('text-grey-800')

export const bodyStyles = classNames(textPrimary, 'text-base', background)

export const brandUnderlineClasses = classNames(
  'transition-all border-b-2 hover:border-b-4 border-brand',
)
