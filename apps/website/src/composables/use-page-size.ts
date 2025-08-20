import { ref, onMounted, onUnmounted, Ref, nextTick, watch } from "vue";

export function usePageSize(containerRef: Ref<HTMLElement | null>) {
  const pageSize = ref(1);
  
  const updatePageSize = async () => {
    await nextTick();
    
    const container = containerRef.value as HTMLElement | null;
    if (!container) {
      return;
    }
    
    const firstChild = container.querySelector<HTMLElement>(":first-child");
    if (!firstChild) {
      return;
    }
    
    const containerWidth = container.clientWidth;
    const packWidth = firstChild.offsetWidth;
    
    const newPageSize = Math.max(1, Math.floor(containerWidth / packWidth));
    if (newPageSize !== pageSize.value) {
      pageSize.value = newPageSize;
    }
  };

  watch(containerRef, (newContainer) => {
    if (newContainer) {
      updatePageSize();
    }
  }, { immediate: true });

  onMounted(() => {
    setTimeout(updatePageSize, 0);
    setTimeout(updatePageSize, 100);
    
    window.addEventListener("resize", updatePageSize);
  });

  onUnmounted(() => {
    window.removeEventListener("resize", updatePageSize);
  });

  return {
    pageSize,
    updatePageSize, // trigger when the api returns packs
  };
}