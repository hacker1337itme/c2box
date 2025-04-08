# c2box
DISCORD C2 IN GOLANG FOR WIN SYSTEM

# WHAT DOES :

    Custom Command Execution: Run arbitrary commands with $custom_command.
    Change Directory: Change the working directory with $cd <path>.
    Screenshot Capture: Take a screenshot and send it to the channel with $screenshot.
    IP Address Information: Retrieve IP configuration using $ip.
    System Information: Fetch system details with $sysinfo.
    File Retrieval: Grab files and send them with $filegrab <filepath>.
    FodHelper Execution: Run a specific PowerShell script (ensure the safety of this command).
    Shutdown and Restart: Control the machine's shutdown and restart with $shutdown and $restart.
    List Files: List files in the current directory with $list.
    Print Current Working Directory: Print the current working directory with $pwd.
    Ping Command: Ping a specific hostname with $ping <hostname>.
    Create Files: Create files with optional content using $createfile <filename> [content].
    Delete Files: Delete specified files with $delete <filename>.

Important Instructions:

    Token and Channel ID: Be sure to fill in your Discord bot's token and the channel ID where you want the bot to send messages.
    Dependencies: Install the necessary Go packages.

# BUILD
      
```shell
go get github.com/bwmarrin/discordgo
go get github.com/go-vgo/robotgo
go mod init c2box.go
go mod tidy
go build -o c2box c2box.go
```
