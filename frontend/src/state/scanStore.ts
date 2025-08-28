import { create, StateCreator } from 'zustand';

interface ScanState {
  scanResultText: string;
  setScanResultText: (text: string) => void;
}

export const useScanStore = create<ScanState>((set) => ({
  scanResultText: "No scan initiated yet.",
  setScanResultText: (text: string) => set({ scanResultText: text }),
}));