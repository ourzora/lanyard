import React from 'react'

type Props = {
  title: string
  description?: string
  children: React.ReactNode
}

function Section({ title, description, children }: Props) {
  return (
    <div className="flex flex-col gap-4">
      <h1 className="font-bold text-xl sm:text-2xl">{title}</h1>
      {description !== undefined && <p>{description}</p>}
      {children}
    </div>
  )
}

export default React.memo(Section)
