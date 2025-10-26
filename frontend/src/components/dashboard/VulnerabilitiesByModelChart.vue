<script setup>
import { Bar } from 'vue-chartjs';
import { Chart as ChartJS, Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale, Colors } from 'chart.js';

ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale, Colors);

const props = defineProps({
  chartData: {
    type: Object,
    required: true,
  },
});

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  indexAxis: 'y', // Makes it a horizontal bar chart
  scales: {
    x: {
      stacked: true,
      grid: { color: 'rgba(255, 255, 255, 0.1)' },
    },
    y: {
      stacked: true,
      grid: { display: false },
    },
  },
  plugins: {
    legend: {
        position: 'bottom',
    },
    title: {
        display: true,
        text: 'Vulnerabilities by Model',
    }
  }
};
</script>

<template>
  <div style="height: 300px;">
    <Bar :data="props.chartData" :options="chartOptions" />
  </div>
</template>