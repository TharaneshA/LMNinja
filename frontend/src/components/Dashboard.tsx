"use client";
import { useAppStore } from "@/lib/state";
import { LoadModel } from "@wailsjs/go/app/App";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { toast } from "sonner";
import RedTeamChat from "./RedTeamChat";
import { Bot } from "lucide-react";

export default function Dashboard() {
  const { connections, activeModel, setActiveModel } = useAppStore();
  const handleModelSelect = async (connectionId: string) => {
    if (!connectionId) { setActiveModel(null); return; }
    try {
      toast.info("Loading model...");
      const loadedModelMeta = await LoadModel(connectionId);
      setActiveModel(loadedModelMeta);
      toast.success(`Model "${loadedModelMeta.name}" is now active!`);
    } catch (error) {
      toast.error("Failed to load model.", { description: String(error) });
      setActiveModel(null);
    }
  };
  return (
    <div className="p-6 space-y-4 flex flex-col h-full">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Red Teaming</h1>
          <p className="text-gray-400 mt-1">{activeModel ? `Interactive session with ${activeModel.name}` : "Select a model to begin"}</p>
        </div>
        <div className="flex items-center gap-3">
          <span className="text-sm text-gray-400">Active Model:</span>
          <Select value={activeModel?.id || ""} onValueChange={handleModelSelect}>
            <SelectTrigger className="w-[280px] bg-[#1B2636] border-[#2A3B52]"><SelectValue placeholder="Select a model to load..." /></SelectTrigger>
            <SelectContent>
              {connections.length > 0 ? (
                connections.map((conn) => (<SelectItem key={conn.id} value={conn.id}>{conn.name}</SelectItem>))
              ) : (<div className="p-2 text-sm text-gray-400 text-center">No connections configured.</div>)}
            </SelectContent>
          </Select>
        </div>
      </div>
      <div className="flex-1 flex flex-col border-t border-[#2A3B52] pt-4">
        {!activeModel ? (
          <div className="flex-1 flex items-center justify-center">
            <div className="text-center p-8 bg-[#1B2636] rounded-lg border border-[#2A3B52]">
              <Bot className="w-12 h-12 mx-auto text-gray-400 mb-4" />
              <h2 className="text-xl font-semibold">No Model Loaded</h2>
              <p className="text-gray-400 mt-2">Go to LLM Connectors to add a model, then select it here.</p>
            </div>
          </div>
        ) : (<RedTeamChat />)}
      </div>
    </div>
  );
}