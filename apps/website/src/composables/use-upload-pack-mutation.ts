import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { uploadPack } from '@/api/stickerpack-upload'
import type { ProgressEvent, StickerpackRequest } from '@/api/stickerpack-upload'
import { useErrorPopup } from './use-error-popup'

export function useUploadPackMutation() {
  const queryClient = useQueryClient()
  const uploadError = useErrorPopup()

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
      uploadError.show(error)
    },
  })

  return {
    ...mutation,
    uploadError,
  }
}
