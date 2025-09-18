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
            if (!connectionId) {
                this.activeConnectionId = null;
                this.activeModel = null;
                this.messages = [];
                return;
            }
            if (this.isLoadingModel) return;

            this.isLoadingModel = true;
            this.messages = [];
            // THIS IS THE CORRECT WAY TO HANDLE LOADING MESSAGES
            const loadingMessage = $message.loading("Loading model...", { duration: 0 });

            try {
                const loadedModelMeta = await LoadModel(connectionId);
                this.activeConnectionId = connectionId;
                this.activeModel = loadedModelMeta;
                
                // Update the message to success and then remove it
                loadingMessage.type = 'success';
                loadingMessage.content = `Model "${loadedModelMeta.name}" is now active!`;
                setTimeout(() => loadingMessage.destroy(), 2000);

            } catch (error) {
                // Update the message to error and keep it
                loadingMessage.type = 'error';
                loadingMessage.content = `Failed to load model: ${error}`;
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