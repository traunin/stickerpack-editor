import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { ref } from 'vue'
import { deletePack } from '@/api/pack-delete'
import type { PacksResponse } from '@/api/packs'

export function useDeletePackMutation() {
  const queryClient = useQueryClient()
  const deletionError = ref<string>('')

  const mutation = useMutation({
    mutationFn: deletePack,
    onSuccess: (_, packName) => {
      queryClient.setQueriesData<PacksResponse>(
        { queryKey: ['packs', 'user'] },
        (old) => old ?
            {
              ...old,
              packs: old.packs.filter(pack => pack.name !== packName),
              total: old.total - 1,
            } :
          old,
      )

      queryClient.setQueriesData<PacksResponse>(
        { queryKey: ['packs', 'public'] },
        (old) => old ?
            {
              ...old,
              packs: old.packs.filter(pack => pack.name !== packName),
              total: old.total > 0 ? old.total - 1 : 0,
            } :
          old,
      )

      deletionError.value = ''
    },
    onError: (error) => {
      deletionError.value = error.message
      setTimeout(() => {
        deletionError.value = ''
      }, 4000)
    },
  })

  return {
    ...mutation,
    deletionError,
  }
}
