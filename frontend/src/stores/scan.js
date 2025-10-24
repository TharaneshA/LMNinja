import { defineStore } from 'pinia';
import useConnectionStore from 'stores/connections.js';
import {
    GetPromptsForScan,
    LoadModel,
    SendMessage,
    EvaluateCompliance, 
    UnloadModel,
    CreateScanRecord,
    SaveScanResult,
    FinalizeScan,
} from 'wailsjs/go/app/App';

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
            return state.scanLog.map(log => {
                const type = log.type.toUpperCase();
                const contentClass = type === 'EVAL' ? `log-content log-content--${log.verdict}` : 'log-content';
                return `<div class="log-line log-line--${type}"><span class="log-tag">[${type}]</span> <span class="${contentClass}">${log.content}</span></div>`;
            }).join('');
        }
    },
    actions: {
        addLog(logEntry) {
            const sanitizedContent = logEntry.content.replace(/</g, "&lt;").replace(/>/g, "&gt;");
            this.scanLog.push({ ...logEntry, content: sanitizedContent });
            if (this.scanLog.length > 500) {
                this.scanLog.shift();
            }
        },
        async startScan() {
            if (!this.targetModelId || this.selectedCategories.length === 0) {
                $message.error("Please select a target model and at least one attack category.");
                return;
            }
            this.isScanning = true;
            this.scanLog = [];
            this.completedSteps = 0;
            this.addLog({ type: 'INFO', content: 'Scan initiated...' });

            const connectionStore = useConnectionStore();
            const targetModelName = connectionStore.serverProfile[this.targetModelId]?.name || 'Unknown Model';
            let scanId = null;

            try {
                scanId = await CreateScanRecord(targetModelName);
                this.addLog({ type: 'INFO', content: `Scan record created: ${scanId}` });
                this.addLog({ type: 'INFO', content: `Fetching up to ${this.scanLimit} prompts...` });
                const promptsToRun = await GetPromptsForScan(this.selectedCategories, this.scanLimit);
                this.totalSteps = promptsToRun.length;

                if (this.totalSteps === 0) { throw new Error("No attack prompts found."); }
                this.addLog({ type: 'SUCCESS', content: `Found ${this.totalSteps} prompts.` });
                this.addLog({ type: 'INFO', content: `Loading target model: ${targetModelName}...` });
                await LoadModel(this.targetModelId);
                this.addLog({ type: 'SUCCESS', content: 'Target model loaded.' });

                for (const prompt of promptsToRun) {
                    if (!this.isScanning) {
                        this.addLog({ type: 'WARN', content: 'Scan stopped by user.' });
                        await FinalizeScan(scanId, 'CANCELLED');
                        break;
                    }

                    this.addLog({ type: 'SENT', content: prompt });
                    const response = await SendMessage(prompt);
                    this.addLog({ type: 'RECV', content: response });
                    
                    const complianceEvalJson = await EvaluateCompliance(prompt, response);
                    const evaluation = JSON.parse(complianceEvalJson);
                    
                    await SaveScanResult(scanId, prompt, response, complianceEvalJson);

                    this.addLog({
                        type: 'EVAL',
                        content: `Attack Status: ${evaluation.verdict} (${evaluation.reason})`,
                        verdict: evaluation.verdict 
                    });
                    this.completedSteps++;
                }

                if (this.isScanning) {
                    await FinalizeScan(scanId, 'COMPLETED');
                    this.addLog({ type: 'SUCCESS', content: 'Scan complete and results saved.' });
                }

            } catch (error) {
                this.addLog({ type: 'ERROR', content: `A critical error occurred: ${error}` });
                $message.error(`Scan failed: ${error}`);
                if (scanId) { await FinalizeScan(scanId, 'FAILED'); }
            } finally {
                this.addLog({ type: 'INFO', content: 'Unloading target model...' });
                if (this.targetModelId) { await UnloadModel(this.targetModelId); }
                this.addLog({ type: 'INFO', content: 'Cleanup complete.' });
                this.isScanning = false;
            }
        },
        stopScan() {
            this.isScanning = false;
        }
    }
});

export default useScanStore;