
import { API_URL } from './config'

interface DeletePackResponse {
  success: boolean
}

export async function deletePack(name: string) {
  try {
    const res = await fetch(`${API_URL}/user/packs/${name}`, {
      method: 'DELETE',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
    })

    if (!res.ok) {
      const text = await res.text()
      let errorMessage: string

      try {
        const errJson = JSON.parse(text)
        errorMessage = errJson.message || JSON.stringify(errJson)
      } catch {
        errorMessage = text
      }

      throw new Error(`Failed to delete pack: ${errorMessage}`)
    }

    const data: DeletePackResponse = await res.json()
    return data
  } catch (err: any) {
    console.error('deletePack error:', err)
    throw new Error(err?.message || 'Unknown error while deleting pack')
  }
}
