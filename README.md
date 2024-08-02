## Experimental CLI app for WSO2 Identity Server / Asgardeo Management 

is-cli is a command-line interface tool for managing and interacting with Identity Server/Asgardeo integrations.

## Features
- Authenticate as a machine or user
- Manage applications (only list, create yet)
- Interactive mode
- Keychain support for storing credentials
- Logging

## Stack
- Golang
- Cobra
- Bubbletea (for interactive mode)

## Installation

### Prerequisites

- Go 1.16 or higher
- Make sure `$HOME/bin` is in your PATH

### Steps

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/is-cli.git
   cd is-cli
   ```

2. Build and install the CLI:
   ```
   make install
   ```

3. Verify the installation:
   ```
   is-cli --version
   ```

If you encounter any issues, ensure that `$HOME/bin` is in your PATH by adding the following line to your shell configuration file (`~/.zshrc` for Zsh or `~/.bash_profile` for Bash):

```
export PATH=$PATH:$HOME/bin
```

Then, reload your shell configuration:
```
source ~/.zshrc  # or ~/.bash_profile for Bash
```

## Usage

Here are some example commands:

```
is-cli login
is-cli applications list
is-cli applications create
```

## Contributing

We welcome contributions! Please feel free to submit a Pull Request.

## Support

If you encounter any problems or have any questions, please open an issue on the GitHub repository.


