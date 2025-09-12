import { type Emote } from '@/types/sticker'
import { computed, ref, type Ref, watch } from 'vue'

interface TenorSearchResponse {
  results: {
    id: string
    title: string
    media_formats: {
      tinygif?: { url: string }
      gif: { url: string }
      mediumgif?: { url: string }
      nanogif?: { url: string }
    }
  }[]
  next: string
}

const TENOR_API_KEY = import.meta.env.VITE_TENOR_API_KEY

async function fetchTenor(query: string, pos = '', limit = 10) {
  const url = new URL('https://tenor.googleapis.com/v2/search')
  url.searchParams.set('q', query)
  url.searchParams.set('key', TENOR_API_KEY)
  url.searchParams.set('limit', String(limit))
  if (pos) url.searchParams.set('pos', pos)

  const res = await fetch(url.toString())
  if (!res.ok) {
    throw new Error(`Failed to fetch from Tenor: ${res.status}`)
  }

  const json: TenorSearchResponse = await res.json()
  const items = json.results ?? []

  return {
    items: items.map((e): Emote => ({
      id: e.media_formats.gif.url.replace("https://media.tenor.com/", ""),
      name: e.title || e.id,
      preview: e.media_formats.tinygif?.url ?? e.media_formats.nanogif?.url ?? '',
      full: e.media_formats.gif?.url ?? e.media_formats.mediumgif?.url ?? '',
    })),
    next: json.next ?? '',
  }
}

export function useTenorSearch(query: Ref<string>, pageSize: number) {
  const emotes = ref<Emote[]>([])
  const error = ref<string | null>(null)

  const pos = ref('')
  const nextPos = ref('')
  const prevPos = ref<string[]>([])

  async function load(q: string, p: string) {
    error.value = null
    emotes.value = []
    try {
      const { items, next } = await fetchTenor(q, p, pageSize)
      emotes.value = items
      nextPos.value = next
    } catch (e: unknown) {
      error.value = String(e)
      nextPos.value = ''
    }
  }

  watch(query, (q) => {
    pos.value = ''
    nextPos.value = ''
    prevPos.value = []
    load(q, pos.value)
  }, { immediate: true })

  async function next() {
    if (nextPos.value) {
      prevPos.value.push(pos.value)
      pos.value = nextPos.value
      load(query.value, pos.value)
    }
  }

  async function prev() {
    if (prevPos.value.length > 0) {
      nextPos.value = pos.value
      pos.value = prevPos.value.pop()!
      load(query.value, pos.value)
    }
  }

  const hasNext = computed(() => nextPos.value !== '')
  const hasPrev = computed(() => prevPos.value.length > 0)

  return { emotes, error, hasNext, hasPrev, next, prev }
}