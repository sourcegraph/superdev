# superdev

## Run instructions
1. Start server
```bash
go run . server
```

2. Start UI
```bash
cd ui
npm i
npm start
```

4. Export your Anthropoic API key to env

3. Issue a cURL command to start a thread
```bash
curl -X POST http://localhost:8080/run \ 
  -H "Content-Type: application/json" \
  -d '{
    "repository_link": "git@github.com:sourcegraph/test-mcp.git",
    "prompt": "Who are you?",
    "docker_image": "superdev-wrapped-image"
```