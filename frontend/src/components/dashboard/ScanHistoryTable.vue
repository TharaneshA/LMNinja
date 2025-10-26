<script setup>
import { NDataTable, NButton, NIcon } from 'naive-ui';
import { h } from 'vue';
import Detail from '@/components/icons/Detail.vue';

const props = defineProps({
  scanData: { type: Array, required: true },
});
const emit = defineEmits(['view-report']);

const createColumns = ({ onViewReport }) => [
  { title: 'Date & Time', key: 'startTime', width: 180, sorter: 'default' },
  { title: 'Target Model', key: 'targetModelName', sorter: 'default' },
  { title: 'Status', key: 'status', width: 120 },
  { title: 'Prompts', key: 'totalPrompts', width: 100, sorter: 'default' },
  { title: 'Vulnerabilities', key: 'vulnerabilitiesFound', width: 140, sorter: 'default' },
  { title: 'Pass Rate', key: 'passRate', width: 110 },
  {
    title: 'Actions',
    key: 'actions',
    render(row) {
      return h(NButton, {
        strong: true,
        tertiary: true,
        size: 'small',
        onClick: () => onViewReport(row)
      }, { default: () => 'View Report', icon: () => h(NIcon, { component: Detail }) });
    }
  }
];

const columns = createColumns({
  onViewReport: (scan) => emit('view-report', scan)
});

</script>

<template>
  <n-data-table
    :columns="columns"
    :data="props.scanData"
    :pagination="{ pageSize: 10 }"
    :single-line="false"
  />
</template>