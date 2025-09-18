<script setup>
import { computed, onMounted, ref } from 'vue';
import { useThemeVars } from 'naive-ui';
import { EventsOn, WindowIsMaximised, WindowToggleMaximise } from 'wailsjs/runtime/runtime.js';
import { isMacOS } from '@/utils/platform.js';
import usePreferencesStore from 'stores/preferences.js';
import useTabStore from 'stores/tab.js';
import { extraTheme } from '@/utils/extra_theme.js';
import iconUrl from '@/assets/images/icon.png';
import Ribbon from './components/sidebar/Ribbon.vue';
import ConnectionPane from './components/sidebar/ConnectionPane.vue';
import ContentPane from './components/content/ContentPane.vue';
import ResizeableWrapper from './components/common/ResizeableWrapper.vue';
import ContentValueTab from './components/content/ContentValueTab.vue';
import ContentServerPane from './components/content/ContentServerPane.vue';
import ToolbarControlWidget from './components/common/ToolbarControlWidget.vue';

const props = defineProps({ loading: Boolean });

const themeVars = useThemeVars();
const prefStore = usePreferencesStore();
const tabStore = useTabStore();
const exThemeVars = computed(() => extraTheme(true));
const maximised = ref(false);

const logoWrapperWidth = computed(() => `calc(50px + ${prefStore.behavior.asideWidth}px - 4px)`);

onMounted(async () => {
    maximised.value = await WindowIsMaximised();
});
EventsOn('window_changed', (info) => {
    if (info) maximised.value = info.maximised;
});
</script>

<template>
    <n-spin :show="props.loading" :style="{ backgroundColor: themeVars.bodyColor, height: '100vh', width: '100vw' }" :theme-overrides="{ opacitySpinning: 0.5 }">
        <div id="app-content-wrapper" class="flex-box-v">
            <div id="app-toolbar" style="height: 38px; --wails-draggable: drag" @dblclick="WindowToggleMaximise">
                <div id="app-toolbar-title" :style="{ width: logoWrapperWidth, minWidth: logoWrapperWidth }">
                    <n-space :size="5" align="center">
                        <n-avatar :size="28" :src="iconUrl" color="transparent" />
                        <div style="font-weight: 800">LMNinja</div>
                        <n-text v-if="tabStore.nav === 'browser' && tabStore.currentTabName" class="ellipsis" strong>
                            - {{ tabStore.currentTabName }}
                        </n-text>
                    </n-space>
                </div>
                <div v-show="tabStore.nav === 'browser'" class="app-toolbar-tab flex-item-expand">
                    <ContentValueTab />
                </div>
                <div class="flex-item-expand" style="min-width: 15px"></div>
                <ToolbarControlWidget v-if="!isMacOS()" :maximised="maximised" :size="38" />
            </div>

            <div id="app-content" class="flex-box-h flex-item-expand">
                <Ribbon v-model:value="tabStore.nav" />
                <div class="content-area flex-box-h flex-item-expand">
                    <ResizeableWrapper v-model:size="prefStore.behavior.asideWidth" :min-size="250">
                        <ConnectionPane class="app-side flex-item-expand" />
                    </ResizeableWrapper>
                    <ContentServerPane v-if="tabStore.nav === 'server' || (tabStore.nav === 'browser' && !tabStore.currentTab)" class="flex-item-expand" />
                    <ContentPane v-else-if="tabStore.nav === 'browser' && tabStore.currentTab" :key="tabStore.currentTab.name" :server="tabStore.currentTab.name" class="flex-item-expand" />
                </div>
            </div>
        </div>
    </n-spin>
</template>

<style lang="scss" scoped>
#app-content-wrapper { width: 100vw; height: 100vh; overflow: hidden; background-color: v-bind('themeVars.bodyColor'); }
#app-toolbar {
    background-color: v-bind('exThemeVars.titleColor');
    border-bottom: 1px solid v-bind('exThemeVars.splitColor');
    display: flex; align-items: center;
    color: v-bind('themeVars.textColorBase');
}
#app-toolbar-title {
    padding-left: v-bind('isMacOS() ? "70px" : "10px"');
    box-sizing: border-box; align-self: center;
}
.app-toolbar-tab { align-self: flex-end; margin-bottom: -1px; margin-left: 3px; overflow: hidden; }
#app-content { height: calc(100% - 38px); }
.content-area { overflow: hidden; }
.app-side {
    height: 100%;
    background-color: v-bind('exThemeVars.sidebarColor');
    border-right: 1px solid v-bind('exThemeVars.splitColor');
}
</style>