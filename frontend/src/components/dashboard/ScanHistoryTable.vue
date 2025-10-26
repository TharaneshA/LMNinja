<script setup>
import { NDataTable } from 'naive-ui';
import { ref } from 'vue';

const props = defineProps({
  scanData: {
    type: Array,
    required: true,
  },
});


const columns = [
  { title: 'Date & Time', key: 'date', width: 180 },
  { title: 'Target Model', key: 'targetModel' },
  { title: 'Attack Categories', key: 'categories' },
  { title: 'Prompts Run', key: 'prompts', width: 120 },
  { title: 'Vulnerabilities', key: 'vulns', width: 140 },
  { title: 'Pass Rate', key: 'passRate', width: 100 },
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