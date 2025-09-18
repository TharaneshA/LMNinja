import { useState } from "react";
import { Link, useLocation } from "react-router-dom";
import { cn } from "@/lib/utils";
import { BarChart3, Zap, Code, Shield, Settings, ChevronLeft, ChevronRight } from "lucide-react";

const navigation = [
  { name: "Red Teaming", href: "/", icon: BarChart3 },
  { name: "LLM Connectors", href: "/connectors", icon: Zap },
  { name: "Attack Library", href: "/attacks", icon: Shield },
  { name: "Developer Hub", href: "/developer", icon: Code },
  { name: "Settings", href: "/settings", icon: Settings },
];

export default function Sidebar() {
  const [collapsed, setCollapsed] = useState(false);
  const location = useLocation();

  return (
    <div className={cn("bg-[#1B2636] border-r border-[#2A3B52] flex flex-col transition-all duration-300", collapsed ? "w-16" : "w-64")}>
      <div className="p-4 border-b border-[#2A3B52] flex items-center justify-between">
        {!collapsed && (
          <div className="flex items-center space-x-2">
            <div className="w-8 h-8 bg-[#4A90E2] rounded-lg flex items-center justify-center"><Shield className="w-5 h-5 text-white" /></div>
            <span className="font-semibold text-lg">LMNinja</span>
          </div>
        )}
        <button onClick={() => setCollapsed(!collapsed)} className="p-1 hover:bg-[#2A3B52] rounded-md transition-colors">
          {collapsed ? <ChevronRight className="w-4 h-4" /> : <ChevronLeft className="w-4 h-4" />}
        </button>
      </div>
      <nav className="flex-1 p-4 space-y-2">
        {navigation.map((item) => {
          const isActive = location.pathname === item.href;
          return (
            <Link key={item.name} to={item.href} className={cn("flex items-center space-x-3 px-3 py-2 rounded-lg transition-colors", isActive ? "bg-[#4A90E2] text-white" : "text-gray-300 hover:bg-[#2A3B52] hover:text-white")}>
              <item.icon className="w-5 h-5 flex-shrink-0" />
              {!collapsed && <span>{item.name}</span>}
            </Link>
          );
        })}
      </nav>
      <div className="p-4 border-t border-[#2A3B52]"><div className={cn("text-xs text-gray-400", collapsed ? "text-center" : "")}>{collapsed ? "v1.0" : "LMNinja v1.0.0"}</div></div>
    </div>
  );
}