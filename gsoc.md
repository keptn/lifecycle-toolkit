#### Keptn

##### GitHub issue self-assignment bot

- Description: The goal is to create a GitHub issue assignment self-service bot for contributors who are not yet part of the organization but would like to work on issues marked for external handling. The bot should be able to chek if the user is part of the organization, examine if the pre-conditions for self-assignment are met (specific labels; rules about issues already assigned/PRs opened) and assign the issue. Additionally the bot should be able to track the state of the issue by adding/removing the specific labels.
- Expected Outcome:
  - Implement GitHub bot in TypeScript
  - Bot is able to assign GitHub issues to contributors following the pre-defined set of rules
  - Bot is able to track the status of GitHub issues with lables
  - Introduce documentation abouth how to use and configure the bot
- Recommended Skills: GitHub API, TypeScript, Webhooks
- Expected project size: Large
- Mentor(s):
  - Ondrej Dubaj (@odubajDT, ondrej.dubaj@dynatrace.com) - primary
  - Rakshit Gondwal (@rakshitgondwal, john@email.address)
- Upstream Issue (URL): https://github.com/keptn/lifecycle-toolkit/issues/2823