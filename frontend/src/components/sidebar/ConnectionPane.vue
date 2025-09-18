<script setup>
import { ref } from 'vue';
import { useThemeVars } from 'naive-ui';
import useDialogStore from 'stores/dialog.js';
import ConnectionTree from './ConnectionTree.vue';
import IconButton from '@/components/common/IconButton.vue';
import AddLink from '@/components/icons/AddLink.vue';
import AddGroup from '@/components/icons/AddGroup.vue';

const dialogStore = useDialogStore();
const themeVars = useThemeVars(); 
const filterPattern = ref('');
</script>

<template>
    <div class="nav-pane-container flex-box-v">
        <div class="pane-header" :style="{ color: themeVars.textColor1 }">LLM Connectors</div>
        
        <ConnectionTree :filter-pattern="filterPattern" />
        
        <div class="nav-pane-bottom nav-pane-func flex-box-h">
            <icon-button 
                :icon="AddLink" 
                size="20" 
                t-tooltip="New Connection" 
                @click="dialogStore.openNewDialog()"
            />
            <icon-button 
                :icon="AddGroup" 
                size="20" 
                t-tooltip="New Group" 
                @click="dialogStore.openNewGroupDialog()"
            />
            <n-divider vertical />
            <n-input v-model:value="filterPattern" :placeholder="$t('interface.filter')" clearable size="small" />
        </div>
    </div>
</template>

<style lang="scss" scoped>
.nav-pane-container { 
    height: 100%; 
    overflow: hidden; 
}
.pane-header {
    padding: 8px 12px;
    font-weight: bold;
    font-size: 14px;
    border-bottom: 1px solid v-bind('themeVars.borderColor');
}
.nav-pane-bottom {
    border-top: 1px solid v-bind('themeVars.borderColor');
    padding: 5px;
    align-items: center;
    gap: 5px;
}
</style>