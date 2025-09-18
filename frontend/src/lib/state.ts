import { create } from 'zustand'
import { storage } from '@wailsjs/go/models'

export interface AppState {
  connections: storage.ConnectionMetadata[]
  activeModel: storage.ConnectionMetadata | null
  setConnections: (connections: storage.ConnectionMetadata[]) => void
  setActiveModel: (model: storage.ConnectionMetadata | null) => void
}

export const useAppStore = create<AppState>((set) => ({
  connections: [],
  activeModel: null,
  setConnections: (connections) => set({ connections }),
  setActiveModel: (model) => set({ activeModel: model }),
}))