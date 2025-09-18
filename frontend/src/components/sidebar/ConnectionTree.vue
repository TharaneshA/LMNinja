<script setup>
import { h, ref } from 'vue';
import { NIcon } from 'naive-ui';
import useConnectionStore from 'stores/connections.js';
import useBrowserStore from 'stores/browser.js';
import { ConnectionType } from '@/consts/connection_type.js';
import Folder from '@/components/icons/Folder.vue';
import Server from '@/components/icons/Server.vue';
import Connect from '@/components/icons/Connect.vue';
import IconButton from '@/components/common/IconButton.vue';

const connectionStore = useConnectionStore();
const browserStore = useBrowserStore();

const props = defineProps({ filterPattern: String });
const expandedKeys = ref([]);
const selectedKeys = ref([]);

const nodeProps = ({ option }) => ({
    onClick: () => {
        selectedKeys.value = [option.key];
    },
    onDblclick: () => {
        if (option.type === ConnectionType.Server) {
            browserStore.openConnection(option.name);
        }
    },
});

const renderPrefix = ({ option }) => {
    if (option.type === ConnectionType.Group) {
        return h(NIcon, { size: 20 }, () => h(Folder, { open: expandedKeys.value.includes(option.key) }));
    }
    if (option.type === ConnectionType.Server) {
        const isConnected = browserStore.isConnected(option.name);
        return h(NIcon, { size: 20, color: isConnected ? '#4A90E2' : undefined }, () => h(Server, { inverse: isConnected }));
    }
};

const renderSuffix = ({ option }) => {
    if (option.type === ConnectionType.Server && !browserStore.isConnected(option.name)) {
        return h(IconButton, {
            icon: Connect, size: 16,
            onClick: () => browserStore.openConnection(option.name),
            tTooltip: "Connect",
        });
    }
    return null;
};
</script>

<template>
    <div class="connection-tree-wrapper">
        <!-- This v-if is critical: it prevents rendering until the store is ready and populated -->
        <n-tree
            v-if="connectionStore.connections && connectionStore.connections.length > 0"
            :data="connectionStore.connections"
            :pattern="props.filterPattern"
            :node-props="nodeProps"
            :render-prefix="renderPrefix"
            :render-suffix="renderSuffix"
            v-model:expanded-keys="expandedKeys"
            v-model:selected-keys="selectedKeys"
            block-line
            class="fill-height"
            virtual-scroll
        />
        <n-empty v-else description="No Data" class="empty-content" />
    </div>
</template>
<style scoped>
.connection-tree-wrapper { height: 100%; overflow: hidden; }
.empty-content {
    height: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
}
</style>