import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { deletePack } from '@/api/pack-delete'
import type { PacksResponse } from '@/api/packs'
import { useErrorPopup } from './use-error-popup'

export function useDeletePackMutation() {
  const queryClient = useQueryClient()
  const deletionError = useErrorPopup()

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

      deletionError.clear()
    },
    onError: (error) => {
      deletionError.show(error)
    },
  })

  return {
    ...mutation,
    deletionError,
  }
}
