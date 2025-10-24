<script setup>
import { computed } from 'vue';
import { NEmpty, NH3, NProgress, NDivider, NLog, useThemeVars } from 'naive-ui';
import useScanStore from 'stores/scan.js';

const scanStore = useScanStore();
const themeVars = useThemeVars();


const coloredLogHtml = computed(() => {
    return scanStore.scanLog.map(log => {
        const type = log.type.toUpperCase();
        
        const sanitizedContent = log.content.replace(/</g, "&lt;").replace(/>/g, "&gt;");
        
        // Wrap the line in a div with a class for styling.
        return `<div class="log-line log-line--${type}"><span class="log-tag">[${type}]</span> <span class="log-content">${sanitizedContent}</span></div>`;
    }).join(''); // Join with an empty string, as each line is a div block.
});
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

      <div class="log-output-container">
          <div class="log-output" v-html="coloredLogHtml"></div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.monitor-content, .results-view, .empty-state {
  height: 100%;
  display: flex;
  flex-direction: column;
}
.monitor-content {
  background-color: v-bind('themeVars.bodyColor');
}
.empty-state {
  justify-content: center;
  align-items: center;
}
.progress-section {
  flex-shrink: 0;
}
.progress-text {
  text-align: center;
  margin-top: 8px;
  color: v-bind('themeVars.textColor3');
}

.log-output-container {
  flex-grow: 1;
  overflow-y: auto; /* The container handles scrolling */
  background-color: #000000; /* Pure black background */
  border: 1px solid v-bind('themeVars.borderColor');
  border-radius: 6px;
  padding: 10px;
  font-family: 'Courier New', Courier, monospace;
  font-size: 13px;
  color: #FFFFFF; /* Default text color is white */
}

.log-output {
    display: flex;
    flex-direction: column;
    gap: 4px;
}

:deep(.log-line) {
    line-height: 1.5;
    white-space: pre-wrap;
    word-break: break-all;
}

:deep(.log-tag) {
    font-weight: bold;
    display: inline-block;
    width: 80px; /* Aligns the tags nicely */
    flex-shrink: 0;
}

:deep(.log-line--SENT .log-tag) { color: #f0a020; } /* Yellow */
:deep(.log-line--RECV .log-tag) { color: #76b7f7; } /* Light Blue */
:deep(.log-line--EVAL .log-tag) { color: #e88080; } /* Red */
:deep(.log-line--UNSAFE .log-tag) { color: #e88080; } /* Red */
:deep(.log-line--SUCCESS .log-tag) { color: #63e2b7; } /* Green */
:deep(.log-line--SAFE .log-tag) { color: #63e2b7; } /* Green */
:deep(.log-line--WARN .log-tag) { color: #f0a020; } /* Yellow */
:deep(.log-line--INFO .log-tag) { color: #999999; } /* Grey */
:deep(.log-line--ERROR .log-tag) { color: #e88080; } /* Red */

</style>