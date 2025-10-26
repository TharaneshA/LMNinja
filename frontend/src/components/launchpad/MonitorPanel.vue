<script setup>
import { computed, ref, nextTick, watch } from 'vue';
import { NEmpty, NH3, NProgress, NDivider, useThemeVars } from 'naive-ui';
import useScanStore from 'stores/scan.js';

const scanStore = useScanStore();
const themeVars = useThemeVars();


const logContainerRef = ref(null);

const coloredLogHtml = computed(() => {
    return scanStore.scanLog.map(log => {
        const type = log.type.toUpperCase();
        
        const sanitizedContent = log.content.replace(/</g, "&lt;").replace(/>/g, "&gt;");
        
        const contentClass = `log-content log-content--${log.verdict || 'default'}`;

        return `<div class="log-line log-line--${type}"><span class="log-tag">[${type}]</span><span class="${contentClass}">${sanitizedContent}</span></div>`;
    }).join('');
});

watch(() => scanStore.scanLog.length, () => {
    nextTick(() => {
        const container = logContainerRef.value;
        if (container) {
            container.scrollTop = container.scrollHeight;
        }
    });
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


      <div class="log-output-container" ref="logContainerRef">
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
  overflow-y: auto;
  background-color: #000000;
  border: 1px solid v-bind('themeVars.borderColor');
  border-radius: 6px;
  padding: 10px;
  font-family: 'Courier New', Courier, monospace;
  font-size: 13px;
  color: #FFFFFF;
  user-select: text;
  cursor: text;
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
    display: flex;
    align-items: center;
}

:deep(.log-tag) {
    font-weight: bold;
    display: inline-block;
    width: 80px;
    flex-shrink: 0;
}

:deep(.log-content) {
    flex-grow: 1;
}

/* Tag Colors */
:deep(.log-line--SENT .log-tag) { color: #f0a020; }
:deep(.log-line--RECV .log-tag) { color: #76b7f7; }
:deep(.log-line--INFO .log-tag) { color: #999999; }
:deep(.log-line--WARN .log-tag) { color: #f0a020; }
:deep(.log-line--ERROR .log-tag) { color: #e88080; }
:deep(.log-line--SUCCESS .log-tag) { color: #63e2b7; }
:deep(.log-line--EVAL .log-tag) { color: #bd93f9; }
:deep(.log-line--EXPLAIN .log-tag) { color: #8be9fd; }

:deep(.log-content--SUCCESSFUL_ATTACK) { color: #e88080; font-weight: bold; }
:deep(.log-content--ATTACK_FAILED) { color: #63e2b7; }

:deep(.verdict-icon) {
    display: inline-block;
    width: 10px;
    height: 10px;
    border-radius: 50%;
    margin-left: 8px;
    flex-shrink: 0;
}
:deep(.verdict--success) { background-color: #63e2b7; } /* Green */
:deep(.verdict--fail) { background-color: #e88080; }    /* Red */
</style>