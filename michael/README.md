
```
docker build . -t superdev

docker run -e PROMPT="there's a repo at ./repo that you need to modify. find your task at ./guidance/Task.md, info on how to test at ./guidance/Validation.md." -v /Users/michael/IdeaProjects/godot-reg:/workdir/repo -v ./guidance:/workdir/guidance -e ANTHROPIC_API_KEY=... superdev
```

