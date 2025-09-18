<script setup>
import { computed } from 'vue';
import { useThemeVars } from 'naive-ui';
import { Quit, WindowMinimise, WindowToggleMaximise } from 'wailsjs/runtime/runtime.js';
import WindowMin from '@/components/icons/WindowMin.vue';
import WindowMax from '@/components/icons/WindowMax.vue';
import WindowClose from '@/components/icons/WindowClose.vue';
import WindowRestore from '@/components/icons/WindowRestore.vue';

const themeVars = useThemeVars();
const props = defineProps({
    size: { type: Number, default: 38 },
    maximised: { type: Boolean },
});

const buttonSize = computed(() => `${props.size}px`);
</script>

<template>
    <n-space :size="0" :wrap-item="false" align="center" justify="center">
        <div class="btn-wrapper" @click="WindowMinimise"><WindowMin /></div>
        <div class="btn-wrapper" @click="WindowToggleMaximise">
            <WindowRestore v-if="maximised" />
            <WindowMax v-else />
        </div>
        <div class="btn-wrapper close-btn" @click="Quit"><WindowClose /></div>
    </n-space>
</template>

<style lang="scss" scoped>
.btn-wrapper {
    width: v-bind('buttonSize');
    height: v-bind('buttonSize');
    display: flex;
    align-items: center;
    justify-content: center;
    --wails-draggable: no-drag;
    &:hover {
        background-color: v-bind('themeVars.closeColorHover');
        cursor: pointer;
    }
    &:active {
        background-color: v-bind('themeVars.closeColorPressed');
    }
}
.close-btn:hover {
    background-color: v-bind('themeVars.errorColorHover');
}
</style>