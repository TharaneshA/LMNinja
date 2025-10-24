import { defineStore } from 'pinia';
import { ConnectionType } from '@/consts/connection_type.js';
import { GetConnections, SaveConnection, DeleteConnection, RenameConnection } from 'wailsjs/go/app/App';

const useConnectionStore = defineStore('connections', {
    state: () => ({
        connections: [],
        serverProfile: {},
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
            // --- DEBUG STEP 1 ---
            console.log("[Connections Store] processConnectionData called with:", rawConnections);

            const tree = [];
            const groups = {};
            const profiles = {};

            if (!Array.isArray(rawConnections)) {
                console.error("[Connections Store] ERROR: rawConnections is not an array!");
                this.connections = [];
                this.serverProfile = {};
                return;
            }

            rawConnections.forEach(conn => {
                profiles[conn.id] = { ...conn };
                const connectionNode = {
                    key: conn.id,
                    label: conn.name,
                    name: conn.name,
                    type: ConnectionType.Server,
                    ...conn,
                };
                if (conn.group) {
                    if (!groups[conn.group]) {
                        groups[conn.group] = {
                            key: conn.group,
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
            
            Object.keys(groups).sort().forEach(groupName => {
                tree.push(groups[groupName]);
            });
            
            // --- DEBUG STEP 2 ---
            console.log("[Connections Store] Processed tree structure:", JSON.parse(JSON.stringify(tree)));
            console.log("[Connections Store] Processed server profiles:", JSON.parse(JSON.stringify(profiles)));

            this.connections = tree;
            this.serverProfile = profiles;
        },

        async initConnections(force = false) {
            console.log("[Connections Store] ACTION: initConnections called.");
            if (!force && this.connections.length > 0) {
                console.log("[Connections Store] Connections already exist, skipping init.");
                return;
            }
            try {
                const rawConnections = await GetConnections() || [];
                this.processConnectionData(rawConnections);
            } catch (error) {
                console.error("[Connections Store] FATAL: Failed to initialize connections:", error);
                throw error;
            }
        },
        
        
        async saveConnection(connectionData, apiKey) {
            try {
                const updatedList = await SaveConnection(connectionData, apiKey);
                this.processConnectionData(updatedList);
                return { success: true };
            } catch (e) {
                return { success: false, msg: e.toString() };
            }
        },

        async deleteConnection(id) {
             try {
                const updatedList = await DeleteConnection(id);
                this.processConnectionData(updatedList);
                return { success: true };
            } catch (e) {
                return { success: false, msg: e.toString() };
            }
        },

        async renameConnection(id, newName) {
            try {
                const updatedList = await RenameConnection(id, newName);
                this.processConnectionData(updatedList);
                return { success: true };
            } catch (e) {
                return { success: false, msg: e.toString() };
            }
        },
    },
});

export default useConnectionStore;