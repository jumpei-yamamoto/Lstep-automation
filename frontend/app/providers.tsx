'use client'
import { RecoilRoot } from 'recoil'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { useState } from 'react'

export default function Providers({ children }: { children: React.ReactNode }) {
  const [qc] = useState(() => new QueryClient())
  return (
    <RecoilRoot>
      <QueryClientProvider client={qc}>{children}</QueryClientProvider>
    </RecoilRoot>
  )
}