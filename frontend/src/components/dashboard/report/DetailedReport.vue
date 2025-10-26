<script setup>
import { ref, onMounted } from 'vue';
import { NButton, NIcon, NPageHeader, NSpin, NGrid, NGi, NCard, NStatistic, NDataTable, NDivider, NAlert, useThemeVars } from 'naive-ui';
import { GetScanResults } from 'wailsjs/go/app/App';
import html2pdf from 'html2pdf.js';
import AttentionHeatmap from './AttentionHeatmap.vue';
import Download from '@/components/icons/Download.vue';
import { computed } from 'vue';

const props = defineProps({
  scan: Object,
});
const emit = defineEmits(['back']);

const themeVars = useThemeVars();
const loading = ref(true);
const results = ref([]);
const selectedResult = ref(null);

const columns = [
  { title: 'Prompt', key: 'prompt', ellipsis: { tooltip: true } },
  { title: 'Response', key: 'response', ellipsis: { tooltip: true } },
  { title: 'Verdict', key: 'verdict' },
  { title: 'Score', key: 'score', width: 100 },
];

onMounted(async () => {
  try {
    const rawResults = await GetScanResults(props.scan.id);
    results.value = rawResults.map(r => {
      const evalData = JSON.parse(r.evaluationJson);
      return {
        ...r,
        verdict: evalData.verdict,
        score: `${(evalData.score * 100).toFixed(1)}%`,
        evaluation: evalData,
      };
    });
  } catch (error) {
    $message.error(`Failed to load scan details: ${error}`);
  } finally {
    loading.value = false;
  }
});

const rowProps = (row) => ({
  style: 'cursor: pointer;',
  onClick: () => { selectedResult.value = row; }
});

const vulnerabilities = computed(() => results.value.filter(r => r.verdict === 'SUCCESSFUL_ATTACK'));
const weakCategories = computed(() => {
    return `Analysis shows the model is most vulnerable to the types of attacks run in this scan.`;
});

const exportToPdf = () => {
  const element = document.getElementById('pdf-export-content');
  const opt = {
    margin: 15,
    filename: `LMNinja_Report_${props.scan.targetModelName.replace(' ', '_')}.pdf`,
    image: { type: 'jpeg', quality: 0.98 },
    html2canvas: { scale: 2 },
    jsPDF: { unit: 'mm', format: 'a4', orientation: 'portrait' }
  };
  html2pdf().from(element).set(opt).save();
};
</script>

<template>
  <div class="detailed-report-container">
    <n-page-header @back="emit('back')">
      <template #title>Scan Report: {{ scan.targetModelName }}</template>
      <template #subtitle>{{ new Date(scan.startTime).toLocaleString() }}</template>
      <template #extra>
        <n-button @click="exportToPdf" type="primary" :loading="isExporting">
          <template #icon><n-icon :component="Download" /></template>
          Export as PDF
        </n-button>
      </template>
    </n-page-header>

    <n-spin :show="loading">
      <div class="ui-content">
        <n-grid cols="4" :x-gap="16" :y-gap="16" style="margin-top: 24px;">
          <n-gi><n-statistic label="Total Prompts" :value="scan.totalPrompts" /></n-gi>
          <n-gi><n-statistic label="Vulnerabilities Found" :value="scan.vulnerabilitiesFound" /></n-gi>
          <n-gi><n-statistic label="Pass Rate" :value="scan.passRate" /></n-gi>
          <n-gi><n-statistic label="Status" :value="scan.status" /></n-gi>
        </n-grid>
        <n-divider />
        <n-alert title="Instructions" type="info" :bordered="false" style="margin-bottom: 24px;">
          Click a row to see the detailed analysis, including the Attention Heatmap.
        </n-alert>
        <n-data-table :columns="columns" :data="results" :row-key="row => row.id" :row-props="rowProps" />
        <template v-if="selectedResult">
          <n-divider />
          <h3 class="n-h3">Detailed Analysis for Prompt #{{ selectedResult.id }}</h3>
          <AttentionHeatmap 
            :prompt="selectedResult.prompt" 
            :explainability="selectedResult.evaluation.explainability" 
          />
        </template>
      </div>
    </n-spin>

    <div id="pdf-export-wrapper">
        <div id="pdf-export-content">
            <h1>LMNinja Security Scan Report</h1>
            <hr>
            <h2>Scan Summary</h2>
            <p><strong>Target Model:</strong> {{ scan.targetModelName }}</p>
            <p><strong>Scan Date:</strong> {{ new Date(scan.startTime).toLocaleString() }}</p>
            <p><strong>Total Prompts Tested:</strong> {{ scan.totalPrompts }}</p>
            <p><strong>Vulnerabilities Found:</strong> {{ scan.vulnerabilitiesFound }}</p>
            <p><strong>Final Pass Rate:</strong> {{ scan.passRate }}</p>
            <hr>
            <h2>Executive Summary & Recommendations</h2>
            <p>The scan identified <strong>{{ vulnerabilities.length }}</strong> successful attack(s) out of {{ results.length }} prompts. {{ weakCategories }} It is recommended to investigate the failed prompts below and consider further fine-tuning or implementing stricter guardrails for the identified weak areas.</p>
            <hr>
            <h2>Detailed Vulnerability Report</h2>
            <div v-if="vulnerabilities.length > 0">
                <div v-for="result in vulnerabilities" :key="result.id" class="pdf-item">
                    <h3>Prompt #{{ result.id }} - FAILED</h3>
                    <h4>Verdict: {{ result.verdict.replace('_', ' ') }}</h4>
                    <p><strong>Distraction Score:</strong> {{ (result.evaluation.explainability.distraction_score * 100).toFixed(2) }}%</p>
                    <h4>Attention Heatmap:</h4>
                    <AttentionHeatmap :prompt="result.prompt" :explainability="result.evaluation.explainability" />
                    <h4 style="margin-top: 15px;">Model Response:</h4>
                    <p class="pdf-response">{{ result.response }}</p>
                </div>
            </div>
            <div v-else>
                <p><strong>No vulnerabilities were found in this scan.</strong></p>
            </div>
        </div>
    </div>
  </div>
</template>

<style scoped>
.detailed-report-container { padding: 16px 24px; height: 100%; overflow-y: auto; }
h2 { margin: 0 0 10px 0; font-size: 24px; font-weight: 600; }


#pdf-export-wrapper {
  position: absolute;
  left: -9999px;
  top: -9999px;
  width: 210mm; 
}
#pdf-export-content {
    color: #333;
    background-color: #fff;
    padding: 15mm;
    box-sizing: border-box;
}
.pdf-item {
    page-break-inside: avoid;
    margin-bottom: 20px;
    border: 1px solid #ccc;
    padding: 10px;
    border-radius: 4px;
}
.pdf-response {
    background-color: #f5f5f5;
    padding: 10px;
    border-radius: 4px;
    white-space: pre-wrap;
    word-break: break-all;
}
</style>