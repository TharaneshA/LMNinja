import React, { useState } from 'react';
import { App } from '../../wailsjs/go/app/App';

interface LaunchpadProps {
  setResultText: (text: string) => void;
}

const Launchpad: React.FC<LaunchpadProps> = ({ setResultText }) => {
  const [scanResult, setScanResult] = useState<string>("No scan initiated yet.");

  const startScan = async () => {
    setScanResult("Starting red team scan... Please wait.");
    setResultText("Starting red team scan... Please wait.");
    try {
      const result = await App.StartRedTeamScan();
      setScanResult(result);
      setResultText(result);
    } catch (error) {
      setScanResult(`Error during scan: ${error}`);
      setResultText(`Error during scan: ${error}`);
    }
  };

  return (
    <div className="launchpad p-6 bg-gray-800 rounded-lg shadow-xl">
      <h1 className="text-4xl font-bold text-white mb-6">Launch a Red Team Scan</h1>
      <p className="text-lg text-gray-300 mb-8">
        Initiate a red team scan against your target LLM. This will run a series of predefined attacks and evaluate the LLM's responses.
      </p>

      <button
        onClick={startScan}
        className="bg-green-600 hover:bg-green-700 text-white font-bold py-3 px-6 rounded-lg transition duration-300 mb-8"
      >
        Start Red Team Scan
      </button>

      <div className="scan-results bg-gray-700 p-5 rounded-lg shadow-md">
        <h2 className="text-xl font-semibold text-purple-400 mb-3">Scan Result</h2>
        <p className="text-gray-200 whitespace-pre-wrap">{scanResult}</p>
      </div>
    </div>
  );
};

export default Launchpad;