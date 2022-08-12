import React, { ReactNode } from 'react'

type Props = {
  title: ReactNode
  description?: ReactNode
  children: ReactNode
}

function Section({ title, description, children }: Props) {
  return (
    <div className="flex flex-col gap-y-4">
      <h1 className="font-bold text-xl sm:text-2xl">{title}</h1>
      {description !== undefined && <p>{description}</p>}
      {children}
    </div>
  )
}

export default React.memo(Section)
