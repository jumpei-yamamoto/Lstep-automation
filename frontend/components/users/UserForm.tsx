'use client'
import { useState } from 'react'

type Props = { 
  onSubmit: (v: { name: string; email: string }) => Promise<void>; 
  submitting?: boolean 
}

export default function UserForm({ onSubmit, submitting }: Props) {
  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [error, setError] = useState<string | null>(null)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError(null)
    try {
      await onSubmit({ name, email })
    } catch (err: any) {
      setError(err.message ?? 'エラーが発生しました')
    }
  }

  return (
    <form onSubmit={handleSubmit} className="max-w-md space-y-4 rounded-lg bg-white p-6 shadow">
      <div>
        <label className="block text-sm font-medium">名前</label>
        <input 
          className="mt-1 w-full rounded border px-3 py-2" 
          value={name} 
          onChange={e => setName(e.target.value)} 
        />
      </div>
      <div>
        <label className="block text-sm font-medium">メール</label>
        <input 
          className="mt-1 w-full rounded border px-3 py-2" 
          value={email} 
          onChange={e => setEmail(e.target.value)} 
        />
      </div>
      {error && <p className="text-sm text-red-600">{error}</p>}
      <button
        disabled={submitting}
        className="inline-flex items-center rounded bg-black px-4 py-2 text-white disabled:opacity-50"
      >
        {submitting ? '送信中...' : '登録'}
      </button>
    </form>
  )
}