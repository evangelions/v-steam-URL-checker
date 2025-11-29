
# Steam ID Checker

This is a command-line tool written in Go that allows you to check the availability of Steam IDs.
The program will check each Steam ID from your `targets.txt` file. If an ID is available, it will be saved to the `output.txt` file. The program will also track your progress, so you can stop and resume your session at any time.

## Features

- **Check Steam ID availability**: Check if a Steam ID is available or already claimed.
- **Session management**: Create new sessions or resume existing ones.
- **Progress tracking**: The tool remembers where it left off in case of interruptions.
- **Target generator**: Generates thousands of 3c/3l/4c/4l usernames in `targets.txt`.
  
## Download 
get the win64 folder from the [release](https://github.com/evangelions/Vsteam/releases/tag/v1.4) page.

## Usage

When you run the tool, you'll be greeted with a menu to select your desired action. You can:

1. **Start a New Session**: Create a new session and begin checking Steam IDs from a `targets.txt` file.
2. **Resume an Existing Session**: Choose an existing session to continue checking IDs.
3. **Generate Random IDs**: Generates a list of random 3c/3l/4c/4l in the `targets.txt` file.
4. **Exit**: Exit the program.

## Issues
- checker will currently give false positives for shadowbanned accounts (working on fix)
