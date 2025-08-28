"use client"

import { useState } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Plus, Edit, Trash2, CheckCircle, XCircle, Cloud, Server, TestTube } from "lucide-react"

const existingConnections = [
  { id: 1, name: "GPT-4 Production", type: "OpenAI API", status: "connected" },
  { id: 2, name: "Claude-3 Testing", type: "Anthropic API", status: "connected" },
  { id: 3, name: "Local Llama", type: "Ollama", status: "error" },
  { id: 4, name: "Mistral Local", type: "Ollama", status: "connected" },
]

export function Connectors() {
  const [selectedConnection, setSelectedConnection] = useState<any>(null)
  const [connectionType, setConnectionType] = useState("cloud")
  const [formData, setFormData] = useState({
    name: "",
    provider: "",
    apiKey: "",
    serverAddress: "",
    modelPath: "",
  })

  const handleConnectionSelect = (connection: any) => {
    setSelectedConnection(connection)
    setFormData({
      name: connection.name,
      provider: connection.type,
      apiKey: "",
      serverAddress: "",
      modelPath: "",
    })
  }

  const handleTestConnection = () => {
    // Simulate connection test
    console.log("Testing connection...")
  }

  const handleSaveConnection = () => {
    // Simulate saving connection
    console.log("Saving connection...", formData)
  }

  return (
    <div className="p-6">
      {/* Header */}
      <div className="mb-6">
        <h1 className="text-3xl font-bold">LLM Connectors</h1>
        <p className="text-gray-400 mt-1">Manage your language model connections</p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Connections List */}
        <Card className="bg-[#1B2636] border-[#2A3B52] lg:col-span-1">
          <CardHeader>
            <CardTitle className="flex items-center justify-between">
              <span>Connections</span>
              <Button size="sm" className="bg-[#4A90E2] hover:bg-[#3A7BC8]">
                <Plus className="w-4 h-4 mr-1" />
                New
              </Button>
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {existingConnections.map((connection) => (
                <div
                  key={connection.id}
                  className={`p-3 rounded-lg cursor-pointer transition-colors ${
                    selectedConnection?.id === connection.id
                      ? "bg-[#4A90E2] bg-opacity-20 border border-[#4A90E2]"
                      : "bg-[#0F1419] hover:bg-[#2A3B52]"
                  }`}
                  onClick={() => handleConnectionSelect(connection)}
                >
                  <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-3">
                      {connection.status === "connected" ? (
                        <CheckCircle className="w-4 h-4 text-green-500" />
                      ) : (
                        <XCircle className="w-4 h-4 text-red-500" />
                      )}
                      <div>
                        <p className="text-sm font-medium">{connection.name}</p>
                        <p className="text-xs text-gray-400">{connection.type}</p>
                      </div>
                    </div>
                    <div className="flex items-center space-x-1">
                      <Button size="sm" variant="ghost">
                        <Edit className="w-3 h-3" />
                      </Button>
                      <Button size="sm" variant="ghost">
                        <Trash2 className="w-3 h-3" />
                      </Button>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Connection Form */}
        <Card className="bg-[#1B2636] border-[#2A3B52] lg:col-span-2">
          <CardHeader>
            <CardTitle>{selectedConnection ? "Edit Connection" : "Add New Connection"}</CardTitle>
          </CardHeader>
          <CardContent>
            <Tabs value={connectionType} onValueChange={setConnectionType}>
              <TabsList className="grid w-full grid-cols-2 mb-6">
                <TabsTrigger value="cloud" className="flex items-center space-x-2">
                  <Cloud className="w-4 h-4" />
                  <span>Cloud API</span>
                </TabsTrigger>
                <TabsTrigger value="local" className="flex items-center space-x-2">
                  <Server className="w-4 h-4" />
                  <span>Local Model</span>
                </TabsTrigger>
              </TabsList>

              <TabsContent value="cloud" className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="name">Connection Name</Label>
                  <Input
                    id="name"
                    placeholder="e.g., GPT-4 Production"
                    value={formData.name}
                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                    className="bg-[#0F1419] border-[#2A3B52]"
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="provider">API Provider</Label>
                  <Select
                    value={formData.provider}
                    onValueChange={(value) => setFormData({ ...formData, provider: value })}
                  >
                    <SelectTrigger className="bg-[#0F1419] border-[#2A3B52]">
                      <SelectValue placeholder="Select provider" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="openai">OpenAI</SelectItem>
                      <SelectItem value="anthropic">Anthropic</SelectItem>
                      <SelectItem value="google">Google</SelectItem>
                      <SelectItem value="cohere">Cohere</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="apiKey">API Key</Label>
                  <Input
                    id="apiKey"
                    type="password"
                    placeholder="Enter your API key"
                    value={formData.apiKey}
                    onChange={(e) => setFormData({ ...formData, apiKey: e.target.value })}
                    className="bg-[#0F1419] border-[#2A3B52]"
                  />
                </div>
              </TabsContent>

              <TabsContent value="local" className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="localName">Connection Name</Label>
                  <Input
                    id="localName"
                    placeholder="e.g., Local Llama"
                    value={formData.name}
                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                    className="bg-[#0F1419] border-[#2A3B52]"
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="modelType">Model Type</Label>
                  <Select>
                    <SelectTrigger className="bg-[#0F1419] border-[#2A3B52]">
                      <SelectValue placeholder="Select model type" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="ollama">Ollama</SelectItem>
                      <SelectItem value="llamacpp">Llama.cpp</SelectItem>
                      <SelectItem value="textgen">Text Generation WebUI</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="serverAddress">Server Address</Label>
                  <Input
                    id="serverAddress"
                    placeholder="http://localhost:11434"
                    value={formData.serverAddress}
                    onChange={(e) => setFormData({ ...formData, serverAddress: e.target.value })}
                    className="bg-[#0F1419] border-[#2A3B52]"
                  />
                </div>
              </TabsContent>
            </Tabs>

            <div className="flex space-x-3 mt-6">
              <Button
                variant="outline"
                onClick={handleTestConnection}
                className="flex items-center space-x-2 bg-transparent"
              >
                <TestTube className="w-4 h-4" />
                <span>Test Connection</span>
              </Button>
              <Button onClick={handleSaveConnection} className="bg-[#4A90E2] hover:bg-[#3A7BC8]">
                Save Connection
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
