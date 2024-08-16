# github-activity-notify

GitHubのコントリビューション数をLine Notifyで知らせます。

![](https://raw.githubusercontent.com/Dodai-Dodai/github-activity-notify/branchImage/note/example.jpeg)

## 使い方

LINE notifyのTOKENとGitHubのAccess Tokenを`.env`にセットして使います

`.env`
```
LINE_TOKEN=YOUR_LINE_ACCESS_TOKEN
GITHUB_TOKEN=YOUR_GITHUB_ACCESS_TOKEN
GITHUB_USER=YOUE_GITHUB_USERNAME
```

```sh
$ go run .
```

## Extra

DockerHub:https://hub.docker.com/repository/docker/dodaidodai/github-activity-notify/general