import { defineStore } from 'pinia';
import useTabStore from 'stores/tab.js';
import useChatStore from 'stores/chat.js';
import useConnectionStore from 'stores/connections.js';

const useBrowserStore = defineStore('browser', {
    state: () => ({
        connectedServers: new Set(),
    }),
    getters: {
        anyConnectionOpened: (state) => state.connectedServers.size > 0,
        isConnected: (state) => (serverName) => state.connectedServers.has(serverName),
    },
    actions: {
        async openConnection(name) {
            const connStore = useConnectionStore();
            // Find the profile from the flat list using the name
            const profile = Object.values(connStore.serverProfile).find(p => p.name === name);
            
            if (!profile) {
                $message.error(`Connection profile for "${name}" not found.`);
                return;
            }

            // Mark as connected and create a tab for it
            this.connectedServers.add(name);
            const tabStore = useTabStore();
            tabStore.upsertTab({ server: name, forceSwitch: true });
            
            // **THE CRITICAL LINK**: Tell the chat store to load this connection's model.
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