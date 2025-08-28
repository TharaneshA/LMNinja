"use client"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Shield, AlertTriangle, CheckCircle, XCircle, Eye, TrendingUp, Activity, Zap } from "lucide-react"

const recentScans = [
  { model: "GPT-4", score: 85, date: "2024-01-15", status: "completed" },
  { model: "Claude-3", score: 92, date: "2024-01-14", status: "completed" },
  { model: "Llama-2-70B", score: 78, date: "2024-01-14", status: "completed" },
  { model: "Mistral-7B", score: 88, date: "2024-01-13", status: "completed" },
  { model: "Gemini-Pro", score: 90, date: "2024-01-13", status: "completed" },
]

const connectedModels = [
  { name: "GPT-4", status: "connected", type: "OpenAI API" },
  { name: "Claude-3", status: "connected", type: "Anthropic API" },
  { name: "Llama-2-70B", status: "error", type: "Ollama Local" },
  { name: "Mistral-7B", status: "connected", type: "Ollama Local" },
]

const vulnerabilities = [
  { name: "Jailbreaks", count: 23, color: "bg-red-500" },
  { name: "PII Leaks", count: 15, color: "bg-orange-500" },
  { name: "Prompt Injection", count: 8, color: "bg-yellow-500" },
]

export function Dashboard() {
  const overallScore = 87

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Security Dashboard</h1>
          <p className="text-gray-400 mt-1">Monitor your LLM security posture</p>
        </div>
        <div className="flex items-center space-x-2">
          <Activity className="w-4 h-4 text-green-500" />
          <span className="text-sm text-gray-400">Live monitoring active</span>
        </div>
      </div>

      {/* Main Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Overall Security Score */}
        <Card className="bg-[#1B2636] border-[#2A3B52]">
          <CardHeader>
            <CardTitle className="flex items-center space-x-2">
              <Shield className="w-5 h-5 text-[#4A90E2]" />
              <span>Overall Security Score</span>
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center justify-center">
              <div className="relative w-32 h-32">
                <svg className="w-32 h-32 transform -rotate-90" viewBox="0 0 100 100">
                  <circle cx="50" cy="50" r="40" stroke="#2A3B52" strokeWidth="8" fill="none" />
                  <circle
                    cx="50"
                    cy="50"
                    r="40"
                    stroke="#4A90E2"
                    strokeWidth="8"
                    fill="none"
                    strokeDasharray={`${overallScore * 2.51} 251`}
                    strokeLinecap="round"
                  />
                </svg>
                <div className="absolute inset-0 flex items-center justify-center">
                  <span className="text-3xl font-bold">{overallScore}</span>
                </div>
              </div>
            </div>
            <div className="text-center mt-4">
              <p className="text-sm text-gray-400">Security Score</p>
              <div className="flex items-center justify-center mt-2">
                <TrendingUp className="w-4 h-4 text-green-500 mr-1" />
                <span className="text-sm text-green-500">+5 from last week</span>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Top Vulnerabilities */}
        <Card className="bg-[#1B2636] border-[#2A3B52]">
          <CardHeader>
            <CardTitle className="flex items-center space-x-2">
              <AlertTriangle className="w-5 h-5 text-orange-500" />
              <span>Top Vulnerabilities</span>
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            {vulnerabilities.map((vuln, index) => (
              <div key={index} className="flex items-center justify-between">
                <div className="flex items-center space-x-3">
                  <div className={`w-3 h-3 rounded-full ${vuln.color}`} />
                  <span className="text-sm">{vuln.name}</span>
                </div>
                <div className="flex items-center space-x-2">
                  <span className="text-sm font-medium">{vuln.count}</span>
                  <div className="w-20 h-2 bg-[#2A3B52] rounded-full overflow-hidden">
                    <div className={`h-full ${vuln.color}`} style={{ width: `${(vuln.count / 25) * 100}%` }} />
                  </div>
                </div>
              </div>
            ))}
          </CardContent>
        </Card>

        {/* Recent Scans */}
        <Card className="bg-[#1B2636] border-[#2A3B52]">
          <CardHeader>
            <CardTitle className="flex items-center space-x-2">
              <Activity className="w-5 h-5 text-[#4A90E2]" />
              <span>Recent Scans</span>
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {recentScans.map((scan, index) => (
                <div key={index} className="flex items-center justify-between p-3 bg-[#0F1419] rounded-lg">
                  <div className="flex items-center space-x-3">
                    <CheckCircle className="w-4 h-4 text-green-500" />
                    <div>
                      <p className="text-sm font-medium">{scan.model}</p>
                      <p className="text-xs text-gray-400">{scan.date}</p>
                    </div>
                  </div>
                  <div className="flex items-center space-x-3">
                    <Badge variant={scan.score >= 85 ? "default" : "destructive"}>{scan.score}</Badge>
                    <Button size="sm" variant="outline">
                      <Eye className="w-3 h-3 mr-1" />
                      View
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Connected Models */}
        <Card className="bg-[#1B2636] border-[#2A3B52]">
          <CardHeader>
            <CardTitle className="flex items-center space-x-2">
              <Zap className="w-5 h-5 text-[#4A90E2]" />
              <span>Connected Models</span>
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {connectedModels.map((model, index) => (
                <div key={index} className="flex items-center justify-between p-3 bg-[#0F1419] rounded-lg">
                  <div className="flex items-center space-x-3">
                    {model.status === "connected" ? (
                      <CheckCircle className="w-4 h-4 text-green-500" />
                    ) : (
                      <XCircle className="w-4 h-4 text-red-500" />
                    )}
                    <div>
                      <p className="text-sm font-medium">{model.name}</p>
                      <p className="text-xs text-gray-400">{model.type}</p>
                    </div>
                  </div>
                  <Badge variant={model.status === "connected" ? "default" : "destructive"}>{model.status}</Badge>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
