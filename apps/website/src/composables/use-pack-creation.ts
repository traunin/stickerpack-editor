import { computed, ref, toRaw } from 'vue'
import { useRouter } from 'vue-router'
import type { ProgressEvent } from '@/api/job'
import { useCreatePackMutation } from '@/composables/use-create-pack-mutation'
import { useCreatedPackStore } from '@/stores/use-created-pack'
import type { PackParameters } from '@/types/pack'
import type { Sticker } from '@/types/sticker'

export function usePackCreation() {
  const router = useRouter()
  const createdPack = useCreatedPackStore()
  const createPackMutation = useCreatePackMutation()

  const progress = ref<ProgressEvent | null>(null)

  const isCreating = computed(() => createPackMutation.isPending.value)
  const createError = createPackMutation.uploadError

  async function createPack(packParams: PackParameters, stickers: Sticker[]) {
    if (packParams.name === '' || packParams.name === null) {
      createError.show('Pack name is required')
      return
    }
    if (packParams.hasWatermark === undefined) {
      createError.show('Watermark info is required')
      return
    }
    try {
      const response = await createPackMutation.mutateAsync({
        request: {
          pack_name: packParams.name ?? '',
          title: packParams.title,
          emotes: stickers.map(e => toRaw(e)),
          has_watermark: packParams.hasWatermark,
          is_public: packParams.isPublic,
        },
        onProgress: (progressEvent: ProgressEvent) => {
          progress.value = progressEvent
        },
      })

      createdPack.createdPack = response.pack

      await router.push({
        name: 'packCreated',
        params: { id: response.pack.id },
      })
    } catch (error) {
      console.error(error)
      throw error
    }
  }

  return {
    progress,
    isCreating,
    createError,
    createPack,
  }
}
