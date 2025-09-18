import { defineStore } from 'pinia';
import useTabStore from 'stores/tab.js';
import useChatStore from 'stores/chat.js';
import useConnectionStore from 'stores/connections.js';

const useBrowserStore = defineStore('browser', {
    state: () => ({
        // Using a Set for efficient tracking of "connected" sessions
        connectedServers: new Set(),
    }),
    getters: {
        anyConnectionOpened(state) {
            return state.connectedServers.size > 0;
        },
        isConnected: (state) => (serverName) => {
            return state.connectedServers.has(serverName);
        },
    },
    actions: {
        async openConnection(name) {
            const connStore = useConnectionStore();
            const profile = Object.values(connStore.serverProfile).find(p => p.name === name);
            
            if (!profile) {
                $message.error(`Connection profile for "${name}" not found.`);
                return;
            }

            // In LMNinja, "connecting" means creating a workspace/tab for it.
            // The actual model loading happens when selected in the chat view.
            this.connectedServers.add(name);

            const tabStore = useTabStore();
            tabStore.upsertTab({ server: name, forceSwitch: true });
            
            // Pre-emptively load the model for the new tab
            const chatStore = useChatStore();
            await chatStore.loadAndSetActiveModel(profile.id);
        },

        async closeConnection(name) {
            if (!this.isConnected(name)) return;
            
            const tabStore = useTabStore();
            const chatStore = useChatStore();

            this.connectedServers.delete(name);
            tabStore.removeTabByName(name);

            // If the closed tab was the active one, clear the active model
            if (chatStore.activeModel?.name === name) {
                chatStore.loadAndSetActiveModel(null);
            }
        },
    },
});

export default useBrowserStore;