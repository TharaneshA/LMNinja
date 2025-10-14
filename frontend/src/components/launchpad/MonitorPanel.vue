<script setup>
import { NEmpty, NH3, NProgress, NDivider, NLog, useThemeVars } from 'naive-ui';
import useScanStore from 'stores/scan.js';

const scanStore = useScanStore();
const themeVars = useThemeVars();
</script>

<template>
  <div class="monitor-content">
    <div v-if="scanStore.scanLog.length === 0" class="empty-state">
      <n-empty size="huge" description="Configure and launch a scan to see real-time results here." />
    </div>

    <div v-else class="results-view">
      <div class="progress-section">
        <n-h3>Scan Progress</n-h3>
        <n-progress
          type="line"
          :percentage="scanStore.progress"
          :indicator-placement="'inside'"
          :processing="scanStore.isScanning"
          status="success"
        />
        <p class="progress-text">{{ scanStore.progressText }}</p>
      </div>

      <n-divider />

      <n-log 
        class="log-output" 
        :log="scanStore.formattedLog" 
        :rows="25"
        hl
        language="text"
        :theme-overrides="{ 
          color: '#000000',
        }" 
      />
    </div>
  </div>
</template>

<style lang="scss" scoped>
.monitor-content {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: v-bind('themeVars.bodyColor');
}
.empty-state {
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}
.results-view {
  height: 100%;
  display: flex;
  flex-direction: column;
}
.progress-section {
  flex-shrink: 0;
}
.progress-text {
  text-align: center;
  margin-top: 8px;
  color: v-bind('themeVars.textColor3');
}
.log-output {
  flex-grow: 1;
  overflow: hidden;
  border: 1px solid v-bind('themeVars.borderColor');
  border-radius: 6px;
  padding: 10px;

  :deep(.n-log-content) {
    font-family: 'Courier New', Courier, monospace;
    font-size: 13px;
    line-height: 1.5;
  }

  :deep(.log-line) {
    color: #FFFFFF;
  }

  :deep(.log-line--SENT) {
    color: #f0a020; /* Yellow */
  }
  :deep(.log-line--RECV) {
    color: #76b7f7; /* Light Blue */
  }
  :deep(.log-line--EVAL) {
    color: #e88080; /* Red */
  }
  :deep(.log-line--SUCCESS) {
    color: #63e2b7; /* Green */
  }
  :deep(.log-line--WARN) {
    color: #f0a020; /* Yellow */
  }
  :deep(.log-line--INFO) {
    color: #999999; /* Grey */
  }
}
</style>