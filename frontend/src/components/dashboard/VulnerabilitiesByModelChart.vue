<script setup>
import { Bar } from 'vue-chartjs';
import { Chart as ChartJS, Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale, Colors } from 'chart.js';

ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale, Colors);

// Mock data - in the future, this will come from your backend
const chartData = {
  labels: ['Llama-3-8B', 'Gemini-1.5-Flash', 'Qwen3-0.6B', 'Mistral-7B'],
  datasets: [
    {
      label: 'Jailbreaking',
      backgroundColor: '#e88080', // Red
      data: [45, 12, 5, 22],
    },
    {
      label: 'Prompt Injection',
      backgroundColor: '#f0a020', // Yellow
      data: [20, 30, 2, 15],
    },
    {
      label: 'PII Leaks',
      backgroundColor: '#76b7f7', // Blue
      data: [5, 8, 1, 3],
    },
  ],
};

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
    <Bar :data="chartData" :options="chartOptions" />
  </div>
</template>