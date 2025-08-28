import { useState } from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Launchpad from "./pages/Launchpad";
import Sidebar from "./components/Sidebar";
import Settings from "./pages/Settings";
import Dashboard from "./pages/Dashboard";

function App() {
    const [scanResultText, setScanResultText] = useState("");

    return (
        <Router>
            <div id="App" className="flex h-screen bg-gray-100">
                <Sidebar />
                <main className="flex-1 p-6 overflow-auto">
                    <Routes>
                        <Route path="/" element={<Dashboard />} />
                        <Route
                            path="/launchpad"
                            element={
                                <Launchpad
                                    setScanResultText={setScanResultText}
                                    setResultText={setScanResultText}
                                />
                            }
                        />
                        <Route path="/settings" element={<Settings />} />
                        <Route path="*" element={<Dashboard />} />
                    </Routes>
                </main>
            </div>
        </Router>
    );
}

export default App;