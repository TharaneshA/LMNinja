import { createContext, useContext, useEffect, useState } from "react";

type ThemeProviderProps = {
  children: React.ReactNode;
  defaultTheme?: string;
  storageKey?: string;
};

type ThemeProviderState = {
  theme: string;
  setTheme: (theme: string) => void;
};

const ThemeProviderContext = createContext<ThemeProviderState>({
  theme: "dark",
  setTheme: () => null,
});

export default function ThemeProvider({
  children,
  defaultTheme = "dark",
  storageKey = "vite-ui-theme",
}: ThemeProviderProps) {
  const [theme, setTheme] = useState(() => (localStorage.getItem(storageKey) || defaultTheme));

  useEffect(() => {
    const root = window.document.documentElement;
    root.classList.remove("light", "dark");
    root.classList.add(theme);
  }, [theme]);

  const value = {
    theme,
    setTheme: (theme: string) => {
      localStorage.setItem(storageKey, theme);
      setTheme(theme);
    },
  };

  return (
    <ThemeProviderContext.Provider value={value}>
      {children}
    </ThemeProviderContext.Provider>
  );
}

export const useTheme = () => {
  const context = useContext(ThemeProviderContext);
  if (context === undefined) {
    throw new Error("useTheme must be used within a ThemeProvider");
  }
  return context;
};