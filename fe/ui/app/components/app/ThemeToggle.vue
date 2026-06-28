<script setup lang="ts">
const colorMode = useColorMode()

const isDark = computed(() => {
  return colorMode.value === 'dark'
})

const toggleTheme = (event: MouseEvent) => {
  // Check if View Transitions API is supported
  const isAppearanceTransition = typeof document !== 'undefined' 
    && typeof document.startViewTransition === 'function' 
    && !window.matchMedia('(prefers-reduced-motion: reduce)').matches

  const nextMode = isDark.value ? 'light' : 'dark'

  if (!isAppearanceTransition) {
    colorMode.preference = nextMode
    return
  }

  const x = event.clientX
  const y = event.clientY
  const endRadius = Math.hypot(
    Math.max(x, innerWidth - x),
    Math.max(y, innerHeight - y)
  )

  const transition = document.startViewTransition(() => {
    colorMode.preference = nextMode
  })

  transition.ready.then(() => {
    const clipPath = [
      `circle(0px at ${x}px ${y}px)`,
      `circle(${endRadius}px at ${x}px ${y}px)`
    ]

    document.documentElement.animate(
      {
        clipPath
      },
      {
        duration: 750,
        easing: 'cubic-bezier(0.4, 0, 0.2, 1)',
        pseudoElement: '::view-transition-new(root)'
      }
    )
  })
}
</script>

<template>
  <div class="flex items-center">
    <!-- Light mode button (hidden in dark mode) -->
    <UButton
      class="dark:hidden"
      icon="i-lucide-sun"
      color="neutral"
      variant="ghost"
      aria-label="Theme"
      @click="toggleTheme"
    />
    <!-- Dark mode button (hidden in light mode) -->
    <UButton
      class="hidden dark:inline-flex"
      icon="i-lucide-moon"
      color="neutral"
      variant="ghost"
      aria-label="Theme"
      @click="toggleTheme"
    />
  </div>
</template>
