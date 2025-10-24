<script setup>
import { NCard, NStatistic, NIcon } from 'naive-ui';
import { computed } from 'vue';

const props = defineProps({
  label: String,
  value: [String, Number],
  icon: Object,
  color: String,
});

const iconColor = computed(() => props.color || 'inherit');
</script>

<template>
  <n-card>
    <n-statistic :label="label">
      <template #prefix>
        <n-icon :component="icon" :color="iconColor" :size="24" />
      </template>
      <span :style="{ color: color }">{{ value }}</span>
    </n-statistic>
  </n-card>
</template>

<style scoped>
.n-statistic {
  text-align: center;
}
span {
  font-size: 28px;
  font-weight: 600;
}
</style>