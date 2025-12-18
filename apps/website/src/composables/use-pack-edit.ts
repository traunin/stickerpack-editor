import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import type { ProgressEvent } from '@/api/job'
import type { EditPackRequest } from '@/api/pack-edit'
import { useCreatedPackStore } from '@/stores/use-created-pack'
import { useEditPackMutation } from './use-edit-pack-mutation'

export function usePackEdit() {
  const router = useRouter()
  const editedPack = useCreatedPackStore()
  const editPackMutation = useEditPackMutation()

  const progress = ref<ProgressEvent | null>(null)

  const isEditing = computed(() => editPackMutation.isPending.value)
  const editError = editPackMutation.editError

  async function editPack(name: string, editRequest: EditPackRequest) {
    if (name === '') {
      editError.show('Pack name is required')
      return
    }

    try {
      const response = await editPackMutation.mutateAsync({
        request: editRequest,
        name,
        onProgress: (progressEvent: ProgressEvent) => {
          progress.value = progressEvent
        },
      })

      editedPack.createdPack = response.pack

      await router.push({
        name: 'packCreated',
        params: { id: response.pack.name },
      })
    } catch (error) {
      console.error(error)
      throw error
    }
  }

  return {
    progress,
    isEditing,
    editError,
    editPack,
  }
}
