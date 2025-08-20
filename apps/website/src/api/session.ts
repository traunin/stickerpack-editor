import { API_URL } from './config'

export interface User {
  auth_date: number
  first_name: string
  hash: string
  id: number
  photo_url: string
  username: string
}

export async function createSession(user: User) {
  const res = await fetch(`${API_URL}/session`, {
    method: 'POST',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(user),
  })

  return res.ok
}

export async function deleteSession() {
  const res = await fetch(`${API_URL}/session`, {
    method: 'DELETE',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
  })

  return res.ok
}
