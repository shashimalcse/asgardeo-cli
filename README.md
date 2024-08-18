<div align="center">

<picture>
  <source media="(prefers-color-scheme: light)" srcset="/docs/asgardeo-cli.png">
  <img alt="is cli logo" src="/docs/asgardeo-cli.png" width="50%" height="50%">
</picture>

asgardeo-cli is a experimental (non-official) cli app for managing and interacting with Asgardeo integrations.

</div>

## Features
- Authenticate as a machine (Client Credentials) or User (Device Flow)
- Manage applications
  - List applications
  - Create applications (Support Templates)
  - Delete applications
- Interactive mode
- Keychain support for storing credentials
- Logging

## Installation

### Prerequisites

- Go 1.16 or higher
- Make sure `$HOME/bin` is in your PATH

### Steps

1. Clone the repository:
   ```
   git clone https://github.com/shashimalcse/asgardeo-cli.git
   cd asgardeo-cli
   ```

2. Build and install the CLI:
   ```
   make install
   ```

3. Verify the installation:
   ```
   asgardeo --version
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

### Authenticating to Your Tenant

Authenticating to your Identity Server/ Asgardeo tenant is required for most functions of the CLI. It can be initiated by running:
```
asgardeo login
```

There are two ways to authenticate:

As a user - Recommended when invoking on a personal machine or other interactive environment. Facilitated by device authorization flow.
As a machine - Recommended when running on a server or non-interactive environments (ex: CI). Facilitated by client credentials flow.

> Authenticating as a user is not supported for Asgardeo tenants.

## Commands:

### Apps

- `asgardeo apps list` - List your applications
- `asgardeo apps create` - Create a new application
- `asgardeo apps delete <app-id>` - Delete an application

### API Resources

- `asgardeo apis list` - List your API resources


![Screenshot 2024-08-02 at 15 41 42](https://github.com/user-attachments/assets/c76a1b8e-740a-4ad7-a014-1a880b5a4f16)
![Screenshot 2024-08-02 at 15 43 22](https://github.com/user-attachments/assets/ebc9f872-65c7-4609-bd7f-926af2bac076)

## Contributing

We welcome contributions! Please feel free to submit a Pull Request.

## Support

If you encounter any problems or have any questions, please open an issue on the GitHub repository.


