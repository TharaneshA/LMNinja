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
                const sanitizedContent = log.content.replace(/</g, "&lt;").replace(/>/g, "&gt;");
                const contentClass = `log-content log-content--${log.verdict || 'default'}`;
                return `<div class="log-line log-line--${type}"><span class="log-tag">[${type}]</span><span class="${contentClass}">${sanitizedContent}</span></div>`;
            }).join('');
        }
    },
    actions: {
        addLog(logEntry) {
            this.scanLog.push(logEntry);
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
            this.totalSteps = 0;
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
                    

                    const verdictText = `Verdict: ${evaluation.verdict}`;
                    
                    this.addLog({
                        type: 'EVAL',
                        content: verdictText, 
                        verdict: evaluation.verdict 
                    });

                    const distractionScore = (evaluation.explainability?.distraction_score * 100).toFixed(2);
                    const scoreLabel = distractionScore > 50 ? 'High' : (distractionScore > 20 ? 'Medium' : 'Low');
                    this.addLog({
                        type: 'EXPLAIN',
                        content: `Distraction Score: ${distractionScore}% (${scoreLabel})`
                    });
                    
                    this.completedSteps++;

                    await new Promise(resolve => setTimeout(resolve, 0));
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