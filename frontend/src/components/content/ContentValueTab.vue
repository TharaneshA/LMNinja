<script setup>
import { computed } from 'vue';
import { useThemeVars } from 'naive-ui';
import useTabStore from 'stores/tab.js';
import useConnectionStore from 'stores/connections.js';
import { extraTheme } from '@/utils/extra_theme.js';
import usePreferencesStore from 'stores/preferences.js';
import Server from '@/components/icons/Server.vue';

const themeVars = useThemeVars();
const tabStore = useTabStore();
const connectionStore = useConnectionStore();
const prefStore = usePreferencesStore();

const exThemeVars = computed(() => extraTheme(true));

const tabMarkColor = (tabName) => {
  if (!tabName) return '';
  const profile = Object.values(connectionStore.serverProfile).find(p => p.name === tabName);
  return profile?.markColor || '';
};

// This is the key fix: We check if tabStore.currentTab exists before accessing its properties.
const currentTabMarkColor = computed(() => {
  return tabStore.currentTab ? tabMarkColor(tabStore.currentTab.name) : '';
});

const tabClass = (tab, idx) => {
    const isActive = tabStore.activatedIndex === idx;
    const color = tabMarkColor(tab.name);
    return [
        'value-tab',
        isActive ? 'value-tab-active' : 'value-tab-inactive',
        isActive && color ? 'value-tab-active_mark' : ''
    ];
};

const handleCloseTab = (tabIndex) => {
    const tab = tabStore.tabs[tabIndex];
    if (tab) tabStore.closeTab(tab.name);
};

const handleSwitchTab = (tabIndex) => {
    tabStore.switchTab(tabIndex);
};
</script>

<template>
    <n-tabs
        :value="tabStore.activatedIndex"
        type="card"
        :closable="true"
        size="small"
        @close="handleCloseTab"
        @update:value="handleSwitchTab"
        :theme-overrides="{
            tabGapSmallCard: '0px',
            tabColor: '#0000',
            tabBorderColor: '#0000',
            tabFontWeightActive: '800',
            tabTextColorCard: themeVars.closeIconColor,
        }"
    >
        <n-tab v-for="(t, i) in tabStore.tabs" :key="t.name" :name="i" :class="tabClass(t, i)">
            <n-space :size="5" :wrap-item="false" align="center">
                <n-icon size="18"><Server stroke-width="4" /></n-icon>
                <n-ellipsis style="max-width: 150px">{{ t.title }}</n-ellipsis>
            </n-space>
        </n-tab>
    </n-tabs>
</template>

<style lang="scss">
.value-tab {
    --wails-draggable: none;
    position: relative;
    border: 1px solid v-bind('exThemeVars.splitColor') !important;
    border-bottom: none !important;
}
.value-tab-active {
    background-color: v-bind('themeVars.tabColor') !important;
    &_mark {
        border-top: 3px solid v-bind('currentTabMarkColor') !important;
    }
}
.value-tab-inactive {
    background-color: v-bind('exThemeVars.titleColor') !important;
    &:hover { background-color: v-bind('exThemeVars.splitColor') !important; }
}
</style>