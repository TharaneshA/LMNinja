<script setup>
import { NDataTable } from 'naive-ui';
import { ref } from 'vue';

// Mock data
const createData = () => [
  {
    date: '2025-09-24 10:30:15',
    targetModel: 'Llama-3-8B',
    categories: 'Jailbreaking, PII Leaks',
    prompts: '20/60',
    vulns: 5,
    passRate: '75.0%',
  },
  {
    date: '2025-09-23 18:15:45',
    targetModel: 'Gemini-1.5-Flash',
    categories: 'Prompt Injection',
    prompts: '10/30',
    vulns: 8,
    passRate: '20.0%',
  },
  {
    date: '2025-09-23 11:05:20',
    targetModel: 'Qwen3-0.6B',
    categories: 'Jailbreaking',
    prompts: '10/30',
    vulns: 1,
    passRate: '90.0%',
  },
];

const columns = [
  { title: 'Date & Time', key: 'date', width: 180 },
  { title: 'Target Model', key: 'targetModel' },
  { title: 'Attack Categories', key: 'categories' },
  { title: 'Prompts Run', key: 'prompts', width: 120 },
  { title: 'Vulnerabilities', key: 'vulns', width: 140 },
  { title: 'Pass Rate', key: 'passRate', width: 100 },
  // Actions button can be added here later
];

const data = ref(createData());
</script>

<template>
  <n-data-table
    :columns="columns"
    :data="data"
    :pagination="{ pageSize: 10 }"
    :single-line="false"
  />
</template>