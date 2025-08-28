import { App } from '../../wailsjs/go/app/App';

export const llmService = {
  async startRedTeamScan(): Promise<string> {
    try {
      const result = await App.StartRedTeamScan();
      return result;
    } catch (error) {
      console.error("Error starting red team scan:", error);
      throw new Error(`Failed to start red team scan: ${error}`);
    }
  },

  // Add other LLM-related service methods here as needed
};