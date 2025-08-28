"use client"

import { useState } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Badge } from "@/components/ui/badge"
import { Code, Activity, Cpu, Search, Play, Square, RefreshCw } from "lucide-react"

const mockLogs = [
  {
    timestamp: "2024-01-15 14:32:15",
    level: "INFO",
    message: "[LM NINJA SERVER] Success! HTTP server listening on port 1234",
  },
  { timestamp: "2024-01-15 14:32:15", level: "INFO", message: "[LM NINJA SERVER] Supported endpoints:" },
  {
    timestamp: "2024-01-15 14:32:15",
    level: "INFO",
    message: "[LM NINJA SERVER] -> GET http://localhost:1234/v1/models",
  },
  {
    timestamp: "2024-01-15 14:32:15",
    level: "INFO",
    message: "[LM NINJA SERVER] -> POST http://localhost:1234/v1/chat/completions",
  },
  {
    timestamp: "2024-01-15 14:32:15",
    level: "INFO",
    message: "[LM NINJA SERVER] -> POST http://localhost:1234/v1/completions",
  },
  {
    timestamp: "2024-01-15 14:32:15",
    level: "INFO",
    message: "[LM NINJA SERVER] -> POST http://localhost:1234/v1/embeddings",
  },
  {
    timestamp: "2024-01-15 14:32:15",
    level: "INFO",
    message: "[LM NINJA SERVER] Logs are saved into /Users/user/.cache/lm-ninja/server-logs",
  },
  { timestamp: "2024-01-15 14:32:15", level: "INFO", message: "Server started." },
  {
    timestamp: "2024-01-15 14:32:20",
    level: "DEBUG",
    message: "llama_kv_cache_init: Metal KV buffer size = 80.00 MiB",
  },
  {
    timestamp: "2024-01-15 14:32:20",
    level: "DEBUG",
    message: "llama_new_context_with_model: KV self size = 80.00 MiB, K (f16): 40.00 MiB, V (f16): 40.00 MiB",
  },
  { timestamp: "2024-01-15 14:32:20", level: "WARN", message: "Model loading took longer than expected: 5.2s" },
  { timestamp: "2024-01-15 14:32:25", level: "ERROR", message: "Failed to connect to model: connection timeout" },
]

export function DeveloperHub() {
  const [logs, setLogs] = useState(mockLogs)
  const [searchTerm, setSearchTerm] = useState("")
  const [selectedLevel, setSelectedLevel] = useState("ALL")
  const [isRunning, setIsRunning] = useState(true)
  const [cpuUsage, setCpuUsage] = useState(15)
  const [memoryUsage, setMemoryUsage] = useState(2.4)

  const filteredLogs = logs.filter((log) => {
    const matchesSearch = log.message.toLowerCase().includes(searchTerm.toLowerCase())
    const matchesLevel = selectedLevel === "ALL" || log.level === selectedLevel
    return matchesSearch && matchesLevel
  })

  const getLevelColor = (level: string) => {
    switch (level) {
      case "ERROR":
        return "text-red-400"
      case "WARN":
        return "text-yellow-400"
      case "DEBUG":
        return "text-blue-400"
      default:
        return "text-green-400"
    }
  }

  return (
    <div className="p-6 h-full flex flex-col">
      {/* Header */}
      <div className="mb-6">
        <h1 className="text-3xl font-bold">Developer Hub</h1>
        <p className="text-gray-400 mt-1">Monitor application state and logs</p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-4 gap-6 flex-1">
        {/* Main Log Panel */}
        <div className="lg:col-span-3 flex flex-col">
          <Card className="bg-[#1B2636] border-[#2A3B52] flex-1 flex flex-col">
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle className="flex items-center space-x-2">
                  <Code className="w-5 h-5 text-[#4A90E2]" />
                  <span>Real-time Log Stream</span>
                </CardTitle>
                <div className="flex items-center space-x-2">
                  <Button
                    size="sm"
                    variant={isRunning ? "destructive" : "default"}
                    onClick={() => setIsRunning(!isRunning)}
                  >
                    {isRunning ? <Square className="w-4 h-4" /> : <Play className="w-4 h-4" />}
                    {isRunning ? "Stop" : "Start"}
                  </Button>
                  <Button size="sm" variant="outline">
                    <RefreshCw className="w-4 h-4" />
                  </Button>
                </div>
              </div>

              {/* Search and Filter */}
              <div className="flex space-x-2 mt-4">
                <div className="relative flex-1">
                  <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400" />
                  <Input
                    placeholder="Search logs..."
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                    className="pl-10 bg-[#0F1419] border-[#2A3B52]"
                  />
                </div>
                <div className="flex space-x-1">
                  {["ALL", "INFO", "WARN", "ERROR", "DEBUG"].map((level) => (
                    <Button
                      key={level}
                      size="sm"
                      variant={selectedLevel === level ? "default" : "outline"}
                      onClick={() => setSelectedLevel(level)}
                    >
                      {level}
                    </Button>
                  ))}
                </div>
              </div>
            </CardHeader>
            <CardContent className="flex-1 overflow-hidden">
              <div className="h-full bg-[#0F1419] rounded-lg p-4 font-mono text-sm overflow-auto scrollbar-thin">
                {filteredLogs.map((log, index) => (
                  <div key={index} className="flex space-x-4 py-1 hover:bg-[#1B2636] px-2 -mx-2 rounded">
                    <span className="text-gray-500 text-xs whitespace-nowrap">{log.timestamp}</span>
                    <Badge variant="outline" className={`text-xs ${getLevelColor(log.level)} border-current`}>
                      {log.level}
                    </Badge>
                    <span className="text-gray-300 flex-1">{log.message}</span>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Sidebar Panels */}
        <div className="space-y-6">
          {/* Application State */}
          <Card className="bg-[#1B2636] border-[#2A3B52]">
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <Activity className="w-5 h-5 text-[#4A90E2]" />
                <span>App State</span>
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-400">Status</span>
                <Badge className="bg-green-500">Running</Badge>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-400">Current Task</span>
                <span className="text-sm">Idle</span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-400">Active Models</span>
                <span className="text-sm">3</span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-400">Uptime</span>
                <span className="text-sm">2h 15m</span>
              </div>
            </CardContent>
          </Card>

          {/* Resource Usage */}
          <Card className="bg-[#1B2636] border-[#2A3B52]">
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <Cpu className="w-5 h-5 text-[#4A90E2]" />
                <span>Resources</span>
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <div className="flex items-center justify-between mb-2">
                  <span className="text-sm text-gray-400">CPU Usage</span>
                  <span className="text-sm">{cpuUsage}%</span>
                </div>
                <div className="w-full bg-[#0F1419] rounded-full h-2">
                  <div
                    className="bg-[#4A90E2] h-2 rounded-full transition-all duration-300"
                    style={{ width: `${cpuUsage}%` }}
                  />
                </div>
              </div>

              <div>
                <div className="flex items-center justify-between mb-2">
                  <span className="text-sm text-gray-400">Memory</span>
                  <span className="text-sm">{memoryUsage} GB</span>
                </div>
                <div className="w-full bg-[#0F1419] rounded-full h-2">
                  <div
                    className="bg-green-500 h-2 rounded-full transition-all duration-300"
                    style={{ width: `${(memoryUsage / 8) * 100}%` }}
                  />
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}
