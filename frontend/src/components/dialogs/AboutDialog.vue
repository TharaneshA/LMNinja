<script setup>
import iconUrl from '@/assets/images/icon.png'
import useDialog from 'stores/dialog.js'
import { BrowserOpenURL } from 'wailsjs/runtime/runtime.js'
import { GetAppInfo } from 'wailsjs/go/app/App' // Corrected Import Path
import { onMounted, ref } from 'vue'

const dialogStore = useDialog()
const version = ref('')

onMounted(() => {
    // Call the new, correct Go method
    GetAppInfo().then((info) => {
        version.value = `v${info.version}`
    })
})

const onOpenSource = () => {
    BrowserOpenURL('https://github.com/TharaneshA/LMNinja') // Your repo
}

const onOpenWebsite = () => {
    // You can create a website later if you want
    BrowserOpenURL('https://github.com/TharaneshA/LMNinja')
}
</script>

<template>
    <n-modal v-model:show="dialogStore.aboutDialogVisible" :show-icon="false" preset="dialog" transform-origin="center">
        <n-space :size="10" :wrap="false" :wrap-item="false" align="center" vertical>
            <n-avatar :size="120" :src="iconUrl" color="#0000"></n-avatar>
            <div style="font-weight: bold; font-size: 18px; margin: 5px;">LMNinja</div>
            <n-text>{{ version }}</n-text>
            <n-space :size="5" :wrap="false" :wrap-item="false" align="center">
                <n-text style="cursor: pointer; text-decoration: underline;" @click="onOpenSource">Source Code</n-text>
                <n-divider vertical />
                <n-text style="cursor: pointer; text-decoration: underline;" @click="onOpenWebsite">Official Website</n-text>
            </n-space>
            <div style="font-size: 12px; color: grey;">
                Copyright © 2024 Tharanesh A. All rights reserved.
            </div>
        </n-space>
    </n-modal>
</template>