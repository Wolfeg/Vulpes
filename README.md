# Vulpes - Send Gitlab pipeline events to Telegram chat

## Description

App listens for Gitlab pipeline webhooks and send alerts to Telegram

## Settings

You need to set this environment variables:

| Env name             | Description                    |
| ---------------------|--------------------------------|
| `JIRAURI`            | Jira url (for link formatting) |
| `JIRAPROJECTCODE`    | Jira project code              |
| `GITLABURI`          | Gitlab url                     |
| `GITLABTGCHATID`     | Telegram Chat ID               |
| `GITLABTGTOKEN`      | Telegram bot token             |
| `GITLABTGLISTENPORT` | Port for listening webhooks    |
| `GITLABTGSECRET`     | Gitlab webhook secret phrase   |

## Launching

Example with Docker Compose

```yml
version: "3.8"
services:
  vulpes:
    container_name: vulpes
    image: "vulpes:latest"
    env_file: vulpes.env
    ports:
      - "3005:3005"
    restart: unless-stopped
```
