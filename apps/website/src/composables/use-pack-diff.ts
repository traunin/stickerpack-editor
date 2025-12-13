import { ref, watch } from 'vue'
import type { EditPackRequest } from '@/api/pack-edit'
import type { PackParameters } from '@/types/pack'
import type { Sticker } from '@/types/sticker'
import type { Ref } from 'vue'

export function usePackDiff(
  stickers: Ref<Sticker[]>,
  params: Ref<PackParameters>,
  originalPack: { title: string, isPublic: boolean, stickers: Sticker[] },
) {
  const edits: EditPackRequest = {
    deleted_stickers: [],
    added_stickers: [],
    emoji_updates: [],
    position_updates: [],
  }

  const hasEdits = ref(false)

  watch(
    [stickers, params],
    () => {
      // very inefficient, but alas
      const origIds = new Set(originalPack.stickers.map(s => s.id))
      const origMap = new Map(originalPack.stickers.map(s => [s.id, s]))
      const currentMap = new Map(stickers.value.map(s => [s.id, s]))
      // reset
      hasEdits.value = false
      edits.deleted_stickers = []
      edits.added_stickers = []
      edits.emoji_updates = []
      edits.position_updates = []
      // added
      for (const sticker of stickers.value) {
        if (!origIds.has(sticker.id)) {
          hasEdits.value = true
          edits.added_stickers.push(sticker)
        }
      }
      // deleted
      for (const sticker of originalPack.stickers) {
        if (!currentMap.has(sticker.id)) {
          hasEdits.value = true
          edits.deleted_stickers.push(sticker.id)
        }
      }
      // emoji
      for (const sticker of stickers.value) {
        const original = origMap.get(sticker.id)
        if (original && original.emoji_list.join('') !== sticker.emoji_list.join('')) {
          hasEdits.value = true
          edits.emoji_updates.push({
            id: sticker.id,
            emojis: [...sticker.emoji_list],
          })
        }
      }
      // drag
      for (let index = stickers.value.length - 1; index >= 0; index--) {
        const sticker = stickers.value[index]
        const origIndex = originalPack.stickers.findIndex(s => s.id === sticker.id)

        // only track changes in existing
        if (origIndex !== -1 && origIndex !== index) {
          hasEdits.value = true
          edits.position_updates.push({
            id: sticker.id,
            position: index,
          })
        }
      }
      // title
      if (originalPack.title !== params.value.title) {
        hasEdits.value = true
        edits.updated_title = params.value.title
      }
      // isPublic
      if (originalPack.isPublic !== params.value.isPublic) {
        hasEdits.value = true
        edits.updated_is_public = params.value.isPublic
      }
    },
    { deep: true },
  )
  return { edits, hasEdits }
}
