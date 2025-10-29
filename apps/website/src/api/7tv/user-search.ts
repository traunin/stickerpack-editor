import type { User7TV } from '@/types/user-7tv'

interface UserSearchResponse {
  data: {
    search: {
      all: {
        users: {
          items: {
            id: string
            mainConnection: {
              platformDisplayName: string
              platformAvatarUrl: string
            }
          }[]
        }
      }
    }
  }
}

const queryString = `
query GlobalSearch($query: String!, $perPage: Int!) {
  search {
    all(query: $query, page: 1, perPage: $perPage) {
        users {
          items{
            id
            mainConnection {
                platformDisplayName
                platformAvatarUrl
            }
          }
        }
      }
    }
  }
`

export async function fetchSearchedUsers(query: string, results = 10) {
  const res = await fetch('https://api.7tv.app/v4/gql', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      operationName: 'GlobalSearch',
      query: queryString,
      variables: {
        query,
        perPage: results,
      },
    }),
  })

  if (!res.ok) {
    throw new Error(`Failed to fetch users: ${res.status}`)
  }

  const json: UserSearchResponse = await res.json()
  const users = json.data.search.all.users.items

  return users.map(
    (e): User7TV => ({
      id: e.id,
      name: e.mainConnection.platformDisplayName,
      avatarURL: e.mainConnection.platformAvatarUrl,
    }),
  )
}
