#### Keptn

##### GitHub issue self-assignment bot

- Description: The goal is to create a self-service issue assignment bot for GitHub contributors who are not yet part of the organization but would like to work on issues marked for external handling. The bot should be able to check if the user is part of the organization, examine if the pre-conditions for self-assignment are met (configurable labels or rules about issues already assigned/PRs opened), and assign the issue. Additionally, the bot should be able to track the state of the issue by adding/removing specific labels. The bot should be part of the CI and executed as an action on an issue change.
- Expected Outcome:
  - Implement GitHub bot in TypeScript/Golang
  - Bot does not require an external hosting
  - Bot is able to assign GitHub issues to contributors following the pre-defined set of rules
  - Bot is able to track the status of GitHub issues with labels
  - Introduce documentation about how to use and configure the bot
- Recommended Skills: GitHub API, TypeScript/Golang, Webhooks
- Expected project size: Large
- Mentor(s):
  - Ondrej Dubaj (@odubajDT, ondrej.dubaj@dynatrace.com) - primary
  - Florian Bacher (@bacherfl, florian.bacher@dynatrace.com)
- Upstream Issue (URL): https://github.com/keptn/lifecycle-toolkit/issues/2823