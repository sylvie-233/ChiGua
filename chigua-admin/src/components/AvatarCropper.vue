<script setup lang="ts">
import { ref, watch } from 'vue'

const props = defineProps<{ imageUrl: string }>()
const emit = defineEmits<{ cropped: [blob: Blob] }>()

const canvasRef = ref<HTMLCanvasElement>()
const cropSize = 200
const scale = ref(1)
const offsetX = ref(0)
const offsetY = ref(0)
const img = new Image()
let dragging = false
let dragStart = { x: 0, y: 0 }

img.onload = () => {
  scale.value = Math.max(cropSize / img.width, cropSize / img.height)
  offsetX.value = 0
  offsetY.value = 0
  draw()
}

watch(() => props.imageUrl, (url) => {
  if (url) img.src = url
}, { immediate: true })

const draw = () => {
  const canvas = canvasRef.value
  if (!canvas || !img.complete) return
  const ctx = canvas.getContext('2d')!
  const w = img.width * scale.value
  const h = img.height * scale.value
  const cx = (canvas.width - w) / 2 + offsetX.value
  const cy = (canvas.height - h) / 2 + offsetY.value

  ctx.clearRect(0, 0, canvas.width, canvas.height)
  // 暗色遮罩
  ctx.fillStyle = 'rgba(0,0,0,0.5)'
  ctx.fillRect(0, 0, canvas.width, canvas.height)
  // 裁剪区域
  ctx.save()
  ctx.beginPath()
  ctx.arc(cropSize, cropSize, cropSize, 0, Math.PI * 2)
  ctx.clip()
  ctx.drawImage(img, cx, cy, w, h)
  ctx.restore()
}

const getCroppedBlob = (): Promise<Blob> => {
  return new Promise(resolve => {
    const out = document.createElement('canvas')
    out.width = cropSize * 2
    out.height = cropSize * 2
    const ctx = out.getContext('2d')!
    ctx.beginPath()
    ctx.arc(cropSize, cropSize, cropSize, 0, Math.PI * 2)
    ctx.clip()
    const w = img.width * scale.value
    const h = img.height * scale.value
    const cx = (400 - w) / 2 + offsetX.value
    const cy = (400 - h) / 2 + offsetY.value
    ctx.drawImage(img, cx, cy, w, h)
    out.toBlob(b => resolve(b!), 'image/jpeg', 0.9)
  })
}

const onWheel = (e: WheelEvent) => {
  e.preventDefault()
  scale.value = Math.max(0.5, Math.min(3, scale.value - e.deltaY * 0.001))
  draw()
}

const onMouseDown = (e: MouseEvent) => {
  dragging = true
  dragStart = { x: e.clientX - offsetX.value, y: e.clientY - offsetY.value }
}
const onMouseMove = (e: MouseEvent) => {
  if (!dragging) return
  offsetX.value = e.clientX - dragStart.x
  offsetY.value = e.clientY - dragStart.y
  draw()
}
const onMouseUp = () => { dragging = false }

const handleCrop = async () => {
  const blob = await getCroppedBlob()
  emit('cropped', blob)
}

defineExpose({ handleCrop })
</script>

<template>
  <div class="cropper-wrapper">
    <canvas
      ref="canvasRef"
      :width="cropSize * 2"
      :height="cropSize * 2"
      @wheel.prevent="onWheel"
      @mousedown="onMouseDown"
      @mousemove="onMouseMove"
      @mouseup="onMouseUp"
      @mouseleave="onMouseUp"
      style="cursor: grab; max-width: 100%; border-radius: 4px;"
    />
    <p style="text-align: center; color: #999; font-size: 12px; margin-top: 8px;">
      滚轮缩放 · 拖拽移动
    </p>
  </div>
</template>
