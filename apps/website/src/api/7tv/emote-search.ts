import type { Emote } from '@/types/sticker'

interface EmoteSearchResponse {
  data: {
    search: {
      all: {
        emotes: {
          items: {
            id: string
            defaultName: string
          }[]
          pageCount: number
          totalCount: number
        }
      }
    }
  }
}

const queryString = `
query GlobalSearch($query: String!, $page: Int!, $pageSize: Int!) {
  search {
    all(query: $query, page: $page, perPage: $pageSize) {
      emotes {
        items {
          id
          defaultName
        }
        pageCount
        totalCount
      }
    }
  }
}
`

export async function fetchSearchedEmotes(query: string, page = 1, pageSize = 10) {
  const res = await fetch('https://api.7tv.app/v4/gql', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      operationName: 'GlobalSearch',
      query: queryString,
      variables: {
        query,
        isDefaultSetSet: false,
        defaultSetId: '',
        page,
        pageSize,
      },
    }),
  })

  if (!res.ok) {
    throw new Error(`Failed to fetch emotes: ${res.status}`)
  }

  const json: EmoteSearchResponse = await res.json()
  const emotes = json.data.search.all.emotes

  return {
    items: emotes.items.map(
      (e): Emote => ({
        id: e.id,
        name: e.defaultName,
        preview: `https://cdn.7tv.app/emote/${e.id}/2x.webp`,
        full: `https://cdn.7tv.app/emote/${e.id}/4x.webp`,
      }),
    ),
    pageCount: emotes.pageCount,
  }
}
