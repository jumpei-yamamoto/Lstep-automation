'use client'
import { useState } from 'react'
import UserForm from '@/components/users/UserForm'
import { createUser } from '@/services/users'
import { useRouter } from 'next/navigation'

export default function Page() {
  const [loading, setLoading] = useState(false)
  const router = useRouter()

  return (
    <div className="mx-auto max-w-2xl p-6">
      <h1 className="mb-6 text-2xl font-semibold">ユーザー登録</h1>
      <UserForm
        submitting={loading}
        onSubmit={async (v) => {
          setLoading(true)
          try {
            await createUser(v)
            router.push('/users')
          } finally {
            setLoading(false)
          }
        }}
      />
    </div>
  )
}