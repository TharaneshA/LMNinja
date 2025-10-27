; ===================================================================
; LMNinja Inno Setup Script - CORRECTED VERSION
; Packages the Wails executable and the self-contained PyInstaller backend.
; ===================================================================

; --- Preprocessor Defines ---
; Define all your app constants here. This is the correct way.
#define MyAppName "LMNinja"
#define MyAppVersion "0.1.0"
#define MyAppPublisher "Tharanesh A"
#define MyAppURL "https://github.com/TharaneshA/LMNinja"
#define MyOutputBaseFilename "LMNinja_Setup_v0.1.0"

[Setup]
; --- Application Identity ---
; Use the macros defined above.
AppId={{e2a3e3b8-3dbd-43f9-acca-bda38eb9c005}}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppPublisher={#MyAppPublisher}
AppPublisherURL={#MyAppURL}
AppSupportURL={#MyAppURL}
AppUpdatesURL="{#MyAppURL}/releases"

; --- Installation Directories & Output ---
DefaultDirName={autopf}\{#MyAppName}
DisableProgramGroupPage=yes
OutputBaseFilename={#MyOutputBaseFilename}
OutputDir=.\installers

; --- Compression & Style ---
Compression=lzma
SolidCompression=yes
WizardStyle=modern

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Tasks]
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"; Flags: unchecked

[Files]

; 1. Copy the main Wails application executable.
Source: "C:\project\LMNinja\build\bin\{#MyAppName}.exe"; DestDir: "{app}"; Flags: ignoreversion

; 2. Copy the ENTIRE self-contained Python engine directory.
Source: "C:\project\LMNinja\build\bin\lmninja-engine\*"; DestDir: "{app}\lmninja-engine"; Flags: recursesubdirs createallsubdirs

[Icons]
; Creates the Start Menu entry. This now correctly uses the preprocessor macro.
Name: "{autoprograms}\{#MyAppName}"; Filename: "{app}\{#MyAppName}.exe"
; Creates the optional Desktop icon. This now correctly uses the preprocessor macro.
Name: "{autodesktop}\{#MyAppName}"; Filename: "{app}\{#MyAppName}.exe"; Tasks: desktopicon

[Run]
; Offers to run the application after installation. This now correctly uses the preprocessor macro.
Filename: "{app}\{#MyAppName}.exe"; Description: "{cm:LaunchProgram,{#StringChange(MyAppName, '&', '&&')}}"; Flags: nowait postinstall skipifsilent