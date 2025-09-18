import { defineStore } from 'pinia';
import { ConnectionType } from '@/consts/connection_type.js';
import { GetConnections, SaveConnection, DeleteConnection } from 'wailsjs/go/app/App';
import { get } from 'lodash';

const useConnectionStore = defineStore('connections', {
    state: () => ({
        connections: [],     // Tree structure for the sidebar UI
        serverProfile: {},   // Flat map for quick access: { [connId]: connMeta }
    }),
    getters: {
        groups() {
            return this.connections
                .filter(c => c.type === ConnectionType.Group)
                .map(g => g.label);
        },
    },
    actions: {
        processConnectionData(rawConnections) {
            const tree = [];
            const groups = {};
            const profiles = {};

            if (!Array.isArray(rawConnections)) {
                console.error("Connection data is not an array:", rawConnections);
                return { tree: [], profiles: {} };
            }

            rawConnections.forEach(conn => {
                profiles[conn.id] = { ...conn };

                const connectionNode = {
                    key: (conn.group ? conn.group + '/' : '/') + conn.name,
                    label: conn.name,
                    name: conn.name,
                    type: ConnectionType.Server,
                    ...conn,
                };

                if (conn.group) {
                    if (!groups[conn.group]) {
                        groups[conn.group] = {
                            key: conn.group + '/',
                            label: conn.group,
                            type: ConnectionType.Group,
                            children: [],
                        };
                    }
                    groups[conn.group].children.push(connectionNode);
                } else {
                    tree.push(connectionNode);
                }
            });

            Object.values(groups).forEach(group => tree.push(group));
            
            this.connections = tree;
            this.serverProfile = profiles;
        },

        async initConnections(force = false) {
            if (!force && this.connections.length > 0) return;
            try {
                const rawConnections = await GetConnections() || [];
                this.processConnectionData(rawConnections);
            } catch (error) {
                console.error("Failed to initialize connections:", error);
                this.connections = [];
                this.serverProfile = {};
                throw error; // Re-throw to be caught by the startup handler
            }
        },
        
        async saveConnection(connectionData, apiKey) {
            try {
                // Your Go backend now expects the full metadata object.
                // If it's a new connection, the ID will be empty. Go will generate it.
                const newMeta = {
                    id: connectionData.id || '', // Pass existing ID for edits
                    name: connectionData.name,
                    group: connectionData.group,
                    provider: connectionData.provider,
                    model: connectionData.model,
                    // CreatedAt will be set by the backend for new connections
                };
                const updatedList = await SaveConnection(newMeta, apiKey);
                this.processConnectionData(updatedList);
                return { success: true };
            } catch (e) {
                console.error("Error saving connection:", e);
                return { success: false, msg: e.toString() };
            }
        },

        async deleteConnection(id) {
             try {
                const updatedList = await DeleteConnection(id);
                this.processConnectionData(updatedList);
                return { success: true };
            } catch (e) {
                console.error("Error deleting connection:", e);
                return { success: false, msg: e.toString() };
            }
        },
    },
});

export default useConnectionStore;