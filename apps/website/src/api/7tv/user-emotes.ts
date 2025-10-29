import type { Emote } from '@/types/sticker'

export interface User7TVEmotesResponse {
  data: {
    users: {
      user: {
        style: {
          activeEmoteSet: {
            emotes: {
              items: {
                alias: string
                emote: {
                  id: string
                  defaultName: string
                }
              }[]
              pageCount: number
              totalCount: number
            }
          }
        }
      }
    }
  }
}

const queryString = `
    query UserActiveEmotes($id: ID!, $page: Int!, $pageSize: Int!) {
      users {
        user(id: $id) {
          style {
            activeEmoteSet {
              emotes(page: $page, perPage: $pageSize) {
                totalCount
                pageCount
                items {
                  alias
                  emote {
                    id
                    defaultName
                  }
                }
              }
            }
          }
        }
      }
    }
  `

export async function fetch7TVUserEmotes(userID: string, page = 1, pageSize = 10) {
  const res = await fetch('https://api.7tv.app/v4/gql', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      operationName: 'UserActiveEmotes',
      query: queryString,
      variables: {
        id: userID,
        page,
        pageSize,
      },
    }),
  })

  if (!res.ok) {
    throw new Error(`Failed to fetch emotes: ${res.status}`)
  }

  const json: User7TVEmotesResponse = await res.json()
  const emotesData = json.data.users.user.style.activeEmoteSet.emotes

  return {
    items: emotesData.items.map(
      (e): Emote => ({
        id: e.emote.id,
        name: e.emote.defaultName,
        preview: `https://cdn.7tv.app/emote/${e.emote.id}/2x.webp`,
        full: `https://cdn.7tv.app/emote/${e.emote.id}/4x.webp`,
      }),
    ),
    pageCount: emotesData.pageCount,
    totalCount: emotesData.totalCount,
  }
}
