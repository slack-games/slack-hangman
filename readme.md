
# Slack Hangman game

Turn based slack hangman game.

## Commands

Slack commands examples:

- ___/hng start___ - start a new game
- ___/hng guess [a-z] - make a guess
- ___/hng current___ - show the current game state
- ___/hng stats___ - show user stats, wins, losses etc [not implemented]
- ___/hng help___ - show user command help and how to play [not implemented]
- ___/hng ping___ - ping request, for development


## TODO

High level todos:

- Move the image URL path to config, currently hardcoded
- Slack message text improvements, make more personal
- Move the words into separate DB or use some other solutions to store the word list
