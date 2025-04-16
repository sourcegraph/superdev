
```
docker build . -t superdev

docker run -e PROMPT="there's a repo at ./repo that you need to modify. find your task at ./guidance/Task.md, info on how to test at ./guidance/Validation.md. come up with a plan for a junior dev first, and then implement that. do only the minimal work required to meet the acceptance criteria." -v /Users/michael/IdeaProjects/godot-reg:/workdir/repo -v ./guidance:/workdir/guidance -e ANTHROPIC_API_KEY=... superdev
```

