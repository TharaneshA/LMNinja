<script setup>
import { computed, ref, onMounted } from 'vue'; 
import { NIcon, NSpace, NFormItem, NSelect, NButton, NH3, NInputNumber, NTooltip } from 'naive-ui';
import useConnectionStore from 'stores/connections.js';
import useScanStore from 'stores/scan.js';
import Zap from '@/components/icons/Zap.vue';
import Pause from '@/components/icons/Pause.vue';
import { GetAttackCategories } from 'wailsjs/go/app/App';

const connectionStore = useConnectionStore();
const scanStore = useScanStore();

const connectionOptions = computed(() => {
    return Object.values(connectionStore.serverProfile).map(conn => ({
        label: conn.name,
        value: conn.id,
    }));
});

const attackCategoryOptions = ref([]);

// Fetch categories when the component is created
onMounted(async () => {
    try {
        const categories = await GetAttackCategories();
        attackCategoryOptions.value = categories.map(cat => ({
            label: cat.name,
            value: cat.id,
        }));
    } catch (error) {
        $message.error(`Failed to load attack categories: ${error}`);
    }
});

// This computed property determines if the "Launch Scan" button should be clickable.
const canStartScan = computed(() => {
    return scanStore.targetModelId && scanStore.selectedCategories.length > 0 && !scanStore.isScanning;
});
</script>

<template>
  <div class="config-content">
    <n-h3>Configure Scan</n-h3>
    
    <n-space vertical :size="20">
      <n-form-item label="Target Model" label-placement="top">
        <n-select
            v-model:value="scanStore.targetModelId"
            :options="connectionOptions"
            placeholder="Select a model to scan..."
            :disabled="scanStore.isScanning"
            clearable
        />
      </n-form-item>

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

      <!-- NEW: Input for number of prompts -->
      <n-form-item label="Number of Prompts per Category" label-placement="top">
        <n-tooltip trigger="hover">
            <template #trigger>
                <n-input-number 
                    v-model:value="scanStore.scanLimit" 
                    :min="1" 
                    :max="30"
                    :disabled="scanStore.isScanning" 
                />
            </template>
            Sets how many prompts to randomly sample from each selected category.
        </n-tooltip>
      </n-form-item>
      
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