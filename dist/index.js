#!/usr/bin/env node
"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const child_process_1 = require("child_process");
const path_1 = __importDefault(require("path"));
const os_1 = __importDefault(require("os"));
const fs_1 = __importDefault(require("fs"));
function getBinaryPath() {
    const platform = os_1.default.platform();
    const arch = os_1.default.arch();
    let binaryName;
    if (platform === 'darwin') {
        binaryName = arch === 'arm64'
            ? 'prometheus-mcp-server-darwin-arm64'
            : 'prometheus-mcp-server-darwin-amd64';
    }
    else if (platform === 'linux') {
        binaryName = arch === 'arm64'
            ? 'prometheus-mcp-server-linux-arm64'
            : 'prometheus-mcp-server-linux-amd64';
    }
    else if (platform === 'win32') {
        binaryName = 'prometheus-mcp-server-windows-amd64.exe';
    }
    else {
        throw new Error(`Unsupported platform: ${platform}`);
    }
    // 尝试多个可能的路径
    const possiblePaths = [
        path_1.default.join(__dirname, '..', 'bin', binaryName),
        path_1.default.join(__dirname, '..', '..', 'bin', binaryName),
        path_1.default.join(process.cwd(), 'bin', binaryName)
    ];
    for (const binPath of possiblePaths) {
        if (fs_1.default.existsSync(binPath)) {
            return binPath;
        }
    }
    throw new Error(`Could not find binary ${binaryName} in any of these locations:\n${possiblePaths.join('\n')}`);
}
try {
    const binaryPath = getBinaryPath();
    console.log(`Using binary: ${binaryPath}`);
    const child = (0, child_process_1.spawn)(binaryPath, process.argv.slice(2), {
        stdio: 'inherit',
        shell: false
    });
    child.on('error', (err) => {
        console.error('Failed to start prometheus-mcp-server:', err);
        process.exit(1);
    });
    child.on('close', (code) => {
        process.exit(code ?? 0);
    });
}
catch (err) {
    console.error('Error:', err instanceof Error ? err.message : String(err));
    process.exit(1);
}
