<script setup>
import { NGrid, NGi, NCard, NTabs, NTabPane, NH2, NSelect, NIcon, NButton, useThemeVars, NSpin } from 'naive-ui';
import { onMounted } from 'vue';
import useDashboardStore from 'stores/dashboard.js'; // <-- Import our new store
import StatCard from './StatCard.vue';
import VulnerabilitiesByModelChart from './VulnerabilitiesByModelChart.vue';
import ScanHistoryTable from './ScanHistoryTable.vue';
import Shield from '@/components/icons/Shield.vue';
import Server from '@/components/icons/Server.vue';
import Zap from '@/components/icons/Zap.vue';
import Checkbox from '@/components/icons/Checkbox.vue';
import Refresh from '@/components/icons/Refresh.vue';

const themeVars = useThemeVars();
const dashboardStore = useDashboardStore(); 

onMounted(() => {
    dashboardStore.fetchDashboardData();
});
</script>

<template>
  <div class="dashboard-container">
    <div class="header-bar">
      <n-h2 style="margin: 0;">Security Dashboard</n-h2>
      <div class="header-actions">
        <n-select
          value="All Time"
          :options="[{label: 'All Time', value: 'all'}]"
          style="width: 200px;"
          disabled
        />
        <n-button circle @click="dashboardStore.fetchDashboardData" :loading="dashboardStore.loading">
          <template #icon>
            <n-icon :component="Refresh" />
          </template>
        </n-button>
      </div>
    </div>

    <n-spin :show="dashboardStore.loading">
      <div class="dashboard-content">
        <n-grid cols="4" :x-gap="16" :y-gap="16" style="margin-top: 16px;">
          <n-gi>
            <StatCard label="Models Scanned" :value="dashboardStore.stats.modelsScanned" :icon="Server" />
          </n-gi>
          <n-gi>
            <StatCard label="Total Scans Run" :value="dashboardStore.stats.totalScans" :icon="Zap" />
          </n-gi>
          <n-gi>
            <StatCard label="Vulnerabilities Found" :value="dashboardStore.stats.vulnerabilitiesFound" :icon="Shield" color="#e88080" />
          </n-gi>
          <n-gi>
            <StatCard label="Overall Pass Rate" :value="`${dashboardStore.stats.overallPassRate.toFixed(1)}%`" :icon="Checkbox" color="#63e2b7" />
          </n-gi>
        </n-grid>

        <n-card style="margin-top: 24px; flex-grow: 1; display: flex; flex-direction: column;">
          <n-tabs type="line" animated default-value="summary">
            <n-tab-pane name="summary" tab="Graphical Summary">
              <n-grid cols="2" :x-gap="24" :y-gap="24">
                <n-gi span="2">
                  <VulnerabilitiesByModelChart :chart-data="dashboardStore.vulnerabilitiesByModel" />
                </n-gi>
              </n-grid>
            </n-tab-pane>
            <n-tab-pane name="history" tab="Scan History">
              <ScanHistoryTable :scan-data="dashboardStore.scanHistory" />
            </n-tab-pane>
          </n-tabs>
        </n-card>
      </div>
    </n-spin>
  </div>
</template>

<style lang="scss" scoped>
.dashboard-container {
  height: 100%;
  padding: 16px 24px;
  display: flex;
  flex-direction: column;
  background-color: v-bind('themeVars.bodyColor');
  overflow: auto;
}
.header-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.header-actions {
  display: flex;
  gap: 16px;
  align-items: center;
}
</style>