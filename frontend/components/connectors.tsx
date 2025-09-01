"use client";

import { useState } from "react";
import { useAppStore } from '@/lib/state';
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Plus, Trash2, CheckCircle, Cloud, Server, Loader2, Wand2, FolderSearch } from "lucide-react";
import { toast } from "sonner";
import { SaveConnection, DeleteConnection, GetProviderModels, SelectGGUFFolder, TestConnection } from "@wailsjs/go/app/App";
import { storage, app } from "@wailsjs/go/models";

type GGUFFile = app.GGUFFile;
type FormData = { id: string; name: string; provider: string; model: string; apiKey: string; };
const emptyForm: FormData = { id: "", name: "", provider: "openai", model: "", apiKey: "" };
const basename = (path: string) => path.split(/[\\/]/).pop() || "";

export function Connectors() {
  const { connections, setConnections, activeModel, setActiveModel } = useAppStore();
  
  const [selectedConnectionId, setSelectedConnectionId] = useState<string | null>(null);
  const [formData, setFormData] = useState<FormData>(emptyForm);
  const [isSaving, setIsSaving] = useState(false);
  const [isTesting, setIsTesting] = useState(false);
  const [isFetchingModels, setIsFetchingModels] = useState(false);
  const [availableModels, setAvailableModels] = useState<string[]>([]);
  const [discoveredGGUFs, setDiscoveredGGUFs] = useState<GGUFFile[]>([]);
  const [isScanning, setIsScanning] = useState(false);
  
  const resetForm = (providerType: 'cloud' | 'local' = 'cloud') => {
    setSelectedConnectionId(null);
    const defaultProvider = providerType === 'cloud' ? 'openai' : 'ollama';
    setFormData({ ...emptyForm, provider: defaultProvider });
    setAvailableModels([]);
    setDiscoveredGGUFs([]);
  };

  const handleConnectionSelect = (connection: storage.ConnectionMetadata) => {
    setSelectedConnectionId(connection.id);
    setFormData({
      id: connection.id, name: connection.name, provider: connection.provider, model: connection.model, apiKey: "",
    });
    if (["openai", "gemini", "anthropic"].includes(connection.provider)) {
      setAvailableModels([connection.model]); setDiscoveredGGUFs([]);
    } else if (connection.provider === "gguf") {
      setDiscoveredGGUFs([{ name: basename(connection.model), path: connection.model }]); setAvailableModels([]);
    } else {
      setAvailableModels([]); setDiscoveredGGUFs([]);
    }
  };

  const handleFormChange = (field: keyof FormData, value: string) => {
    const newFormData = { ...formData, [field]: value };
    if (field === 'provider') {
        newFormData.model = ''; setAvailableModels([]); setDiscoveredGGUFs([]);
    }
    setFormData(newFormData);
  };

  const handleFetchModels = async () => {
    if (!formData.provider || !formData.apiKey) {
      toast.warning("Please select a provider and enter an API key first."); return;
    }
    setIsFetchingModels(true); setAvailableModels([]);
    try {
      const models = await GetProviderModels(formData.provider, formData.apiKey);
      setAvailableModels(models || []);
      if (!models || models.length === 0) {
        toast.error("No models found", { description: "Your API key may be invalid or have no permissions." });
      } else { toast.success(`Found ${models.length} compatible models.`); }
    } catch (error) { toast.error("Failed to fetch models", { description: String(error) });
    } finally { setIsFetchingModels(false); }
  };

  const handleSelectGGUFFolder = async () => {
    setIsScanning(true); toast.info("Opening folder selection dialog...");
    try {
      const files = await SelectGGUFFolder();
      setDiscoveredGGUFs(files);
      if (files.length === 0) { toast.info("No GGUF files found in the selected folder.");
      } else { toast.success(`Discovered ${files.length} GGUF model(s).`); }
    } catch (error) {
      if (!String(error).includes("No directory selected")) {
          toast.error("Failed to scan folder", { description: String(error) });
      }
    } finally { setIsScanning(false); }
  };

  const handleTestOllama = async () => {
    if (!formData.model) {
        toast.warning("Please enter an Ollama model name to test (e.g., llama3)."); return;
    }
    setIsTesting(true);
    try {
        const meta = new storage.ConnectionMetadata({ provider: 'ollama', model: formData.model });
        const result = await TestConnection(meta, "");
        toast.success(result, { description: `Successfully connected to Ollama and got a response from '${formData.model}'.` });
    } catch (error) { toast.error("Ollama Connection Test Failed", { description: String(error) });
    } finally { setIsTesting(false); }
  };

  const handleSaveConnection = async () => {
    if (!formData.name || !formData.provider || !formData.model) {
      toast.warning("Please provide a Connection Name and select a Model."); return;
    }
    if (formData.provider !== 'gguf' && formData.provider !== 'ollama' && !formData.id && !formData.apiKey) {
      toast.warning("API Key is required for new Cloud connections."); return;
    }
    setIsSaving(true);
    try {
      const meta = new storage.ConnectionMetadata({ id: formData.id, name: formData.name, provider: formData.provider, model: formData.model });
      const updatedConnections = await SaveConnection(meta, formData.apiKey);
      setConnections(updatedConnections);
      toast.success(`Connection "${formData.name}" saved successfully!`);
      const savedConn = updatedConnections.find(c => c.id === meta.id || c.name === meta.name);
      if (savedConn) handleConnectionSelect(savedConn);
    } catch (error) { toast.error("Failed to save connection", { description: String(error) });
    } finally { setIsSaving(false); }
  };
  
  const handleDeleteConnection = async (id: string, name: string) => {
    if (!window.confirm(`Are you sure you want to delete the connection "${name}"?`)) return;
    try {
        if (activeModel?.id === id) setActiveModel(null);
        const updatedConnections = await DeleteConnection(id);
        setConnections(updatedConnections);
        toast.success(`Connection "${name}" deleted.`);
        resetForm();
    } catch (error) { toast.error("Failed to delete connection", { description: String(error) }); }
  };

  return (
    <div className="p-6">
      <div className="mb-6"><h1 className="text-3xl font-bold">LLM Connectors</h1><p className="text-gray-400 mt-1">Manage your language model connections</p></div>
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <Card className="bg-[#1B2636] border-[#2A3B52] lg:col-span-1">
          <CardHeader><CardTitle className="flex items-center justify-between"><span>My Models</span><Button size="sm" className="bg-[#4A90E2] hover:bg-[#3A7BC8]" onClick={() => resetForm()}><Plus className="w-4 h-4 mr-1" /> New</Button></CardTitle></CardHeader>
          <CardContent>{connections.length === 0 ? (<div className="text-center text-gray-400 py-8 px-4"><p className="font-semibold">No models connected.</p><p className="text-sm mt-1">Click 'New' to add your first connection.</p></div>) : (<div className="space-y-3">{connections.map((c) => (<div key={c.id} className={`p-3 rounded-lg cursor-pointer transition-colors ${selectedConnectionId === c.id ? "bg-[#4A90E2] bg-opacity-20 border border-[#4A90E2]" : "bg-[#0F1419] hover:bg-[#2A3B52]"}`} onClick={() => handleConnectionSelect(c)}><div className="flex items-center justify-between"><div className="flex items-center space-x-3"><CheckCircle className="w-4 h-4 text-green-500" /><div><p className="text-sm font-medium">{c.name}</p><p className="text-xs text-gray-400 uppercase">{c.provider}: {basename(c.model)}</p></div></div><Button size="icon" variant="ghost" className="h-7 w-7 hover:bg-red-500/20 text-gray-400 hover:text-red-400" onClick={(e) => { e.stopPropagation(); handleDeleteConnection(c.id, c.name)}}><Trash2 className="w-3.5 h-3.5" /></Button></div></div>))}</div>)}</CardContent>
        </Card>
        <Card className="bg-[#1B2636] border-[#2A3B52] lg:col-span-2">
          <CardHeader><CardTitle>{selectedConnectionId ? `Edit "${formData.name}"` : "Add New Connection"}</CardTitle></CardHeader>
          <CardContent>
            <Tabs value={['ollama', 'gguf'].includes(formData.provider) ? 'local' : 'cloud'} onValueChange={(val) => resetForm(val as 'cloud' | 'local')}>
              <TabsList className="grid w-full grid-cols-2 mb-6"><TabsTrigger value="cloud" className="flex items-center space-x-2"><Cloud className="w-4 h-4" /> Cloud API</TabsTrigger><TabsTrigger value="local" className="flex items-center space-x-2"><Server className="w-4 h-4" /> Local Model</TabsTrigger></TabsList>
              
              <TabsContent value="cloud" className="space-y-4">
                <div className="space-y-2"><Label htmlFor="name">Connection Name</Label><Input id="name" placeholder="e.g., My Personal OpenAI Key" value={formData.name} onChange={(e) => handleFormChange('name', e.target.value)} className="bg-[#0F1419] border-[#2A3B52]" /></div>
                <div className="space-y-2"><Label htmlFor="provider">API Provider</Label><Select value={formData.provider} onValueChange={(value) => handleFormChange('provider', value)}><SelectTrigger className="w-full bg-[#0F1419] border-[#2A3B52]"><SelectValue placeholder="Select a provider" /></SelectTrigger><SelectContent><SelectItem value="openai">OpenAI</SelectItem><SelectItem value="anthropic">Anthropic</SelectItem><SelectItem value="gemini">Google Gemini</SelectItem></SelectContent></Select></div>
                <div className="space-y-2"><Label htmlFor="apiKey">API Key</Label><div className="flex items-center gap-2"><Input id="apiKey" type="password" placeholder={selectedConnectionId ? "Leave blank to keep existing key" : "Enter your API key"} value={formData.apiKey} onChange={(e) => handleFormChange('apiKey', e.target.value)} className="bg-[#0F1419] border-[#2A3B52]" /><Button variant="outline" onClick={handleFetchModels} disabled={isFetchingModels || !formData.provider || !formData.apiKey}>{isFetchingModels ? <Loader2 className="w-4 h-4 animate-spin" /> : <Wand2 className="w-4 h-4" />}<span className="ml-2 hidden sm:inline">Fetch Models</span></Button></div></div>
                { (availableModels.length > 0 || isFetchingModels) && (<div className="space-y-2"><Label htmlFor="model">Model</Label><Select value={formData.model} onValueChange={(value) => handleFormChange('model', value)} disabled={isFetchingModels}><SelectTrigger className="w-full bg-[#0F1419] border-[#2A3B52]"><SelectValue placeholder={isFetchingModels ? "Fetching models..." : "Select a model"} /></SelectTrigger><SelectContent>{availableModels.map(modelId => (<SelectItem key={modelId} value={modelId}>{modelId}</SelectItem>))}</SelectContent></Select></div>)}
              </TabsContent>

              <TabsContent value="local" className="space-y-4">
                <div className="space-y-2"><Label>Local Provider</Label><Select value={formData.provider} onValueChange={(value) => handleFormChange('provider', value)}><SelectTrigger className="w-full bg-[#0F1419] border-[#2A3B52]"><SelectValue /></SelectTrigger><SelectContent><SelectItem value="ollama">Ollama</SelectItem><SelectItem value="gguf">GGUF</SelectItem></SelectContent></Select></div>
                
                {formData.provider === 'ollama' && (<>
                  <div className="space-y-2"><Label htmlFor="ollama_name">Connection Name</Label><Input id="ollama_name" placeholder="e.g., My Local Llama3" value={formData.name} onChange={(e) => handleFormChange('name', e.target.value)} className="bg-[#0F1419] border-[#2A3B52]" /></div>
                  <div className="space-y-2"><Label htmlFor="ollama_model">Ollama Model Name</Label><Input id="ollama_model" placeholder="e.g., llama3 (must be pulled in Ollama)" value={formData.model} onChange={(e) => handleFormChange('model', e.target.value)} className="bg-[#0F1419] border-[#2A3B52]" /><p className="text-xs text-gray-400">This must match a model name from `ollama list`.</p></div>
                  <Button variant="outline" onClick={handleTestOllama} disabled={isTesting || !formData.model}>{isTesting ? <Loader2 className="w-4 h-4 animate-spin" /> : <Wand2 className="w-4 h-4" />}<span className="ml-2">Test Ollama</span></Button>
                </>)}

                {formData.provider === 'gguf' && (<>
                    <Button variant="outline" className="w-full" onClick={handleSelectGGUFFolder} disabled={isScanning}>{isScanning ? <Loader2 className="w-4 h-4 animate-spin" /> : <FolderSearch className="w-4 h-4" />}<span className="ml-2">Scan Folder for GGUF Models</span></Button>
                    {(discoveredGGUFs.length > 0) && (<>
                        <div className="space-y-2"><Label htmlFor="gguf_model">Discovered Model File</Label><Select value={formData.model} onValueChange={(value) => handleFormChange('model', value)}><SelectTrigger className="w-full bg-[#0F1419] border-[#2A3B52]"><SelectValue placeholder="Select a discovered GGUF model" /></SelectTrigger><SelectContent>{discoveredGGUFs.map(f => (<SelectItem key={f.path} value={f.path}>{f.name}</SelectItem>))}</SelectContent></Select></div>
                        <div className="space-y-2"><Label htmlFor="gguf_name">Connection Name</Label><Input id="gguf_name" placeholder="e.g., Local Phi-3 Mini" value={formData.name} onChange={(e) => handleFormChange('name', e.target.value)} className="bg-[#0F1419] border-[#2A3B52]" /></div>
                    </>)}
                </>)}
              </TabsContent>
            </Tabs>
            <div className="flex space-x-3 mt-6">
                <Button onClick={handleSaveConnection} disabled={isSaving || !formData.model} className="bg-[#4A90E2] hover:bg-[#3A7BC8]">{isSaving && <Loader2 className="w-4 h-4 animate-spin mr-2" />} {isSaving ? "Saving..." : "Save Connection"}</Button>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}