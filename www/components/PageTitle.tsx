import React from 'react'

export default function PageTitle({ children }: { children: React.ReactNode }) {
  return <div className="font-bold text-2xl md:text-3xl mb-10">{children}</div>
}
