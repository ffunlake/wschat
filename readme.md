# WsChat - WebSocket Chat Application
## Installation

### Prerequisites

1. Go Environment
   - Install Go 1.24 or later from [official Go website](https://golang.org/dl/)
   - Verify installation with `go version`
   - Set up GOPATH environment variable

2. SQLite3
   - For Debian/Ubuntu: `sudo apt-get install sqlite3`
   - For CentOS/RHEL: `sudo yum install sqlite`
   - For macOS: `brew install sqlite3`
   - For Windows: Download from [SQLite website](https://www.sqlite.org/download.html)
   - Verify installation with `sqlite3 --version`

3. Run Installation Script
   ```bash
   # Make script executable
   chmod +x shell/install.sh
   
   # Run installation script
   ./shell/install.sh
   ```
4. Build from Source
   ```bash
   # Build the binary in root directory
   go build -o wschat
   
   # For Windows
   go build -o wschat.exe
   ```
## Running the Application

1. Configure Application
   - Check and modify `config/chat.yaml` configuration file
   - Key configurations:
     ```yaml
     server:
       port: 8080    # WebSocket server port
     database:
       type: sqlite  # Database type
       path: data/wschat.db  # Database file path

2. Start Server
   ```bash
   # Run server in foreground
   export DEEPSEEK_API_KEY=[apikey] & ./wschat server./wschat server

   # Or with custom configuration file
   export DEEPSEEK_API_KEY=[apikey] & ./wschat server --config path_to_config_file
   ```

3. Start Client
   ```bash
   # Connect to local server
   ./wschat client

   # Or with custom configuration file
   ./wschat client  --config path_to_config_file
   ```
## Websocket Client Commands
The client supports the following commands:

1. `/register <name>`
   - Register as a new customer with the given name
   - Example: `/register Alice`

2. `/login <customer_name>`
   - Login with your customer Name
   - Example: `/login Alice`

3. `/chat <message>`
   - Send a message to chat with the AI
   - Example: `/chat Hello, how are you?`

4. `/aichat <message>`
   - Send a message to chat with the AI
   - Example: `/aichat What is the weather like today?`

5. `/history [limit]`
   - View chat history, optionally specify number of messages to show
   - Example: `/history 10`

6. `/feedback <rating> <command>`
   - Provide feedback function (rating: 1-5)
   - Example: `/feedback 5 very good`


Note: Commands are case-sensitive and must start with a forward slash (/). The <> indicates required parameters and [] indicates optional parameters. 

## Http restful api
The following HTTP REST API endpoints are available:

1. GET `/message/list`
   - Retrieves chat history for a specific customer
   - Query Parameters:
     - `customer_name`: Name of the customer to get history for
   - Example: `http://localhost:8080/message/list?customer_name=Alice`
   - Returns JSON array of messages

## Directory Description
- `cmd/`: Contains command-line interface code for server and client
- `config/`: Configuration files and config loading logic
- `data/`: Database files and data storage
- `docs/`: Project documentation and requirements
- `pkg/`: Core packages
  - `database/`: Database connection and management
  - `orm/`: Data models and ORM mappings
  - `services/`: Business logic services
- `shell/`: Shell scripts for installation and maintenance
- `wschat/`: Binary output directory


