#!/usr/bin/env node

import { spawn } from 'child_process';
import path from 'path';
import os from 'os';
import fs from 'fs';

function getBinaryPath(): string {
  const platform = os.platform();
  const arch = os.arch();
  
  let binaryName: string;
  if (platform === 'darwin') {
    binaryName = arch === 'arm64' 
      ? 'prometheus-mcp-server-darwin-arm64'
      : 'prometheus-mcp-server-darwin-amd64';
  } else if (platform === 'linux') {
    binaryName = arch === 'arm64'
      ? 'prometheus-mcp-server-linux-arm64'
      : 'prometheus-mcp-server-linux-amd64';
  } else if (platform === 'win32') {
    binaryName = 'prometheus-mcp-server-windows-amd64.exe';
  } else {
    throw new Error(`Unsupported platform: ${platform}`);
  }

  // 尝试多个可能的路径
  const possiblePaths = [
    path.join(__dirname, '..', 'bin', binaryName),
    path.join(__dirname, '..', '..', 'bin', binaryName),
    path.join(process.cwd(), 'bin', binaryName)
  ];

  for (const binPath of possiblePaths) {
    if (fs.existsSync(binPath)) {
      return binPath;
    }
  }

  throw new Error(`Could not find binary ${binaryName} in any of these locations:\n${possiblePaths.join('\n')}`);
}

try {
  const binaryPath = getBinaryPath();
  console.log(`Using binary: ${binaryPath}`);
  
  const child = spawn(binaryPath, process.argv.slice(2), {
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
} catch (err) {
  console.error('Error:', err instanceof Error ? err.message : String(err));
  process.exit(1);
} 