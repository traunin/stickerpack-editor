import { computed, unref } from 'vue'
import { useRouter } from 'vue-router'
import { useDeletePackMutation } from '@/composables/use-delete-pack-mutation'
import type { MaybeRef } from 'vue'

export function usePackDeletion() {
  const router = useRouter()
  const deletePackMutation = useDeletePackMutation()

  const isDeleting = computed(() => deletePackMutation.isPending.value)
  const deletionError = deletePackMutation.deletionError

  async function deletePack(packName: MaybeRef<string | null>) {
    const name = unref(packName)
    if (name === '' || name === null) {
      deletionError.show('Pack name is required')
      return
    }

    try {
      await deletePackMutation.mutateAsync(name)
      // go to main page
      await router.push('/')
    } catch (error) {
      console.error(error)
      throw error
    }
  }

  return {
    isDeleting,
    deletionError,
    deletePack,
  }
}
