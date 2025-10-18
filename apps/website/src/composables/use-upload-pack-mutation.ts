import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { ref } from 'vue'
import { uploadPack } from '@/api/stickerpack-upload'
import type { ProgressEvent, StickerpackRequest } from '@/api/stickerpack-upload'

export function useUploadPackMutation() {
  const queryClient = useQueryClient()
  const uploadError = ref<string>('')

  const mutation = useMutation({
    mutationFn: ({ request, onProgress }: {
      request: StickerpackRequest
      onProgress?: (progress: ProgressEvent) => void
    }) => uploadPack(request, onProgress),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['packs', 'user'] })
      queryClient.invalidateQueries({ queryKey: ['packs', 'public'] })
    },
    onError: (error) => {
      uploadError.value = error.message
      setTimeout(() => {
        uploadError.value = ''
      }, 4000)
    },
  })

  return {
    ...mutation,
    uploadError,
  }
}
