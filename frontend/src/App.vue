<script setup>
import { onMounted, ref, watch } from 'vue';
import { darkTheme } from 'naive-ui';
import { EventsOn, WindowSetDarkTheme, WindowSetLightTheme } from 'wailsjs/runtime';
import { darkThemeOverrides, themeOverrides } from '@/utils/theme.js';
import usePreferencesStore from './stores/preferences.js';
import useConnectionStore from './stores/connections.js';
import AppContent from './AppContent.vue';
import ConnectionDialog from './components/dialogs/ConnectionDialog.vue';
import GroupDialog from './components/dialogs/GroupDialog.vue';
import PreferencesDialog from './components/dialogs/PreferencesDialog.vue';
import AboutDialog from '@/components/dialogs/AboutDialog.vue';
import { i18n } from '@/utils/i18n.js';

const prefStore = usePreferencesStore();
const connectionStore = useConnectionStore();
const initializing = ref(true);

onMounted(async () => {
    // This is the correct orchestration pattern from tiny-rdm.
    try {
        // All initial data fetching happens here.
        await connectionStore.initConnections();
        // You can add other startup calls here like prefStore.loadPreferences() if needed.
    } catch (e) {
        console.error("Failed during initialization:", e);
        $message.error(`Initialization failed: ${e}`);
    } finally {
        // This flag is the gatekeeper for rendering the main UI.
        initializing.value = false;
    }
});

watch(() => prefStore.isDark, (isDark) => {
    isDark ? WindowSetDarkTheme() : WindowSetLightTheme();
}, { immediate: true });

watch(() => prefStore.currentLanguage, (lang) => {
    i18n.global.locale.value = lang;
});
</script>

<template>
    <n-config-provider
        :theme="prefStore.isDark ? darkTheme : undefined"
        :theme-overrides="prefStore.isDark ? darkThemeOverrides : themeOverrides"
        class="fill-height"
    >
        <n-dialog-provider>
            <!-- Pass the loading flag to AppContent -->
            <AppContent :loading="initializing" />
            <ConnectionDialog />
            <GroupDialog />
            <PreferencesDialog />
            <AboutDialog />
        </n-dialog-provider>
    </n-config-provider>
</template>