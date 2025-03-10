<script setup lang="ts">
import { ref, onMounted } from 'vue'

const props = defineProps<{
  width?: number
  height?: number
}>()

const emit = defineEmits<{
  (e: 'update:code', code: string): void
}>()

const captchaCanvas = ref<HTMLCanvasElement | null>(null)
const captchaCode = ref('')

// 生成随机颜色
const randomColor = (min: number, max: number) => {
  const r = Math.floor(Math.random() * (max - min + 1) + min)
  const g = Math.floor(Math.random() * (max - min + 1) + min)
  const b = Math.floor(Math.random() * (max - min + 1) + min)
  return `rgb(${r}, ${g}, ${b})`
}

// 生成随机码
const generateCode = () => {
  const characters = '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz'
  let code = ''
  for (let i = 0; i < 4; i++) {
    code += characters.charAt(Math.floor(Math.random() * characters.length))
  }
  return code
}

// 绘制验证码
const drawCaptcha = () => {
  const canvas = captchaCanvas.value
  if (!canvas) return

  const ctx = canvas.getContext('2d')
  if (!ctx) return

  // 生成新验证码
  captchaCode.value = generateCode()
  
  // 通知父组件当前验证码
  emit('update:code', captchaCode.value)

  // 清空画布
  ctx.clearRect(0, 0, canvas.width, canvas.height)

  // 绘制背景
  ctx.fillStyle = '#f5f5f5'
  ctx.fillRect(0, 0, canvas.width, canvas.height)

  // 绘制文字
  const fontSize = Math.floor(canvas.height * 0.6)
  ctx.font = `bold ${fontSize}px Arial`
  ctx.textBaseline = 'middle'
  
  // 随机偏移和旋转字符
  for (let i = 0; i < captchaCode.value.length; i++) {
    const x = (canvas.width / captchaCode.value.length) * i + 10
    const y = canvas.height / 2 + Math.random() * 8 - 4
    const rotate = (Math.random() - 0.5) * 0.5
    
    ctx.save()
    ctx.translate(x, y)
    ctx.rotate(rotate)
    ctx.fillStyle = randomColor(50, 120)
    ctx.fillText(captchaCode.value[i], 0, 0)
    ctx.restore()
  }

  // 绘制干扰线
  for (let i = 0; i < 4; i++) {
    ctx.strokeStyle = randomColor(150, 200)
    ctx.beginPath()
    ctx.moveTo(Math.random() * canvas.width, Math.random() * canvas.height)
    ctx.lineTo(Math.random() * canvas.width, Math.random() * canvas.height)
    ctx.stroke()
  }

  // 绘制干扰点
  for (let i = 0; i < 30; i++) {
    ctx.fillStyle = randomColor(150, 200)
    ctx.beginPath()
    ctx.arc(Math.random() * canvas.width, Math.random() * canvas.height, 1, 0, 2 * Math.PI)
    ctx.fill()
  }
}

// 刷新验证码
const refreshCaptcha = () => {
  drawCaptcha()
}

onMounted(() => {
  drawCaptcha()
})

defineExpose({
  refreshCaptcha,
  captchaCode
})
</script>

<template>
  <div class="captcha-container">
    <canvas
      ref="captchaCanvas"
      :width="props.width || 120"
      :height="props.height || 40"
      @click="refreshCaptcha"
      class="captcha-canvas"
    ></canvas>
    <div class="refresh-icon" @click="refreshCaptcha" title="刷新验证码">
      <i class="el-icon-refresh"></i>
    </div>
  </div>
</template>

<style scoped>
.captcha-container {
  position: relative;
  display: inline-block;
  cursor: pointer;
}

.captcha-canvas {
  border-radius: 4px;
  border: 1px solid #e4e7ed;
  vertical-align: middle;
}

.refresh-icon {
  position: absolute;
  right: 2px;
  bottom: 2px;
  font-size: 12px;
  color: #909399;
}

.refresh-icon:hover {
  color: #409EFF;
}
</style> 