import { defineStore } from 'pinia';
import { NTag } from 'naive-ui';
import { h } from 'vue';
import { GetPromptsForScan, LoadModel, SendMessage, EvaluatePrompt, UnloadModel } from 'wailsjs/go/app/App';

const getTagType = (logType) => {
    switch(logType.toUpperCase()) {
        case 'SENT': return 'warning';
        case 'RECV': return 'info';
        case 'EVAL': return 'error';
        case 'INFO': return 'default';
        case 'WARN': return 'warning';
        case 'SUCCESS': return 'success';
        case 'SAFE': return 'success';
        case 'UNSAFE': return 'error';
        default: return 'default';
    }
};

export const useScanStore = defineStore('scan', {
    state: () => ({
        targetModelId: null,
        selectedCategories: [],
        scanLimit: 10,
        isScanning: false,
        scanLog: [],
        totalSteps: 0,
        completedSteps: 0,
    }),
    getters: {
        progress: (state) => {
            if (state.totalSteps === 0) return 0;
            return (state.completedSteps / state.totalSteps) * 100;
        },
        progressText: (state) => {
            return `Step ${state.completedSteps} of ${state.totalSteps}`;
        },
        formattedLog: (state) => {
            const lines = state.scanLog.map(log => {
                return h('div', { style: 'line-height: 1.5;' }, [
                    h(NTag, { 
                        type: getTagType(log.type), 
                        size: 'small', 
                        bordered: false 
                    }, () => log.type.toUpperCase()),
                    h('span', { style: 'margin-left: 10px; white-space: pre-wrap; word-break: break-all;' }, log.content)
                ]);
            });
            return () => h('div', { style: 'display: flex; flex-direction: column; gap: 4px;' }, lines);
        }
    },
    actions: {
        addLog(type, content) {
            this.scanLog.push({ type, content });
            if (this.scanLog.length > 500) this.scanLog.shift();
        },
        async startScan() {
            if (!this.targetModelId || this.selectedCategories.length === 0) {
                $message.error("Please select a target model and at least one attack category.");
                return;
            }
            this.isScanning = true;
            this.scanLog = [];
            this.completedSteps = 0;
            this.addLog('INFO', 'Scan initiated...');
            try {
                this.addLog('INFO', `Fetching up to ${this.scanLimit} prompts...`);
                const promptsToRun = await GetPromptsForScan(this.selectedCategories, this.scanLimit);
                this.totalSteps = promptsToRun.length;
                if (this.totalSteps === 0) {
                    this.addLog('WARN', 'No prompts found. Stopping scan.');
                    this.isScanning = false;
                    return;
                }
                this.addLog('SUCCESS', `Found ${this.totalSteps} prompts.`);
                this.addLog('INFO', `Loading target model...`);
                await LoadModel(this.targetModelId);
                this.addLog('SUCCESS', 'Target model loaded.');
                for (const prompt of promptsToRun) {
                    if (!this.isScanning) {
                        this.addLog('WARN', 'Scan stopped by user.');
                        break;
                    }
                    this.addLog('SENT', prompt);
                    const response = await SendMessage(prompt);
                    this.addLog('RECV', response);
                    const evalResultString = await EvaluatePrompt(prompt);
                    const evaluation = JSON.parse(evalResultString);
                    this.addLog(
                        evaluation.verdict,
                        `Injection Score: ${(evaluation.details.prompt_injection.score * 100).toFixed(2)}%`
                    );
                    this.completedSteps++;
                }
            } catch (error) {
                this.addLog('ERROR', `Scan failed: ${error}`);
            } finally {
                this.addLog('INFO', 'Scan finished. Unloading model...');
                if (this.targetModelId) await UnloadModel(this.targetModelId);
                this.addLog('SUCCESS', 'Cleanup complete.');
                this.isScanning = false;
            }
        },
        stopScan() {
            this.isScanning = false;
        }
    }
});

export default useScanStore;