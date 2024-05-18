# github-activity-notify

GitHubのコントリビューション数をLine Notifyで知らせます。

![](https://raw.githubusercontent.com/Dodai-Dodai/github-activity-notify/branchImage/note/example.jpeg)

## 使い方

環境変数にコントリビューション数を返すAPIのURLとLINE Notifyのトークンを入れて実行

```sh
URL=https://github-contributions-api.deno.dev/[GitHub_User_ID].json \
TOKEN=[Your Line Notify Token] \
go run .
```

こちらのAPIを使っています:https://github.com/kawarimidoll/deno-github-contributions-api

## Extra

DockerHub:https://hub.docker.com/repository/docker/dodaidodai/github-activity-notify/general