<script setup lang="ts">
defineProps<{
  steps: string[]
  stepKeys: string[]
  current: string
}>()
</script>

<template>
  <div class="steps">
    <div
      v-for="(label, i) in steps"
      :key="i"
      class="step"
      :class="{
        'step--active': stepKeys[i] === current,
        'step--done':   stepKeys.indexOf(current) > i,
      }"
    >
      <div class="step__dot">
        <svg v-if="stepKeys.indexOf(current) > i" width="12" height="12" viewBox="0 0 12 12" fill="none">
          <path d="M2 6l3 3 5-5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
        </svg>
        <span v-else>{{ i + 1 }}</span>
      </div>
      <span class="step__label">{{ label }}</span>
    </div>
  </div>
</template>

<style scoped>
.steps {
  display: flex;
  align-items: flex-start;
  gap: 0;
  padding: 20px 0;
  border-bottom: 1px solid var(--border);
  margin-bottom: 32px;
  overflow-x: auto;
}
.step {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  position: relative;
}
.step:not(:last-child)::after {
  content: '';
  position: absolute;
  top: 16px;
  left: 50%;
  width: 100%;
  height: 1px;
  background: var(--border);
  z-index: 0;
}
.step--done:not(:last-child)::after { background: var(--accent); }

.step__dot {
  width: 32px; height: 32px;
  border-radius: 50%;
  border: 2px solid var(--border-strong);
  background: var(--bg-card);
  display: flex; align-items: center; justify-content: center;
  font-size: 13px; font-weight: 600; color: var(--text-3);
  transition: all 0.2s;
  position: relative; z-index: 1;
}
.step--active .step__dot {
  border-color: var(--accent);
  background: var(--accent-bg);
  color: var(--accent);
}
.step--done .step__dot {
  border-color: var(--accent);
  background: var(--accent);
  color: #fff;
}

.step__label {
  font-size: 12px; font-weight: 500;
  color: var(--text-3); white-space: nowrap;
  transition: color 0.2s;
}
.step--active .step__label { color: var(--accent); }
.step--done   .step__label { color: var(--text-2); }
</style>