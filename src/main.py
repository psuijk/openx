#!/usr/bin/env python3
import json
import os
import sys
import subprocess
from pathlib import Path

CONFIG_FILE = Path.home() / ".openx_config.json"

def load_config():
    """Load the configuration file."""
    if not CONFIG_FILE.exists():
        return {}

    try:
        with open(CONFIG_FILE, 'r') as f:
            return json.load(f)
    except:
        return {}


def save_config(config):
    """Save the configuration file."""
    with open(CONFIG_FILE, 'w') as f:
        json.dump(config, f, indent=2)


def add_path(key, path=None):
    """Add a new path to the config."""
    if not key:
        print("Error: Please provide a key name")
        print("Usage: openx --add <key> [path]")
        return 1

    # Use current directory if no path provided
    if path is None:
        path = os.getcwd()

    # Resolve to absolute path
    path = os.path.abspath(os.path.expanduser(path))

    if not os.path.isdir(path):
        print(f"Error: Directory '{path}' does not exist")
        return 1


    config = load_config()
    config[key] = path
    save_config(config)
    print(f"✓ Added: {key} -> {path}")

    return 0


def remove_path(key):

    """Remove a path from the config."""
    if not key:
        print("Error: Please provide a key name")
        print("Usage: openx --remove <key>")
        return 1

    config = load_config()

    if key in config:
        del config[key]
        save_config(config)
        print(f"✓ Removed: {key}")
        return 0

    else:
        print(f"Error: Key '{key}' not found")
        return 1


def list_paths():
    """List all saved paths."""
    config = load_config()

    if config:
        print("Saved paths:")
        for key, path in config.items():
            print(f"  {key} -> {path}")

    else:
        print("No paths saved yet")

    return 0


def open_path(key):
    """Open a saved path in VS Code and Terminal."""

    if not key:
        print("Error: Please provide a key name")
        print("Usage: openx <key>")
        return 1

    config = load_config()

    if key not in config:
        print(f"Error: Key '{key}' not found")
        print("Use 'openx --list' to see available keys")
        return 1

    path = config[key]

    if not os.path.isdir(path):
        print(f"Error: Directory '{path}' no longer exists")
        return 1

    print(f"Opening: {path}")

    # Open in VS Code
    try:
        subprocess.run(['code', path], check=True, capture_output=True)
        print("✓ Opened in VS Code")
    except (subprocess.CalledProcessError, FileNotFoundError):
        print("⚠ VS Code command 'code' not found")
    
    # Open new Terminal window at the path
    applescript = f'''
    tell application "Terminal"
        do script "cd '{path}' && clear"
        activate
    end tell
    '''

    try:
        subprocess.run(['osascript', '-e', applescript], check=True, capture_output=True)
        print("✓ Opened in Terminal")
    except subprocess.CalledProcessError:
        print("⚠ Failed to open Terminal")

    return 0


def show_help():
    """Show help message."""
    print("openx - Quick project directory opener")
    print("")
    print("Usage:")
    print("  openx <key>              Open saved path in VS Code and Terminal")
    print("  openx --add <key> [path] Add current (or specified) directory")
    print("  openx --remove <key>     Remove a saved path")
    print("  openx --list             List all saved paths")
    print("  openx --help             Show this help message")
    print("")
    print("Examples:")
    print("  openx --add pbmares      # Add current directory as 'pbmares'")
    print("  openx --add myapp ~/Projects/myapp")
    print("  openx pbmares            # Open the pbmares directory")
    print("  openx --list             # Show all saved paths")

    return 0


def main():
    """Main entry point."""
    if len(sys.argv) < 2:
        return show_help()

    command = sys.argv[1]

    if command == '--add':
        key = sys.argv[2] if len(sys.argv) > 2 else None
        path = sys.argv[3] if len(sys.argv) > 3 else None
        return add_path(key, path)

    
    elif command in ['--remove', '--rm']:
        key = sys.argv[2] if len(sys.argv) > 2 else None
        return remove_path(key)

    elif command in ['--list', '--ls', '-l']:
        return list_paths()


    elif command in ['--help', '-h']:
        return show_help()

    else:
        # Treat as a key to open
        return open_path(command)

if __name__ == '__main__':
    sys.exit(main())

