import { defineStore } from 'pinia'
import { lang } from '@/langs/index.js'
import { cloneDeep, get, map, set } from 'lodash'
import { enUS, useOsTheme, zhCN } from 'naive-ui'

// NOTE: We are removing Go backend calls for now.
// Preferences can be loaded once at startup and then managed client-side.
// This simplifies the logic and reduces backend chatter.
// We will add persistence logic later if needed.

const osTheme = useOsTheme()

const usePreferencesStore = defineStore('preferences', {
    state: () => ({
        behavior: {
            welcomed: false,
            asideWidth: 300,
        },
        general: {
            theme: 'auto',
            language: 'auto',
            fontFamily: [],
            fontSize: 14,
        },
        // We are removing Redis-specific settings like scanSize, keyIconStyle, etc.
        // We are also removing fields we will add later, like 'editor' and 'cli'.
        decoder: [], // Keep for future use, but it's empty now.
        fontList: [],
    }),
    getters: {
        // ... (Getters from the previous response can be kept, they are mostly UI logic)
        langOption() {
            const options = Object.entries(lang).map(([key, value]) => ({
                value: key,
                label: value['name'],
            }));
            options.splice(0, 0, {
                value: 'auto',
                label: 'preferences.general.system_lang', // You'll need to check your lang files for this key
            });
            return options;
        },
        currentLanguage() {
            let lang = get(this.general, 'language', 'auto');
            if (lang === 'auto') {
                const systemLang = navigator.language || navigator.userLanguage;
                lang = systemLang.split('-')[0];
            }
            return lang || 'en';
        },
        isDark() {
            const th = get(this.general, 'theme', 'auto');
            return th === 'auto' ? osTheme.value === 'dark' : th === 'dark';
        },
        themeLocale() {
            return this.currentLanguage === 'zh' ? zhCN : enUS;
        },
        generalFont() {
            const fontStyle = { fontSize: `${this.general.fontSize}px` };
            if (this.general.fontFamily && this.general.fontFamily.length > 0) {
                fontStyle['fontFamily'] = this.general.fontFamily.map(f => `"${f}"`).join(',');
            }
            return fontStyle;
        },
    },
    actions: {
        // This is now a simple client-side initialization.
        async loadPreferences() {
            // In a real app, you might load this from a Go function that reads a single config file.
            // For now, we'll just use the defaults. This removes the broken backend call.
            console.log("Preferences initialized with client-side defaults.");
        },
        
        // This would be the method to persist settings back to Go.
        async savePreferences() {
            // TODO: Create a Go function `app.SavePreferences(prefs)`
            // For now, this action doesn't do anything on the backend.
            console.log("Simulating save of preferences:", this.$state);
            return true;
        },
    },
});

export default usePreferencesStore;