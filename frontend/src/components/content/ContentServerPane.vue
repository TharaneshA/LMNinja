<script setup>
import AddLink from '@/components/icons/AddLink.vue';
import useDialogStore from 'stores/dialog.js';
import useConnectionStore from 'stores/connections.js';

const dialogStore = useDialogStore();
const connectionStore = useConnectionStore();
</script>

<template>
    <div class="content-container flex-box-v">
        <n-empty v-if="connectionStore.connections.length === 0" description="No LLM server added">
            <template #extra>
                <n-button @click="dialogStore.openNewDialog()">
                    <template #icon><n-icon :component="AddLink" /></template>
                    Add Connection
                </n-button>
            </template>
        </n-empty>
        <n-empty v-else description="Select and open a connection from the left panel">
             <template #extra>
                <n-button @click="dialogStore.openNewDialog()">
                    <template #icon><n-icon :component="AddLink" /></template>
                    Add Connection
                </n-button>
            </template>
        </n-empty>
    </div>
</template>

<style lang="scss" scoped>
.content-container { justify-content: center; height: 100%; }
</style>