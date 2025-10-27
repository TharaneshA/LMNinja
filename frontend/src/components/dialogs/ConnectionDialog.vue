<script setup>
import { computed, reactive, ref, watch } from 'vue';
import useDialog, { ConnDialogType } from 'stores/dialog';
import useConnectionStore from 'stores/connections.js';
import { GetProviderModels, TestConnection, SelectGGUFFile, SelectHuggingFaceFolder } from 'wailsjs/go/app/App';

const dialogStore = useDialog();
const connectionStore = useConnectionStore();

const formRef = ref(null);
const form = reactive({
    id: '', name: '', group: '', provider: 'openai', model: '', apiKey: '',
});

const isEditMode = computed(() => dialogStore.connType === ConnDialogType.EDIT);
const title = computed(() => isEditMode.value ? "Edit LLM Connector" : "New LLM Connector");

const availableModels = ref([]);
const fetchingModels = ref(false);
const testingConnection = ref(false);

watch(() => dialogStore.connDialogVisible, (visible) => {
    if (visible) {
        Object.assign(form, { id: '', name: '', group: '', provider: 'openai', model: '', apiKey: '' });
        availableModels.value = [];

        if (isEditMode.value && dialogStore.connParam) {
            const param = dialogStore.connParam;
            form.id = param.id;
            form.name = param.name;
            form.group = param.group || '';
            form.provider = param.provider;
            form.model = param.model;
            form.apiKey = '••••••••••••••••';
        }
    }
});

const onSaveConnection = async (e) => {
    e.preventDefault();
    try {
        await formRef.value?.validate();
        const connectionData = {
            id: form.id,
            name: form.name,
            group: form.group,
            provider: form.provider,
            model: form.model,
        };
        const apiKey = form.apiKey !== '••••••••••••••••' ? form.apiKey : '';

        const { success, msg } = await connectionStore.saveConnection(connectionData, apiKey);
        if (success) {
            $message.success("Connection saved successfully!");
            onClose();
        } else {
            $message.error(msg);
        }
    } catch (errors) {
        $message.error("Please fill in all required fields.");
    }
};

const handleFetchModels = async () => {
    if (!form.apiKey || form.apiKey === '••••••••••••••••') {
        $message.warning("Please enter an API Key to fetch models.");
        return;
    }
    fetchingModels.value = true;
    try {
        const models = await GetProviderModels(form.provider, form.apiKey);
        
        if (!models) {
             availableModels.value = [];
             $message.error(`Failed to fetch models: Received null list from backend. (Check Go implementation for provider: ${form.provider})`);
             return;
        }
        
        availableModels.value = models.map(m => ({ label: m, value: m }));
        if (models.length > 0) $message.success(`Found ${models.length} models.`);
        else $message.warning(`No compatible models found.`);
    } catch (error) {
        $message.error(`Failed to fetch models: ${error}`);
    } finally {
        fetchingModels.value = false;
    }
};

const handleSelectHFFolder = async () => {
    try {
        const folder = await SelectHuggingFaceFolder();
        if (folder && folder.path) {
            form.model = folder.path;
            if (!form.name) {
                form.name = folder.name;
            }
        } else {
            $message.info("No folder was selected.");
        }
    } catch (error) {
        $message.error(`Failed to select Hugging Face folder: ${error}`);
    }
};


const handleSelectGGUF = async () => { 
    try {
        const ggufFile = await SelectGGUFFile();
        
        if (ggufFile && ggufFile.path) {
            form.model = ggufFile.path;
            if (!form.name) {
                form.name = ggufFile.name.replace(/\.gguf$/i, '');
            }
        } else {
             $message.info("No .gguf file was selected.");
        }

    } catch(error) {
        $message.error(`Failed to select GGUF file: ${error}`);
    }
};


