import { useEffect, useState } from "react";
import { HashRouter, Routes, Route } from "react-router-dom";
import { EventsOn } from "@wailsjs/runtime";
import { GetConnections } from "@wailsjs/go/app/App";
import { useAppStore } from "./lib/state";
import { Toaster } from "./components/ui/sonner";
import { Loader2 } from "lucide-react";
import Sidebar from "./components/Sidebar";
import Dashboard from "./components/Dashboard";
import Connectors from "./components/Connectors";
import AttackLibrary from "./components/AttackLibrary";
import DeveloperHub from "./components/DeveloperHub";
import Settings from "./components/Settings";
import ThemeProvider from "./components/ThemeProvider";

function MainApp() {
    return (
        <div className="flex h-screen w-full">
            <Sidebar />
            <main className="flex-1 h-full overflow-y-auto">
                <Routes>
                    <Route path="/" element={<Dashboard />} />
                    <Route path="/connectors" element={<Connectors />} />
                    <Route path="/attacks" element={<AttackLibrary />} />
                    <Route path="/developer" element={<DeveloperHub />} />
                    <Route path="/settings" element={<Settings />} />
                </Routes>
            </main>
        </div>
    );
}

export default function App() {
    const { setConnections } = useAppStore();
    const [isHydrated, setIsHydrated] = useState(false);

    useEffect(() => {
        EventsOn("wails:ready", async () => {
            try {
                const connections = await GetConnections();
                setConnections(connections || []);
            } catch (error) {
                console.error("Failed to hydrate connections:", error);
            } finally {
                setIsHydrated(true);
            }
        });
    }, [setConnections]);

    return (
        <ThemeProvider>
            <HashRouter>
                {!isHydrated ? (
                    <div className="flex h-screen w-full items-center justify-center bg-[#0a0a0f] text-white">
                        <div className="flex flex-col items-center gap-4">
                            <Loader2 className="w-12 h-12 animate-spin text-[#4A90E2]" />
                            <p className="text-lg text-gray-400">Initializing LMNinja...</p>
                        </div>
                    </div>
                ) : (
                    <MainApp />
                )}
            </HashRouter>
            <Toaster position="bottom-right" />
        </ThemeProvider>
    );
}