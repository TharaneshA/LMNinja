import usePreferencesStore from 'stores/preferences.js';
import { createDiscreteApi, darkTheme } from 'naive-ui';
import { darkThemeOverrides, themeOverrides } from '@/utils/theme.js';

export async function setupDiscreteApi() {
    const prefStore = usePreferencesStore();
    const configProviderProps = {
        theme: prefStore.isDark ? darkTheme : undefined,
        themeOverrides: prefStore.isDark ? darkThemeOverrides : themeOverrides,
    };
    const { message, dialog, notification } = createDiscreteApi(
        ['message', 'notification', 'dialog'], 
        { configProviderProps }
    );
    window.$message = message;
    window.$notification = notification;
    window.$dialog = dialog;
}