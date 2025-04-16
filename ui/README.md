# SuperDev UI

A simple React UI for monitoring SuperDev threads and their outputs.

## Features

- Display all running threads in a grid layout
- Auto-updates thread outputs every second
- Responsive design that works on all screen sizes

## Getting Started

### Prerequisites

- Node.js 14+ and npm

### Installation

```bash
npm install
```

### Running the UI

```bash
npm start
```

The UI will be available at http://localhost:3000 and will automatically connect to the SuperDev server running on port 8080.

### Building for Production

```bash
npm run build
```

## Configuration

By default, the UI connects to `http://localhost:8080`. To change the API base URL, edit the `API_BASE_URL` constant in `src/App.js`.