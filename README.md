# ChatGPT Slack bot
Slack bot that utilizes ChatGPT `davinci-003` model to get real time text completion.
## Configuration
Use the `.env.example` file to set up the configuration of your Slack app and OpenAI API key, and rename the file to `.env`
### Setup Slack
* Create a Slack app and obtain auth token (`SLACK_AUTH_TOKEN` env variable).
* Enable [socket mode](https://api.slack.com/apis/connections/socket) for the application.
* [Configure Event Subscription](https://api.slack.com/apis/connections/events-api#subscribing) and add `app_mention` event to it, this will allow us to listen to messages that mention the bot.
* Create app token by going to General Information -> Generate Tokens and Scopes under App-Level Tokens, make sure to give it `connections:write` scope (`SLACK_APP_TOKEN` env variable).
* Invite the app to the channel where you want to use the bot by typing `/invite @AppName` in the channel.

### Obtain OpenAI API key
Get [API key](https://platform.openai.com/account/api-keys) from your OpenAI account (`OPENAI_API_KEY` env variable).

## Run the program
After setting up the `.env` file run
```
docker build -t chatgpt-slack-bot .
docker run -d chatgpt-slack-bot
```

See in the logs that everything is ok:

`docker logs $[CONTAINER_ID]`

You can see more logs by setting the DEBUG environment variable: `DEBUG="true"`.

You can also modify the [max tokens](https://platform.openai.com/docs/introduction/tokens) used with the `CHAT_GPT_MAX_TOKENS` environment variable.

That's it, now go to the channel and talk with the bot by mention it in the channel :)