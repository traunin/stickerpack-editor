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
      queryClient.setQueriesData<{ pages: PacksResponse[], pageParams: number[] }>(
        { queryKey: ['packs', 'user'] },
        (old) => {
          if (!old) {
            return old
          }
          return {
            ...old,
            pages: old.pages.map(page => ({
              ...page,
              packs: page.packs.filter(pack => pack.name !== packName),
              total: page.total - 1,
            })),
          }
        },
      )

      queryClient.setQueriesData<{ pages: PacksResponse[], pageParams: number[] }>(
        { queryKey: ['packs', 'public'] },
        (old) => {
          if (!old) {
            return old
          }
          return {
            ...old,
            pages: old.pages.map(page => ({
              ...page,
              packs: page.packs.filter(pack => pack.name !== packName),
              total: page.total > 0 ? page.total - 1 : 0,
            })),
          }
        },
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
