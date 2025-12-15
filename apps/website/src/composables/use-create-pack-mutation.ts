import { useMutation, useQueryClient } from '@tanstack/vue-query'
import type { ProgressEvent } from '@/api/job'
import { createPack } from '@/api/pack-create'
import type { CreatePackRequest } from '@/api/pack-create'
import { useErrorPopup } from '@/composables/use-error-popup'

export function useCreatePackMutation() {
  const queryClient = useQueryClient()
  const uploadError = useErrorPopup()

  const mutation = useMutation({
    mutationFn: ({ request, onProgress }: {
      request: CreatePackRequest
      onProgress?: (progress: ProgressEvent) => void
    }) => createPack(request, onProgress),
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
