import { useMutation, useQueryClient } from '@tanstack/vue-query'
import type { ProgressEvent } from '@/api/job'
import { editPack } from '@/api/pack-edit'
import type { EditPackRequest } from '@/api/pack-edit'
import { useErrorPopup } from '@/composables/use-error-popup'

export function useEditPackMutation() {
  const queryClient = useQueryClient()
  const editError = useErrorPopup()

  const mutation = useMutation({
    mutationFn: ({ name, request, onProgress }: {
      name: string
      request: EditPackRequest
      onProgress?: (progress: ProgressEvent) => void
    }) => editPack(name, request, onProgress),
    onSuccess: (_data, variables) => {
      // invalidate edited pack
      queryClient.invalidateQueries({ queryKey: ['pack', variables.name] })
      // preview might change, otherwise... Might be overkill
      queryClient.invalidateQueries({ queryKey: ['packs', 'user'] })
      queryClient.invalidateQueries({ queryKey: ['packs', 'public'] })
    },
    onError: (error) => {
      editError.show(error)
    },
  })

  return {
    ...mutation,
    editError,
  }
}
