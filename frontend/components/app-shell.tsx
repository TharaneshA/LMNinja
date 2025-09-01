"use client";

import { useEffect, useState } from "react";
import { EventsOn } from "@wailsjs/runtime";
import { GetConnections } from "@wailsjs/go/app/App";
import { useAppStore } from "@/lib/state";
import { toast } from "sonner";
import { Sidebar } from "@/components/sidebar";
import { Loader2 } from "lucide-react";

export function AppShell({ children }: { children: React.ReactNode }) {
  const { setConnections } = useAppStore();
  // We now have two states: is Wails ready, and is our data loaded.
  const [isWailsReady, setIsWailsReady] = useState(false);
  const [isHydrated, setIsHydrated] = useState(false);

  useEffect(() => {
    // Listen for the "wails:ready" event. This is the official signal
    // that the Go/JS bridge is initialized and safe to use.
    EventsOn("wails:ready", () => {
      console.log("[DEBUG] Wails is ready.");
      setIsWailsReady(true);
    });
  }, []);

  useEffect(() => {
    // This effect runs ONLY when Wails is confirmed to be ready.
    if (!isWailsReady) return;

    const hydrateConnections = async () => {
      console.log("[DEBUG] Wails is ready, now fetching connections...");
      try {
        const connections = await GetConnections();
        setConnections(connections || []);
        console.log(`[DEBUG] Hydration complete. Loaded ${connections?.length || 0} connections.`);
      } catch (error) {
        toast.error("Fatal Error: Could not load connections", { description: String(error) });
      } finally {
        setIsHydrated(true); // Mark hydration as complete, even if it failed.
      }
    };
    hydrateConnections();
  }, [isWailsReady, setConnections]);

  // Render a loading screen until the app is fully initialized and data is loaded.
  if (!isHydrated) {
    return (
      <div className="flex h-screen w-full items-center justify-center bg-[#0a0a0f] text-white">
        <div className="flex flex-col items-center gap-4">
          <Loader2 className="w-12 h-12 animate-spin text-[#4A90E2]" />
          <p className="text-lg text-gray-400">Initializing LMNinja...</p>
        </div>
      </div>
    );
  }

  // Once ready, render the main application interface.
  return (
    <div className="flex h-screen w-full">
      <Sidebar />
      <main className="flex-1 h-full overflow-y-auto">
        {children}
      </main>
    </div>
  );
}