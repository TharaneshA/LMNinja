import { defineStore } from 'pinia';
import { LoadModel, SendMessage } from 'wailsjs/go/app/App';

const useChatStore = defineStore('chat', {
    state: () => ({
        activeConnectionId: null,
        activeModel: null, // Holds the full metadata of the loaded model
        messages: [],
        isLoadingModel: false,
        isStreamingResponse: false,
    }),
    getters: {
        isModelLoaded: (state) => !!state.activeModel,
    },
    actions: {
        async loadAndSetActiveModel(connectionId) {
            if (!connectionId) {
                this.activeConnectionId = null;
                this.activeModel = null;
                this.messages = [];
                return;
            }

            if (this.isLoadingModel || this.activeConnectionId === connectionId) return;

            this.isLoadingModel = true;
            this.messages = [];
            const loadingKey = "loading-model";
            $message.loading("Loading model...", { duration: 0, key: loadingKey });

            try {
                const loadedModelMeta = await LoadModel(connectionId);
                this.activeConnectionId = connectionId;
                this.activeModel = loadedModelMeta;
                $message.destroy(loadingKey);
                $message.success(`Model "${loadedModelMeta.name}" is now active!`);
            } catch (error) {
                $message.destroy(loadingKey);
                $message.error(`Failed to load model: ${error}`);
                this.activeConnectionId = null;
                this.activeModel = null;
            } finally {
                this.isLoadingModel = false;
            }
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