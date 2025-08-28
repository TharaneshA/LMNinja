import { useState } from 'react';

export const useScanResult = () => {
  const [resultText, setResultText] = useState<string>("No scan initiated yet.");

  return { resultText, setResultText };
};