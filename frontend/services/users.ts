import { api } from '@/lib/apiClient'

type CreateUserReq = { name: string; email: string }
type UserRes = { id: string; name: string; email: string }

export const createUser = (payload: CreateUserReq) =>
  api<UserRes>('/api/v1/users', {
    method: 'POST',
    body: JSON.stringify(payload),
  })