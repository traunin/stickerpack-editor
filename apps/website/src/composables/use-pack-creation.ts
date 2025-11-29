import { computed, ref, toRaw } from 'vue'
import { useRouter } from 'vue-router'
import { useUploadPackMutation } from '@/composables/use-upload-pack-mutation'
import { useCreatedPackStore } from '@/stores/use-created-pack'
import type { PackParameters } from '@/types/pack'
import type { Sticker } from '@/types/sticker'

export function usePackCreation() {
  const router = useRouter()
  const createdPack = useCreatedPackStore()
  const uploadPackMutation = useUploadPackMutation()

  const progress = ref<ProgressEvent | null>(null)

  const isUploading = computed(() => uploadPackMutation.isPending.value)
  const uploadError = uploadPackMutation.uploadError

  async function createPack(packParams: PackParameters, stickers: Sticker[]) {
    if (packParams.name === '' || packParams.name === null) {
      uploadError.show('Pack name is required')
      return
    }
    if (packParams.hasWatermark === undefined) {
      uploadError.show('Watermark info is required')
      return
    }
    try {
      const response = await uploadPackMutation.mutateAsync({
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
    isUploading,
    uploadError,
    createPack,
  }
}
