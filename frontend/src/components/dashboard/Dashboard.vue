<script setup>
import { NGrid, NGi, NCard, NTabs, NTabPane, NH2, NSelect, NIcon, NButton, useThemeVars } from 'naive-ui';
import StatCard from './StatCard.vue';
import VulnerabilitiesByModelChart from './VulnerabilitiesByModelChart.vue';
import ScanHistoryTable from './ScanHistoryTable.vue';
import Shield from '@/components/icons/Shield.vue';
import Server from '@/components/icons/Server.vue';
import Zap from '@/components/icons/Zap.vue';
import Checkbox from '@/components/icons/Checkbox.vue';
import Refresh from '@/components/icons/Refresh.vue';

const themeVars = useThemeVars();
</script>

<template>
  <div class="dashboard-container">
    <!-- 1. Top Header Bar -->
    <div class="header-bar">
      <n-h2 style="margin: 0;">Security Dashboard</n-h2>
      <div class="header-actions">
        <n-select
          value="Last 7 Days"
          :options="[{label: 'Last 7 Days', value: 7}, {label: 'Last 30 Days', value: 30}]"
          style="width: 200px;"
        />
        <n-button circle>
          <template #icon>
            <n-icon :component="Refresh" />
          </template>
        </n-button>
      </div>
    </div>

    <!-- 2. Stat Cards -->
    <n-grid cols="4" :x-gap="16" :y-gap="16" style="margin-top: 16px;">
      <n-gi>
        <StatCard label="Models Tested" value="8" :icon="Server" />
      </n-gi>
      <n-gi>
        <StatCard label="Total Scans Run" value="27" :icon="Zap" />
      </n-gi>
      <n-gi>
        <StatCard label="Vulnerabilities Found" value="154" :icon="Shield" color="#e88080" />
      </n-gi>
      <n-gi>
        <StatCard label="Overall Pass Rate" value="92.3%" :icon="Checkbox" color="#63e2b7" />
      </n-gi>
    </n-grid>

    <!-- 3. Main Dashboard View -->
    <n-card style="margin-top: 24px; flex-grow: 1; display: flex; flex-direction: column;">
      <n-tabs type="line" animated default-value="summary">
        <n-tab-pane name="summary" tab="Graphical Summary">
          <n-grid cols="2" :x-gap="24" :y-gap="24">
            <n-gi span="2">
              <VulnerabilitiesByModelChart />
            </n-gi>
            <!-- Donut chart can be added here later -->
          </n-grid>
        </n-tab-pane>
        <n-tab-pane name="history" tab="Scan History">
          <ScanHistoryTable />
        </n-tab-pane>
      </n-tabs>
    </n-card>
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