<script setup>
import { ref, watch, nextTick, computed } from 'vue';
import { useThemeVars } from 'naive-ui';
import useChatStore from 'stores/chat.js';
import MarkdownIt from 'markdown-it';

const md = new MarkdownIt({
  html: false, // Security: Disable raw HTML tags in markdown
  linkify: true, // Automatically convert URL-like text to links
  typographer: true, // Enable smart quotes and other nice typographic features
});

const chatStore = useChatStore();
const prompt = ref('');
const scrollArea = ref(null);
const themeVars = useThemeVars();

const isChatDisabled = computed(() => {
  return !chatStore.isModelLoaded || chatStore.isStreamingResponse;
});

const isSendDisabled = computed(() => {
  return isChatDisabled.value || prompt.value.trim() === '';
});

const renderMarkdown = (text) => {
  return md.render(text);
};

const scrollToBottom = () => {
  nextTick(() => {
    if (scrollArea.value) {
      scrollArea.value.scrollTop = scrollArea.value.scrollHeight;
    }
  });
};

watch(() => chatStore.messages, scrollToBottom, { deep: true });

const handleSendMessage = () => {
    if (isSendDisabled.value) return;
    chatStore.sendMessage(prompt.value);
    prompt.value = '';
};
</script>

<template>
    <div class="chat-container">
        <!-- Message History Area -->
        <div class="message-history" ref="scrollArea">
            <div v-if="!chatStore.isModelLoaded" class="welcome-message">
                <n-empty description="No model is currently loaded. Select one from the top bar to begin your session." />
            </div>
            <div v-else-if="chatStore.messages.length === 0 && !chatStore.isStreamingResponse" class="welcome-message">
                <n-empty :description="`Ready to chat with ${chatStore.activeModel?.name}.`" />
            </div>

            <!-- Conversation -->
            <div v-for="(msg, index) in chatStore.messages" :key="index" :class="['message', `message-${msg.sender}`]">
                <div class="message-bubble" :class="{ 'error-bubble': msg.isError }">

                    <div v-if="msg.sender === 'model'" 
                         v-html="renderMarkdown(msg.text)" 
                         class="markdown-body">
                    </div>

                    <p v-else>{{ msg.text }}</p>
                </div>
            </div>
            
            <div v-if="chatStore.isStreamingResponse" class="message message-model">
                <div class="message-bubble">
                    <n-spin size="small" />
                </div>
            </div>
        </div>

        <div class="input-area">
            <n-input
                v-model:value="prompt"
                type="textarea"
                :autosize="{ minRows: 1, maxRows: 5 }"
                placeholder="Send a message... (Enter to send, Shift+Enter for newline)"
                @keypress.enter.prevent.exact="handleSendMessage"
                :disabled="isChatDisabled"
            />
            <n-button 
                @click="handleSendMessage" 
                type="primary" 
                :loading="chatStore.isStreamingResponse"
                :disabled="isSendDisabled"
            >
                Send
            </n-button>
        </div>
    </div>
</template>

<style lang="scss" scoped>
@import 'github-markdown-css/github-markdown-dark.css';

.chat-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: v-bind('themeVars.bodyColor');
}
.message-history {
  flex-grow: 1;
  overflow-y: auto;
  padding: 16px;
}
.message {
  display: flex;
  margin-bottom: 12px;
  max-width: 90%;
}
.message-user { 
  margin-left: auto;
  justify-content: flex-end;
}
.message-model { 
  margin-right: auto;
  justify-content: flex-start;
}
.message-bubble {
  padding: 10px 15px;
  border-radius: 18px;
  line-height: 1.5;
  
  user-select: text;
  cursor: text;

  p { 
    margin: 0; 
    white-space: pre-wrap; 
    word-wrap: break-word;
  }
}
.message-user .message-bubble {
  background-color: v-bind('themeVars.primaryColor');
  color: white;
  border-bottom-right-radius: 4px;
}
.message-model .message-bubble {
  background-color: v-bind('themeVars.cardColor');
  border-bottom-left-radius: 4px;
}
.error-bubble {
  background-color: v-bind('themeVars.errorColor');
  color: white;
}
.input-area {
  display: flex;
  gap: 10px;
  align-items: flex-end;
  padding: 8px;
  border-top: 1px solid v-bind('themeVars.borderColor');
  flex-shrink: 0;
}
.welcome-message {
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    text-align: center;
    color: v-bind('themeVars.textColor3');
}


:deep(.markdown-body) {
    background-color: transparent;
    color: v-bind('themeVars.textColorBase');
    font-size: 14px;

    // Style for code blocks
    pre {
        background-color: v-bind('themeVars.codeColor') !important;
        padding: 12px;
        border-radius: 6px;
    }

    // Style for inline code
    code {
        background-color: v-bind('themeVars.codeColor') !important;
        color: #c9d1d9; // A standard light text color for dark code blocks
        border-radius: 4px;
        padding: 2px 4px;
    }

    // Style for blockquotes
    blockquote {
        border-left: 0.25em solid v-bind('themeVars.borderColor');
        color: v-bind('themeVars.textColor3');
    }
}
</style>