import { defineStore } from 'pinia';
import { GetDashboardStats, GetVulnerabilitiesByModel, GetScanHistory } from 'wailsjs/go/app/App';

export const useDashboardStore = defineStore('dashboard', {
    state: () => ({
        loading: false,
        stats: {
            modelsScanned: 0,
            totalScans: 0,
            vulnerabilitiesFound: 0,
            overallPassRate: 0.0,
        },
        vulnerabilitiesByModel: {
            labels: [],
            datasets: [{
                label: 'Vulnerabilities',
                backgroundColor: '#e88080',
                data: [],
            }],
        },
        scanHistory: [],
    }),
    actions: {
        async fetchDashboardData() {
            if (this.loading) return;
            this.loading = true;
            try {
                // Fetch all data in parallel for speed
                const [stats, vulnsByModel, history] = await Promise.all([
                    GetDashboardStats(),
                    GetVulnerabilitiesByModel(),
                    GetScanHistory(),
                ]);

                // 1. Update stats
                this.stats = stats;

                // 2. Process and update chart data
                const labels = vulnsByModel.map(item => item.modelName);
                const data = vulnsByModel.map(item => item.vulnCount);
                this.vulnerabilitiesByModel = {
                    labels,
                    datasets: [{
                        label: 'Vulnerabilities Found',
                        backgroundColor: '#e88080',
                        data,
                    }],
                };

                // 3. Update scan history table data
                this.scanHistory = history.map(item => ({
                    ...item,
                    passRate: item.totalPrompts > 0 
                        ? `${(( (item.totalPrompts - item.vulnerabilitiesFound) / item.totalPrompts) * 100).toFixed(1)}%`
                        : 'N/A',
                }));

            } catch (error) {
                $message.error(`Failed to load dashboard data: ${error}`);
            } finally {
                this.loading = false;
            }
        }
    }
});

export default useDashboardStore;