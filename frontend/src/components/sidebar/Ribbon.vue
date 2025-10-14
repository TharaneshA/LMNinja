<script setup>
import { computed } from 'vue';
import useTabStore from 'stores/tab.js';
import useDialogStore from 'stores/dialog.js';
import { BrowserOpenURL } from 'wailsjs/runtime/runtime.js';
import { NIcon, useThemeVars } from 'naive-ui';
import { extraTheme } from '@/utils/extra_theme.js';
import Server from '@/components/icons/Server.vue';
import Database from '@/components/icons/Database.vue';
import Zap from '@/components/icons/Zap.vue';
import Config from '@/components/icons/Config.vue';
import Github from '@/components/icons/Github.vue';

const props = defineProps({ width: { type: Number, default: 50 } }); 
const tabStore = useTabStore();
const dialogStore = useDialogStore();
const themeVars = useThemeVars();
const exThemeVars = computed(() => extraTheme(true));
const iconSize = computed(() => Math.floor(props.width * 0.55)); 

const menuOptions = [
    { label: 'Red Teaming', key: 'browser', icon: Database },
    { label: 'LLM Connectors', key: 'server', icon: Server },
    { label: 'Launch Pad', key: 'launchpad', icon: Zap },
];


const openGithub = () => BrowserOpenURL('https://github.com/TharaneshA/LMNinja');
</script>

<template>
    <div id="app-ribbon" :style="{ width: props.width + 'px', minWidth: props.width + 'px' }" class="flex-box-v">
        <div class="ribbon-wrapper flex-box-v">
            <n-tooltip v-for="m in menuOptions" :key="m.key" :delay="300" placement="right">
                <template #trigger>
                    <div :class="{ 'ribbon-item-active': tabStore.nav === m.key }" class="ribbon-item clickable" @click="tabStore.nav = m.key">
                        <n-icon :size="iconSize"><component :is="m.icon" :stroke-width="3.5" /></n-icon>
                    </div>
                </template>
                {{ m.label }}
            </n-tooltip>
        </div>
        <div class="flex-item-expand"></div>
        <div class="nav-menu-item flex-box-v">
             <n-tooltip placement="right" :delay="300">
                <template #trigger>
                    <div class="ribbon-item clickable" @click="dialogStore.openPreferencesDialog()">
                        <n-icon :size="iconSize"><Config :stroke-width="3" /></n-icon>
                    </div>
                </template>
                Preferences
             </n-tooltip>
             <n-tooltip placement="right" :delay="300">
                <template #trigger>
                    <div class="ribbon-item clickable" @click="openGithub">
                        <n-icon :size="iconSize"><Github /></n-icon>
                    </div>
                </template>
                GitHub
             </n-tooltip>
        </div>
    </div>
</template>

<style lang="scss" scoped>
#app-ribbon {
    border-right: 1px solid v-bind('exThemeVars.splitColor');
    background-color: v-bind('exThemeVars.ribbonColor');
    color: v-bind('themeVars.textColor2');
    --wails-draggable: drag;
}
.ribbon-wrapper, .nav-menu-item {
    width: 100%;
    align-items: center;
    --wails-draggable: none;
}
.ribbon-item {
    width: 100%;
    height: 50px; // Fixed height for vertical centering
    display: flex;
    justify-content: center;
    align-items: center;
    position: relative;
    color: v-bind('themeVars.textColor3');
    &:hover { color: v-bind('themeVars.primaryColor'); }
}
.ribbon-item-active {
    color: v-bind('themeVars.primaryColor');
    &:before {
        position: absolute; width: 3px; left: 0; top: 25%; bottom: 25%;
        border-radius: 9999px; content: ''; background-color: v-bind('themeVars.primaryColor');
    }
}
</style>