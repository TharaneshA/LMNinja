<script setup>
import { computed } from 'vue';
import { NIcon } from 'naive-ui';
import useChatStore from 'stores/chat.js';
import useConnectionStore from 'stores/connections.js';
import Down from '@/components/icons/Down.vue';

const chatStore = useChatStore();
const connectionStore = useConnectionStore();

const connectionOptions = computed(() => {
    return Object.values(connectionStore.serverProfile).map(conn => ({
        label: conn.name,
        value: conn.id,
    }));
});

const activeModelId = computed(() => chatStore.activeConnectionId);

const handleModelSelect = (connectionId) => {
    if (connectionId) {
        chatStore.loadAndSetActiveModel(connectionId);
    }
};
</script>

<template>
    <div class="model-selector-bar">
        <n-select
            :value="activeModelId"
            :options="connectionOptions"
            :loading="chatStore.isLoadingModel"
            placeholder="Select a model to load..."
            size="large"
            @update:value="handleModelSelect"
        >
            <template #arrow><n-icon :component="Down" /></template>
        </n-select>
    </div>
</template>

<style lang="scss" scoped>
.model-selector-bar {
    padding: 8px;
    border-bottom: 1px solid v-bind('themeVars.borderColor');
}
:deep(.n-base-selection) {
    --n-font-size: 15px !important;
    --n-padding-single: 0 32px 0 12px !important;
    --n-height: 40px !important;
    background-color: v-bind('themeVars.inputColor');
}
</style>