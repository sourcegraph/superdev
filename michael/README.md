
```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o superdev-amprunner

docker build . -t superdev-worker

docker run -v /Users/michael/IdeaProjects/godot-reg:/workdir/repo -v ./guidance:/workdir/guidance -e ./connection:/workdir/connection -e ANTHROPIC_API_KEY=... superdev
```
