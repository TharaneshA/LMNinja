import { createApp, nextTick } from 'vue';
import { createPinia } from 'pinia';
import App from './App.vue';
import './styles/style.scss';
import { i18n } from '@/utils/i18n.js';
import { setupDiscreteApi } from '@/utils/discrete.js';
import usePreferencesStore from 'stores/preferences.js';
import { loadEnvironment } from '@/utils/platform.js';

async function setupApp() {
    const app = createApp(App);
    app.use(i18n);
    app.use(createPinia());

    await loadEnvironment();
    const prefStore = usePreferencesStore();
    await prefStore.loadPreferences(); // Load preferences early
    await setupDiscreteApi(); // Setup globals like $message

    app.config.errorHandler = (err, instance, info) => {
        nextTick().then(() => {
            try {
                const content = err.toString();
                $notification.error(content, {
                    title: 'Application Error',
                    meta: 'See console for more detail (Ctrl+Shift+I)',
                });
                console.error("Vue Error:", err, info);
            } catch (e) {
                // Failsafe
            }
        });
    };
    
    app.mount('#app');
}

setupApp();