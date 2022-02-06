# linux-alert
Every 10 minutes, memory, cpu and storage usage is checked and if they over 80%, sending alert via email.

### Usage
1) Create .env file from .env.example
2) Execute the binary

```
cp .env.example .env
wget -q https://github.com/melorthys/linux-alert/releases/download/v0.0.1/linux-alert-amd64
chmod +x linux-alert
./linux-alert
```

### Terminal Output
```
Getting system metrics..
CPU  = % 4.3
RAM  = % 31.8
DISK = % 4
Sleeping for 10 minutes..
```

### Email Content
```
CPU usage is 90.1% for hostname (192.168.1.50) at Sun, 06 Feb 2022 23:47:45 +01
```
