<script setup>
import { computed, ref } from 'vue';
import { NIcon, NSpace, NFormItem, NSelect, NButton, NH3 } from 'naive-ui';
import useConnectionStore from 'stores/connections.js';
import useScanStore from 'stores/scan.js';
import Zap from '@/components/icons/Zap.vue';
import Pause from '@/components/icons/Pause.vue';

const connectionStore = useConnectionStore();
const scanStore = useScanStore();

// Get the list of saved models for the dropdown
const connectionOptions = computed(() => {
    return Object.values(connectionStore.serverProfile).map(conn => ({
        label: conn.name,
        value: conn.id,
    }));
});

// Placeholder for attack categories.
const attackCategoryOptions = ref([
    { label: 'Jailbreaking', value: 'jailbreaking' },
    { label: 'Prompt Injection', value: 'prompt_injection' },
    { label: 'PII Leaks (Simulated)', value: 'pii_leaks' },
    { label: 'Role Playing Attacks', value: 'role_playing' },
]);

// This computed property determines if the "Launch Scan" button should be clickable.
const canStartScan = computed(() => {
    return scanStore.targetModelId && scanStore.selectedCategories.length > 0 && !scanStore.isScanning;
});
</script>

<template>
  <div class="config-content">
    <n-h3>Configure Scan</n-h3>
    
    <n-space vertical :size="20">
      <!-- Target Model Selector -->
      <n-form-item label="Target Model" label-placement="top">
        <n-select
            v-model:value="scanStore.targetModelId"
            :options="connectionOptions"
            placeholder="Select a model to scan..."
            :disabled="scanStore.isScanning"
            clearable
        />
      </n-form-item>

      <!-- Attack Categories Selector -->
      <n-form-item label="Attack Categories" label-placement="top">
        <n-select
            v-model:value="scanStore.selectedCategories"
            :options="attackCategoryOptions"
            multiple
            placeholder="Select attack categories..."
            :disabled="scanStore.isScanning"
            clearable
        />
      </n-form-item>
      
      <!-- Judge Model Selector (Future Feature) -->
      <n-form-item label="Judge / Evaluator Model" label-placement="top">
         <n-select
            value="NinjaGuard (Default)"
            disabled
            placeholder="Select an evaluator..."
        />
      </n-form-item>
    </n-space>

    <!-- Launch/Stop Button Container -->
    <div class="launch-button-container">
      <n-button 
        v-if="!scanStore.isScanning"
        type="primary" 
        size="large" 
        block 
        @click="scanStore.startScan"
        :disabled="!canStartScan"
      >
        <template #icon>
          <n-icon :component="Zap" />
        </template>
        Launch Scan
      </n-button>
      <n-button 
        v-else
        type="error" 
        size="large" 
        block 
        @click="scanStore.stopScan"
      >
        <template #icon>
          <n-icon :component="Pause" />
        </template>
        Stop Scan
      </n-button>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.config-content {
  display: flex;
  flex-direction: column;
  height: 100%;
}
.launch-button-container {
  margin-top: auto;
  padding-top: 20px;
}
</style>