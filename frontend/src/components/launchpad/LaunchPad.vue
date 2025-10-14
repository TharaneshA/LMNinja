<script setup>
import { computed } from 'vue';
import { useThemeVars } from 'naive-ui';
import { extraTheme } from '@/utils/extra_theme.js';
import ConfigurationPanel from './ConfigurationPanel.vue';
import MonitorPanel from './MonitorPanel.vue';

const themeVars = useThemeVars();
const exThemeVars = computed(() => extraTheme(true));
</script>

<template>
  <div class="launch-pad-container">
    <div class="config-panel">
      <ConfigurationPanel />
    </div>

    <div class="monitor-panel">
      <MonitorPanel />
    </div>
  </div>
</template>

<style lang="scss" scoped>
.launch-pad-container {
  display: flex;
  height: 100%;
  width: 100%;
  overflow: hidden;
  background-color: v-bind('themeVars.bodyColor');
}

.config-panel {
  width: 350px;
  min-width: 300px;
  flex-shrink: 0;
  border-right: 1px solid v-bind('themeVars.borderColor');
  padding: 16px;
  display: flex;
  flex-direction: column;
  background-color: v-bind('exThemeVars.sidebarColor');
}

.monitor-panel {
  flex-grow: 1;
  padding: 16px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background-color: v-bind('themeVars.bodyColor');
}
</style>