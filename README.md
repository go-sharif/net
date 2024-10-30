# go-sharif-net

![go-sharif-net](https://img.shields.io/github/v/release/go-sharif/net?label=Release) ![License](https://img.shields.io/github/license/go-sharif/net)

`go-sharif-net` is a command-line interface (CLI) tool written in Go, designed to simplify the authentication process for accessing the internet at Sharif University of Technology. By automating the login process to `net2.sharif.edu`, this tool ensures seamless and efficient network access management.

<img width="1542" alt="image" src="https://github.com/user-attachments/assets/ceacb96d-a535-47c7-b538-6eddced6ef69">

## Table of Contents

- [Features](#features)
  - [Implemented Features](#implemented-features)
  - [Features to Be Implemented](#features-to-be-implemented)
- [Installation](#installation)
- [Usage](#usage)
  - [Commands](#commands)
  - [Examples](#examples)
- [Configuration](#configuration)
- [Development](#development)
- [Contributing](#contributing)
- [License](#license)

## Features

### Implemented Features

- **Login Command**
  - Authenticate to `net2.sharif.edu` using username and password.
  - Support for specifying configuration files.
  - Option to use IP instead of domain for authentication.
  - Interactive session monitoring with the ability to stay alive and manage session termination.
  
- **Help and Autocompletion**
  - Comprehensive help texts for all commands and flags.
  - Generate shell autocompletion scripts for various shells.

### Features to Be Implemented

| Feature                              | Description                                                   | Progress       |
|--------------------------------------|---------------------------------------------------------------|----------------|
| **Logout Command**                   | Allow users to terminate their network session manually.      | To Do          |
| **Automatic Reconnection**           | Automatically attempt to reconnect if the session drops.      | Done 🚧        |
| **Enhanced Logging**                 | Provide detailed logs for authentication and session status.  | To Do          |
| **Cross-Platform Support**           | Ensure compatibility with Windows, macOS, and Linux.          | To Do |
| **Configuration Management**         | Simplify configuration through environment variables or a GUI. | In Progress 🚧   |

## Installation

### Prerequisites

- Go (version 1.16 or later)

### Steps

1. **Clone the Repository**

   ```bash
   git clone https://github.com/drippypale/go-sharif-net.git
   cd go-sharif-net
   ```

2. **Build the Application**

   ```bash
   go build -o go-sharif-net
   ```

3. **Move the Binary to Your PATH**

   ```bash
   sudo mv go-sharif-net /usr/local/bin/
   ```

   Alternatively, you can use package managers like `brew` or `apt` if supported in the future.

## Usage

`go-sharif-net` provides a simple CLI interface to manage your network authentication. Below are the available commands and examples of how to use them.

### Commands

- **login**
  - Authenticate to the network.
  
- **help**
  - Display help information about commands.
  
- **completion**
  - Generate shell autocompletion scripts.

### Global Flags

- `--config string`  
  Specify the configuration file (default is `$HOME/.go-sharif-net.yaml`).
  
- `--use-ip`  
  Use the IP address specified in the configuration instead of the domain.

### Login Command Flags

- `-u, --username string`  
  Username for login. **(Required)**
  
- `-p, --password string`  
  Password for login. **(Required)**
  
- `-a, --alive`  
  Keep the session alive to monitor and maintain the connection.

### Examples

1. **Login with Username and Password**

   ```bash
   go-sharif-net login -u your_username -p your_password
   ```

2. **Login and Keep the Session Alive**

   ```bash
   go-sharif-net login -u your_username -p your_password -a
   ```

3. **Use a Custom Configuration File**

   ```bash
   go-sharif-net login --config /path/to/config.yaml -u your_username -p your_password
   ```

4. **Generate Autocompletion Script for Bash**

   ```bash
   go-sharif-net completion bash > /etc/bash_completion.d/go-sharif-net
   ```

## Permissions and Running with `sudo`

Certain features of `go-sharif-net`, such as **Automatic Reconnection** and **Internet Connection Monitoring**, require elevated privileges to function correctly. This is because these features utilize ICMP (Internet Control Message Protocol) operations, which typically require root or administrative access.

### Why `sudo` is Necessary

- **ICMP Operations:** Monitoring the internet connection involves sending ICMP packets (ping) to check the connectivity status. Most operating systems restrict ICMP operations to users with root privileges to enhance security.
- **Automatic Reconnection:** To seamlessly reconnect to the network without manual intervention, the CLI needs to continuously monitor the connection status, which involves privileged network operations.

### How to Run `go-sharif-net` with `sudo`

To enable these advanced features, you should run the CLI with `sudo`. Below are examples of how to execute the `login` command with elevated privileges:

```bash
sudo go-sharif-net login -u your_username -p your_password -a


## Configuration

`go-sharif-net` allows customization through a configuration file. By default, it looks for a file named `.go-sharif-net.yaml` in the home directory.

### Example Configuration (`~/.go-sharif-net.yaml`)

```yaml
username: your_username
password: your_password
use_ip: false
session_duration: 8h
```

**Note:** Storing passwords in plaintext configuration files can be a security risk. Consider using environment variables or secure storage mechanisms to handle sensitive information.

## Development

### Project Structure

```plaintext
go-sharif-net/
├── cmd/                # Command implementations
│   └── login.go
├── internal/
│   ├── http/           # HTTP handlers
│   └── ui/             # User interface components
├── util/               # Utility functions
├── go.mod
├── go.sum
└── README.md
```

### Code Highlights

- **Cobra & Viper:** Utilizes Cobra for CLI commands and Viper for configuration management.
- **Session Management:** Handles session login and maintains it with optional monitoring.
- **Concurrency:** Implements goroutines and channels to manage asynchronous tasks like pinging and session checking.


### Building from Source

1. **Clone the Repository**

   ```bash
   git clone https://github.com/drippypale/go-sharif-net.git
   cd go-sharif-net
   ```

2. **Install Dependencies**

   ```bash
   go mod download
   ```

3. **Run the Application**

   ```bash
   go run main.go login -u your_username -p your_password
   ```

4. **Run Tests**

   ```bash
   go test ./...
   ```

## Contributing

Contributions are welcome! Please follow these steps to contribute:

1. **Fork the Repository**

2. **Create a Feature Branch**

   ```bash
   git checkout -b feature/YourFeature
   ```

3. **Commit Your Changes**

   ```bash
   git commit -m "Add your feature"
   ```

4. **Push to the Branch**

   ```bash
   git push origin feature/YourFeature
   ```

5. **Open a Pull Request**

Please ensure your code follows the project's coding standards and includes appropriate tests.

## License

This project is licensed under the [MIT License](LICENSE).

---

**Note:** This README is a living document. As the project evolves, please keep this file updated to reflect the latest changes and features.
