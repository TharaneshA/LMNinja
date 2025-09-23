import { defineStore } from 'pinia';
import { LoadModel, SendMessage } from 'wailsjs/go/app/App';

const useChatStore = defineStore('chat', {
    state: () => ({
        activeConnectionId: null,
        activeModel: null,
        messages: [],
        isLoadingModel: false,
        isStreamingResponse: false,
    }),
    getters: {
        isModelLoaded: (state) => !!state.activeModel,
    },
    actions: {
        async loadAndSetActiveModel(connectionId) {
            // --- DEBUG STEP 0 ---
            console.log(`[ChatStore] ACTION: loadAndSetActiveModel called with ID: ${connectionId}`);

            if (!connectionId) {
                console.log("[ChatStore] ID is null, clearing active model.");
                this.activeConnectionId = null;
                this.activeModel = null;
                this.messages = [];
                return;
            }
            if (this.isLoadingModel) {
                console.log("[ChatStore] Already loading a model, aborting.");
                return;
            }

            this.isLoadingModel = true;
            this.messages = [];
            
            // --- DEBUG STEP 1 ---
            console.log("[ChatStore] Creating loading message...");
            const loadingMessage = $message.loading("Loading model...", { duration: 0 });
            console.log("[ChatStore] Loading message should be visible now.");

            let loadedModelMeta = null;

            try {
                // --- DEBUG STEP 2 ---
                console.log("[ChatStore] AWAITING: Calling Go backend 'LoadModel'...");
                loadedModelMeta = await LoadModel(connectionId);
                // --- DEBUG STEP 3 ---
                console.log("[ChatStore] AWAITED: 'LoadModel' returned successfully. Model meta:", loadedModelMeta);

                // If we reach here, the backend call was successful.
                this.activeConnectionId = connectionId;
                this.activeModel = loadedModelMeta;
                console.log("[ChatStore] State updated with new active model.");
                
            } catch (error) {
                console.error("[ChatStore] ERROR: 'LoadModel' failed.", error);
                loadingMessage.destroy();
                $message.error(`Failed to load model: ${error}`, { duration: 5000 });
                this.activeConnectionId = null;
                this.activeModel = null;
                return;
            } finally {
                this.isLoadingModel = false;
                console.log("[ChatStore] FINALLY: isLoadingModel set to false.");
            }

            // --- DEBUG STEP 4 ---
            console.log("[ChatStore] Destroying the 'loading' message...");
            loadingMessage.destroy();
            console.log("[ChatStore] 'loading' message destroyed.");

            // --- DEBUG STEP 5 ---
            console.log(`[ChatStore] Creating SUCCESS message for model: ${loadedModelMeta.name}`);
            $message.success(`Model "${loadedModelMeta.name}" is now active!`, { duration: 3000 });
            
            // --- DEBUG STEP 6 ---
            // This is the most important log. If you see this, but the UI is frozen,
            // the problem is with the $message.success call itself.
            console.log("[ChatStore] SUCCESS message created. Function is about to exit. UI should be responsive NOW.");
        },

        async sendMessage(prompt) {
            if (!this.activeModel || this.isStreamingResponse) return;
            const userMessage = { sender: 'user', text: prompt.trim() };
            this.messages.push(userMessage);
            this.isStreamingResponse = true;
            try {
                const responseText = await SendMessage(prompt.trim());
                this.messages.push({ sender: 'model', text: responseText });
            } catch (error) {
                this.messages.push({ sender: 'model', text: `Error: ${error}`, isError: true });
            } finally {
                this.isStreamingResponse = false;
            }
        },
    },
});

export default useChatStore;