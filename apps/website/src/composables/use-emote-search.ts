import { ref, type Ref, watch } from 'vue'

export interface Emote {
  id: string
  name: string
  url: string
}

interface EmoteSearchResponse {
  data: {
    search: {
      all: {
        emotes: {
          items: {
            id: string
            defaultName: string
            images: {
              url: string
            }[]
          }[]
          pageCount: number
          totalCount: number
        }
      }
    }
  }
}

async function fetchEmotes(query: string, page = 1, pageSize = 10) {
  const res = await fetch('https://api.7tv.app/v4/gql', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      operationName: 'GlobalSearch',
      query: `
        query GlobalSearch($query: String!) {
          search {
            all(query: $query, page: ${page}, perPage: ${pageSize}) {
              emotes {
                items {
                  id
                  defaultName
                  images {
                    url
                  }
                }
                pageCount
                totalCount
              }
            }
          }
        }
      `,
      variables: {
        query,
        isDefaultSetSet: false,
        defaultSetId: '',
      },
    }),
  })

  if (!res.ok) {
    throw new Error(`Failed to fetch emotes: ${res.status}`)
  }

  const json: EmoteSearchResponse = await res.json()
  const emotes = json.data.search.all.emotes
  const items = emotes.items

  return { items: items.map((e): Emote => ({
    id: e.id,
    name: e.defaultName,
    url: `https://cdn.7tv.app/emote/${e.id}/2x.webp`,
    // url: `${e.images[0]?.url ?? ''}`,
  })), pageCount: emotes.totalCount }
}

export function useEmoteSearch(query: Ref<string>, pageSize: number) {
  const page = ref(1)
  const maxPages = ref(1)
  const emotes = ref<Emote[] | null>()
  const error = ref<string | null>()

  watch(query, async (newQuery) => {
    page.value = 1
    error.value = null
    emotes.value = []

    try {
      const { items, pageCount } = await fetchEmotes(newQuery, 1, pageSize)
      emotes.value = items
      maxPages.value = pageCount
    } catch (e: unknown) {
      error.value = String(e)
      maxPages.value = 1
    }
  }, { immediate: true })

  watch([query, page], async ([q, p]) => {
    error.value = null
    emotes.value = []

    try {
      const { items } = await fetchEmotes(q, p, pageSize)
      emotes.value = items
    } catch (e: unknown) {
      error.value = String(e)
    }
  })

  function next() {
    if (page.value < maxPages.value) {
      page.value++
    }
  }

  function prev() {
    if (page.value > 1) {
      page.value--
    }
  }

  return { emotes, error, page, maxPages, next, prev }
}
