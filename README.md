# Telegram stickerpack editor
A stickerpack editor written in Go with a Vue frontend. Allows you to create and edit (not yet implemented) stickerpacks. Supports importing 7tv emotes and tenor gifs.

# Quickstart
You need to make your own `.env` based on `.env.example` file and use your own bot [token](https://core.telegram.org/bots/api#authorizing-your-bot) and [domain](https://core.telegram.org/widgets/login#linking-your-domain-to-the-bot). In order for the telegram auth widget to work, the website needs to be running on https://yourdomain. You need to reroute the domain to 127.0.0.1 in your HOSTS file. After that, run `docker-compose.dev.yml`.