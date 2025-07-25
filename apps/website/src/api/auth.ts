import { API_URL } from './config'

export interface User {
  auth_date: number
  first_name: string
  hash: string
  id: number
  photo_url: string
  username: string
}

export async function isValidAuth(user: User) {
  const res = await fetch(`${API_URL}/auth`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(user),
  })

  return res.ok
}
