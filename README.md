# NetSleep

**NetSleep** is a small, cross-platform tool written in **Go** that automatically shuts down your computer once network activity drops below a certain threshold — perfect for letting downloads, uploads, or backups finish overnight.


## Installation & Build

### Clone the repository

```bash
git clone https://github.com/GunnarLuebbermann/netsleep.git
cd netsleep
```

### Get dependencies

```bash
go mod tidy
```

### Build (Windows, Linux, macOS)

use the provided build script:
```bash
make
```

The binaries will appear in the `build/` folder.


## Usage

Simply run:

```bash
netsleep.exe
```

**Default behavior:**
- Checks network usage every **10 seconds**
- If average speed stays below **50 KB/s** for **5 minutes**, the PC shuts down automatically

## Configuration

You can customize NetSleep’s behavior by creating a `config.json` file in the same directory as the executable (`netsleep.exe`).

This file allows you to adjust how often the network is checked, what counts as idle traffic, and how long the system should remain idle before shutting down.

### Example configuration

```json
{
    "interval": 10,
    "idle_threshold": 51200,
    "idle_time": 300,
    "shutdown_command": "shutdown /s /t 0"
}
```

## Antivirus Notice

Because NetSleep executes the `shutdown` command, some antivirus software (like **Windows Defender**) may flag it as *potentially unwanted software*.

If you built the binary yourself, it’s safe.  
You can fix false positives by:

- Building **without** stripping debug info:
  ```bash
  go build -o netsleep.exe
  ```
- Adding a Defender **exception** for the project folder  
- Uploading to [VirusTotal](https://www.virustotal.com/) to verify it’s clean  