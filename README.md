# AI Dashboard

[![Created by Freddy Ouellette](https://img.shields.io/badge/Created%20by%20Freddy%20Ouellette-gray)](https://freddyouellette.com) [![GitHub](https://img.shields.io/badge/GitHub-black?logo=github)](https://github.com/freddyouellette/ai-dashboard)

### One Dashboard for all of your personal AI Chatbots.

![AI Dashboard](docs/1.png)

## Installation & Usage
```
make app
```

## Development
There is a dev container already set up for development with VSCode. Just use the action `Remote-Containers: Reopen in Container`. Then you can run the following command:
```
make app-dev
```

## Plugins
* You can add custom API plugins to the `plugins/custom` directory, just add the built `.so` file(s) to the `AI_API_PLUGINS` env var.
* You can see examples of this in [`plugins/ai_apis/anthropic/anthropic.go`](plugins/ai_apis/anthropic/anthropic.go) and [`plugins/ai_apis/openai/openai.go`](plugins/ai_apis/openai/openai.go)

## Problems, Questions, Suggestions? 
* I encourage all issues or suggestions to be submitted through the [**Issues** tab on GitHub](https://github.com/freddyouellette/ai-dashboard/issues).
* Pull requests are welcome.

## Support Me
[![Donate](https://img.shields.io/badge/Donate-fec133?logo=paypal)](https://www.paypal.com/donate/?hosted_button_id=3PJ9XD363CC5E)

Bitcoin: `bc1qs39glh9cwsef0qv40dny6ajnweqe2le7ynfgr2`

Ethereum: `0x5Baba8708b8676afBFF2974b4af4894Fc12aE242`