"use client"

import { useState } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Switch } from "@/components/ui/switch"
import { Label } from "@/components/ui/label"
import { Badge } from "@/components/ui/badge"
import { Palette, Download, Trash2, RefreshCw, CheckCircle, Info } from "lucide-react"

export function Settings() {
  const [darkMode, setDarkMode] = useState(true)
  const [autoUpdates, setAutoUpdates] = useState(true)
  const [notifications, setNotifications] = useState(true)

  const handleExportData = () => {
    console.log("Exporting scan data...")
  }

  const handleClearHistory = () => {
    console.log("Clearing scan history...")
  }

  const handleCheckUpdates = () => {
    console.log("Checking for updates...")
  }

  return (
    <div className="p-6 max-w-4xl">
      {/* Header */}
      <div className="mb-6">
        <h1 className="text-3xl font-bold">Settings</h1>
        <p className="text-gray-400 mt-1">Configure your LMNinja application</p>
      </div>

      <div className="space-y-6">
        {/* Appearance Section */}
        <Card className="bg-[#1B2636] border-[#2A3B52]">
          <CardHeader>
            <CardTitle className="flex items-center space-x-2">
              <Palette className="w-5 h-5 text-[#4A90E2]" />
              <span>Appearance</span>
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center justify-between">
              <div className="space-y-1">
                <Label htmlFor="dark-mode">Dark Mode</Label>
                <p className="text-sm text-gray-400">Use dark theme for the interface</p>
              </div>
              <Switch id="dark-mode" checked={darkMode} onCheckedChange={setDarkMode} />
            </div>

            <div className="flex items-center justify-between">
              <div className="space-y-1">
                <Label htmlFor="notifications">Desktop Notifications</Label>
                <p className="text-sm text-gray-400">Show notifications for scan results</p>
              </div>
              <Switch id="notifications" checked={notifications} onCheckedChange={setNotifications} />
            </div>
          </CardContent>
        </Card>

        {/* Updates Section */}
        <Card className="bg-[#1B2636] border-[#2A3B52]">
          <CardHeader>
            <CardTitle className="flex items-center space-x-2">
              <RefreshCw className="w-5 h-5 text-[#4A90E2]" />
              <span>Updates</span>
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center justify-between">
              <div className="space-y-1">
                <Label>Current Version</Label>
                <div className="flex items-center space-x-2">
                  <Badge variant="outline">v1.0.0</Badge>
                  <CheckCircle className="w-4 h-4 text-green-500" />
                  <span className="text-sm text-green-500">Up to date</span>
                </div>
              </div>
              <Button
                variant="outline"
                onClick={handleCheckUpdates}
                className="flex items-center space-x-2 bg-transparent"
              >
                <RefreshCw className="w-4 h-4" />
                <span>Check for Updates</span>
              </Button>
            </div>

            <div className="flex items-center justify-between">
              <div className="space-y-1">
                <Label htmlFor="auto-updates">Automatic Updates</Label>
                <p className="text-sm text-gray-400">Automatically download and install updates</p>
              </div>
              <Switch id="auto-updates" checked={autoUpdates} onCheckedChange={setAutoUpdates} />
            </div>
          </CardContent>
        </Card>

        {/* Data Section */}
        <Card className="bg-[#1B2636] border-[#2A3B52]">
          <CardHeader>
            <CardTitle className="flex items-center space-x-2">
              <Download className="w-5 h-5 text-[#4A90E2]" />
              <span>Data Management</span>
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center justify-between">
              <div className="space-y-1">
                <Label>Export Scan Data</Label>
                <p className="text-sm text-gray-400">Download all scan results and reports</p>
              </div>
              <Button
                variant="outline"
                onClick={handleExportData}
                className="flex items-center space-x-2 bg-transparent"
              >
                <Download className="w-4 h-4" />
                <span>Export All Data</span>
              </Button>
            </div>

            <div className="flex items-center justify-between">
              <div className="space-y-1">
                <Label>Clear Scan History</Label>
                <p className="text-sm text-gray-400">Remove all stored scan results</p>
              </div>
              <Button variant="destructive" onClick={handleClearHistory} className="flex items-center space-x-2">
                <Trash2 className="w-4 h-4" />
                <span>Clear History</span>
              </Button>
            </div>
          </CardContent>
        </Card>

        {/* About Section */}
        <Card className="bg-[#1B2636] border-[#2A3B52]">
          <CardHeader>
            <CardTitle className="flex items-center space-x-2">
              <Info className="w-5 h-5 text-[#4A90E2]" />
              <span>About</span>
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <p className="text-sm text-gray-400">
                LMNinja is a comprehensive LLM red teaming tool designed for security professionals and developers to
                test and evaluate the safety of language models.
              </p>
              <div className="flex items-center space-x-4 text-sm text-gray-400">
                <span>Version: 1.0.0</span>
                <span>•</span>
                <span>Build: 2024.01.15</span>
                <span>•</span>
                <span>Go/Wails Runtime</span>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
