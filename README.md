# SQL Playground

Uses Go-starter repo

## Installation

Run make commands to setup 

```
make setup          # Installs wgo, templ
make setup_mac      # Installs wgo, templ for mac
```

Running the server
```
docker-compose up   # Start the database
make dev            # Start the server
make dev port=9000  # Start a server at a specific port
```

Build the binary
```
make build
```