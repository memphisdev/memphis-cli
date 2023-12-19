[![Github (6)](https://github.com/memphisdev/memphis/assets/107035359/bc2feafc-946c-4569-ab8d-836bc0181890)](https://www.functions.memphis.dev/)
<p align="center">
<a href="https://memphis.dev/discord"><img src="https://img.shields.io/discord/963333392844328961?color=6557ff&label=discord" alt="Discord"></a>
<a href="https://github.com/memphisdev/memphis/issues?q=is%3Aissue+is%3Aclosed"><img src="https://img.shields.io/github/issues-closed/memphisdev/memphis?color=6557ff"></a> 
  <img src="https://img.shields.io/npm/dw/memphis-dev?color=ffc633&label=installations">
<a href="https://github.com/memphisdev/memphis/blob/master/CODE_OF_CONDUCT.md"><img src="https://img.shields.io/badge/Code%20of%20Conduct-v1.0-ff69b4.svg?color=ffc633" alt="Code Of Conduct"></a> 
<img alt="GitHub release (latest by date)" src="https://img.shields.io/github/v/release/memphisdev/memphis?color=61dfc6">
<img src="https://img.shields.io/github/last-commit/memphisdev/memphis?color=61dfc6&label=last%20commit">
</p>

 <b><p align="center">
  <a href="https://memphis.dev/pricing/">Cloud</a> - <a href="github.com/memphisdev/memphis-dev-academy">Academy</a> - <a href="https://memphis.dev/docs/">Docs</a> - <a href="https://twitter.com/Memphis_Dev">X</a> - <a href="https://www.youtube.com/channel/UCVdMDLCSxXOqtgrBaRUHKKg">YouTube</a>
</p></b>

<div align="center">

  <h4>

**[Memphis.dev](https://memphis.dev)** Is The First Data Streaming Platform Designed For Backend Developers<br>
To Build Event-driven And Real-time Features Faster Than Ever.<br>

<img width="177" alt="cloud_native 2 (5)" src="https://github.com/memphisdev/memphis/assets/107035359/a20ea11c-d509-42bb-a46c-e388c8424101">

  </h4>
  
</div>

## Installation
Mac
```sh
$ brew tap memphisdev/memphis-cli
$ brew install memphis-cli
```

Linux - RPM
```sh
$ wget https://github.com/memphisdev/memphis-cli/releases/latest/download/memphis-cli_linux_amd64.rpm
$ sudo rpm -i memphis-cli_linux_amd64.rpm
```

Linux - APK
```sh
$ wget https://github.com/memphisdev/memphis-cli/releases/latest/download/memphis-cli_linux_arm64.apk
$ apk add memphis-cli_linux_arm64.apk --allow-untrusted
```

Windows - Powershell
```sh
$ powershell -c "Invoke-WebRequest -Uri 'https://github.com/memphisdev/memphis-cli/releases/latest/download/memphis-cli_Windows_x86_64.zip'  -OutFile './memphis-cli_Windows_x86_64.zip'"
$ powershell -c "Expand-Archive memphis-cli_Windows_x86_64.zip -DestinationPath memphis-cli -Force"
```

## Upgrade
Mac
```sh
$ brew upgrade memphis-cli
```

Linux - RPM
```sh
$ wget https://github.com/memphisdev/memphis-cli/releases/latest/download/memphis-cli_linux_amd64.rpm
$ sudo rpm -U memphis-cli_linux_amd64.rpm
```

Linux - APK
```sh
$ wget https://github.com/memphisdev/memphis-cli/releases/latest/download/memphis-cli_linux_arm64.apk
$ apk add memphis-cli_linux_arm64.apk --allow-untrusted
```

Windows - Powershell
```sh
$ powershell -c "Invoke-WebRequest -Uri 'https://github.com/memphisdev/memphis-cli/releases/latest/download/memphis-cli_Windows_x86_64.zip'  -OutFile './memphis-cli_Windows_x86_64.zip'"
$ powershell -c "Expand-Archive memphis-cli_Windows_x86_64.zip -DestinationPath memphis-cli -Force"
```

## Functions
### Create a basic Memphis function template.
```sh
$ mem func init myExampleFunc --lang [nodejs/go/python]
```

## Benchmark
#### Overview
The Memphis Benchmarking Tool is designed to evaluate the performance of Memphis producer/consumer under various configurations. It allows for detailed specification of message size, batch processing, concurrency, and more

### Producer
#### Usage
To benchmark a Memphis producer, use the following command structure:
```sh
$ mem bench producer --message-size 128 --count 1000 --concurrency 1 --host <host> --account-id <account-id(not needed for open-source)> --user <client type user> --password <password>
```

### available flags
- **station:** Specify the station for message production (default: benchmark-station).
- **partition-key:** Define a partition key for message production. Takes precedence over partition-number.
- **partition-number:** The desired partition number to which the messages will be produced, default is 1
- **producer-name:** Name the producer (default: p-bench).
- **message-size:** Set message size in bytes (range: 128 to 8,388,608 [8MB]). Random data generated if empty.
- **count:** Specify the number of messages to produce (default: 1).
- **message:** Define a custom message. Random data generated if empty.
- **message:** Whether to wait for an acknowledgement for every message, default is false
- **concurrency:** Set the number of concurrent producers (default: 1).
- **host:** Specify the Memphis host (default: localhost).
- **account-id:** Account ID for Memphis server (not required for open-source edition).
- **user:** Username for Memphis server access (client type user).
- **password:** Password for Memphis server access.


### Consumer
#### Usage
To benchmark a Memphis consumer, use the following command structure:
```sh
$ mem bench consumer --message-size 128 --count 1000 --concurrency 1 --batch-size 50 --host <host> --account-id <account-id(not needed for open-source)> --user <client type user> --password <password>
```

### available flags
- **station:** Specify the station for message production (default: benchmark-station).
- **partition-key:** Define a partition key for message production. Takes precedence over partition-number.
- **consumer-name:** Set a custom name for the consumer (default: c-bench).
- **group:** Name the consumer's group (default: cg-bench).
- **batch-size:** Determine the batch size (default: 10).
- **batch-max-wait-time:** Max wait time (in milliseconds) for a batch (default: 1).
- **message-size:** Set message size in bytes (range: 128 to 8,388,608 [8MB]). Random data generated if empty.
- **count:** Specify the number of messages to produce (default: 1).
- **message:** Define a custom message. Random data generated if empty.
- **concurrency:** Set the number of concurrent producers (default: 1).
- **producer-name:** Name the producer (default: p-bench).
- **host:** Specify the Memphis host (default: localhost).
- **account-id:** Account ID for Memphis server (not required for open-source edition).
- **user:** Username for Memphis server access (client type user).
- **password:** Password for Memphis server access.
