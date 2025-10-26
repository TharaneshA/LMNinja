# LMNinja

LMNinja is a desktop application built with Wails, Go, and React (TypeScript) designed to help red teamers and security professionals evaluate the security of Large Language Models (LLMs). It provides a framework for running various attacks, analyzing model vulnerabilities, and managing configurations for different LLM providers.

## Features

- **LLM Connector Management**: Easily configure and switch between different LLM providers (e.g., OpenAI, Ollama, custom APIs).
- **Attack Library**: A curated and extensible library of red teaming prompts and attack techniques (e.g., jailbreaks, prompt injections, data exfiltration).
- **Automated Scanning**: Run automated scans against target LLMs using predefined or custom attack scenarios.
- **Vulnerability Reporting**: Generate detailed reports on identified vulnerabilities and model weaknesses.
- **Evaluation Engine**: Tools to evaluate LLM responses against expected security criteria.
- **User-friendly Interface**: A clean and intuitive React-based frontend for easy interaction.

<img width="1002" height="720" alt="image" src="https://github.com/user-attachments/assets/b55af4f0-1834-439e-9a56-5a1d5cae75e4" />

<img width="998" height="720" alt="image" src="https://github.com/user-attachments/assets/c78db3a1-06bc-4eb5-a59f-60273579b2f5" />

<img width="1919" height="1032" alt="image" src="https://github.com/user-attachments/assets/72c7e69c-a8e6-4610-b279-c9d4317008a1" />


## Project Structure

```
lmninja/
├── .github/
│   └── workflows/
│       └── build.yml         # GitHub Action to auto-build for Win/macOS/Linux
├── .gitignore
├── LICENSE
├── README.md                 # Your project's detailed README
├── go.mod
├── go.sum
├── wails.json                # Wails project configuration
│
├── cmd/
│   └── lmninja/
│       └── main.go           # Application entry point (package main)
│
├── app/
│   └── app.go                # The main App struct and all bound Go methods
│
├── internal/                 # All core backend Go packages
│   ├── config/
│   │   └── config.go         # Loading/saving user configurations
│   ├── engine/
│   │   ├── engine.go         # Red Teaming Engine
│   │   └── evaluator.go      # Evaluation Engine
│   ├── llm/
│   │   ├── llm.go            # The universal LLM interface
│   │   ├── openai.go         # OpenAI client
│   │   └── ollama.go         # Ollama client
│   └── storage/
│       └── database.go       # SQLite database logic
│
└── frontend/                 # The React project created by Wails
    ├── node_modules/
    ├── public/
    ├── src/
    │   ├── assets/           # Images (logo), fonts, CSS
    │   ├── components/       # Reusable React components (Button, Chart, Table)
    │   ├── hooks/            # Custom React hooks (e.g., useScanResults)
    │   ├── pages/            # Components for each screen (Dashboard, Settings)
    │   ├── services/         # Logic for calling Go methods
    │   ├── state/            # State management (Zustand or React Context)
    │   ├── App.tsx           # Main React App component
    │   └── main.tsx          # React entry point
    ├── index.html
    ├── package.json
    ├── tsconfig.json
    └── vite.config.ts        # Vite configuration for the frontend
```

## Getting Started

### Prerequisites

- Go (1.21 or later)
- Node.js (18 or later)
- Wails CLI (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/TharaneshA/LMNinja.git
   cd LMNinja
   ```

2. **Install Go dependencies:**
   ```bash
   go mod tidy
   ```

3. **Install frontend dependencies:**
   ```bash
   cd frontend
   npm install
   cd ..
   ```

### Running the Application

To run the application in development mode:

```bash
wails dev
```

This will start the Go backend and the React frontend, with live reloading for both.

### Building the Application

To build a production-ready executable:

```bash
wails build
```

The executable will be generated in the `build/bin` directory.

## Contributing

We welcome contributions! Please see our `CONTRIBUTING.md` (coming soon) for more details.

## License

This project is licensed under the MIT License - see the `LICENSE` file for details.