const handleTestConnection = async () => {
    testingConnection.value = true;
    try {
        const meta = { provider: form.provider, model: form.model };
        const key = form.apiKey !== '••••••••••••••••' ? form.apiKey : '';
        const result = await TestConnection(meta, key);
        $message.success(result);
    } catch (error) {
        $message.error(`Test failed: ${error}`);
    } finally {
        testingConnection.value = false;
    }
};

const onClose = () => dialogStore.closeConnDialog();

const rules = {
    name: { required: true, message: 'Connection name is required', trigger: 'blur' },
    model: { required: true, message: 'Model is required', trigger: 'blur' },
    apiKey: { required: computed(() => ['openai', 'gemini', 'anthropic'].includes(form.provider)), message: 'API Key is required', trigger: 'blur'}
};

const providerOptions = [
    { label: 'OpenAI', value: 'openai' },
    { label: 'Google Gemini', value: 'gemini' },
    { label: 'Anthropic', value: 'anthropic' },
    { label: 'Ollama', value: 'ollama' },
    { label: 'GGUF (File)', value: 'gguf' },
    { label: 'Hugging Face (Folder)', value: 'huggingface' },
];
</script>

<template>
    <n-modal v-model:show="dialogStore.connDialogVisible" preset="dialog" :title="title" :show-icon="false"
             :mask-closable="false" :closable="false" style="width: 600px">
        <n-form ref="formRef" :model="form" :rules="rules" label-placement="top">
            <n-grid :x-gap="10">
                <n-form-item-gi :span="12" label="Connection Name" path="name" required>
                    <n-input v-model:value="form.name" placeholder="e.g., Llama 3 8B Instruct" />
                </n-form-item-gi>
                <n-form-item-gi :span="12" label="Provider" path="provider">
                     <n-select v-model:value="form.provider" :options="providerOptions" />
                </n-form-item-gi>
            </n-grid>

            <!-- Cloud Providers -->
            <template v-if="['openai', 'gemini', 'anthropic'].includes(form.provider)">
                <n-form-item label="API Key" path="apiKey" required>
                    <n-input v-model:value="form.apiKey" type="password" show-password-on="click" placeholder="Enter new key or leave masked" />
                </n-form-item>
                <n-form-item label="Model" path="model" required>
                     <n-input-group>
                         <n-select v-model:value="form.model" :options="availableModels" filterable tag clearable placeholder="e.g., gpt-4o" />
                         <n-button @click="handleFetchModels" :loading="fetchingModels">Fetch Models</n-button>
                     </n-input-group>
                </n-form-item>
            </template>

            <!-- Local Providers -->
            <template v-if="form.provider === 'ollama'">
                <n-form-item label="Model Name" path="model" required>
                    <n-input v-model:value="form.model" placeholder="e.g., llama3 (must be pulled in Ollama)" />
                </n-form-item>
            </template>
            
            <template v-if="form.provider === 'gguf'">
                <n-form-item label="GGUF File Path" path="model" required>
                    <n-input-group>
                         <n-input v-model:value="form.model" placeholder="Select a .gguf file..." readonly />
                         <n-button :focusable="false" @click="handleSelectGGUF">Browse</n-button>
                     </n-input-group>
                </n-form-item>
            </template>
            
            <template v-if="form.provider === 'huggingface'">
                <n-form-item label="Model Folder Path" path="model" required>
                    <n-input-group>
                         <n-input v-model:value="form.model" placeholder="Select a Hugging Face model folder..." readonly />
                         <n-button :focusable="false" @click="handleSelectHFFolder">Browse</n-button>
                     </n-input-group>
                </n-form-item>
            </template>

        </n-form>
         <template #action>
            <n-button @click="handleTestConnection" :loading="testingConnection">Test Connection</n-button>
            <div style="flex-grow: 1"></div>
            <n-button @click="onClose">Cancel</n-button>
            <n-button type="primary" @click="onSaveConnection">Save</n-button>
        </template>
    </n-modal>
</template>