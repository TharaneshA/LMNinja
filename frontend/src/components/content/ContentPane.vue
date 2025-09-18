<script setup>
import ModelSelectorBar from './ModelSelectorBar.vue';
import ChatInterface from './ChatInterface.vue';
import useChatStore from 'stores/chat.js';
import { useThemeVars } from 'naive-ui';
import { computed } from 'vue';

const chatStore = useChatStore();
const themeVars = useThemeVars();

// Data for the bottom bar
const modelProvider = computed(() => chatStore.activeModel?.provider || 'N/A');
const modelName = computed(() => chatStore.activeModel?.model || 'N/A');
</script>

<template>
    <div class="content-wrapper">
        <ModelSelectorBar />
        <div class="chat-area">
            <ChatInterface v-if="chatStore.isModelLoaded" />
            <n-empty v-else description="Select a model from the dropdown to begin your session." class="empty-content" />
        </div>
        <!-- Contextual Info Bar at the bottom -->
        <div class="content-footer flex-box-h">
            <n-text>Provider: {{ modelProvider }}</n-text>
            <n-divider vertical />
            <n-text>Model: {{ modelName }}</n-text>
            <div class="flex-item-expand"></div>
            <!-- Add other info like RAM usage here later -->
        </div>
    </div>
</template>

<style lang="scss" scoped>
.content-wrapper {
    display: flex;
    flex-direction: column;
    height: 100%;
    background-color: v-bind('themeVars.bodyColor');
}
.chat-area {
    flex-grow: 1;
    overflow: hidden;
}
.empty-content {
    height: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
}
.content-footer {
    border-top: 1px solid v-bind('themeVars.borderColor');
    background-color: v-bind('themeVars.cardColor');
    padding: 4px 12px;
    font-size: 12px;
    color: v-bind('themeVars.textColor3');
    height: 30px;
    align-items: center;
    flex-shrink: 0;
}
</style>