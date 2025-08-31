import { create } from 'zustand'
import { config } from '../../wailsjs/go/models'

// This interface defines the shape of our global state
interface AppState {
  connections: config.ConnectionMetadata[]
  activeModel: config.ConnectionMetadata | null
  setConnections: (connections: config.ConnectionMetadata[]) => void
  setActiveModel: (model: config.ConnectionMetadata | null) => void
}

// create the store
export const useAppStore = create<AppState>((set) => ({
  connections: [],
  activeModel: null,
  setConnections: (connections) => set({ connections }),
  setActiveModel: (model) => set({ activeModel: model }),
}))