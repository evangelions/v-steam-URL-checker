package main

import (
    "bufio"
    "fmt"
    "math/rand"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "sync"
    "time"

    "resty.dev/v3"
)

var (
    Reset   = "\033[0m"
    Red     = "\033[31m"
    Green   = "\033[32m"
    Yellow  = "\033[33m"
    Blue    = "\033[34m"
    Magenta = "\033[35m"
    Cyan    = "\033[36m"
    Gray    = "\033[37m"
    White   = "\033[97m"
)

// Your single rotating proxy
const proxyURL = ""

func checkID(id string) bool {
    client := resty.New().
        SetTimeout(10 * time.Second).
        SetProxy(proxyURL)

    url := "https://steamcommunity.com/id/" + id
    resp, err := client.R().Get(url)
    if err != nil {
        fmt.Printf(Yellow+"[%s] Request error: %v\n"+Reset, id, err)
        return false
    }

    body := resp.String()
    return !strings.Contains(body, "The specified profile could not be found.")
}

func pauseTerminal() {
    fmt.Println("\nPress Enter to exit...")
    bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func getAllSessions() ([]string, string) {
    sessionsDir := "sessions"
    files, err := os.ReadDir(sessionsDir)
    if err != nil {
        return nil, ""
    }
    var sessions []string
    var latestSession string
    maxSession := 0
    for _, file := range files {
        if file.IsDir() {
            sessionName := file.Name()
            if num, err := strconv.Atoi(strings.TrimPrefix(sessionName, "SESSION_")); err == nil {
                if num > maxSession {
                    maxSession = num
                    latestSession = sessionName
                }
                sessions = append(sessions, sessionName)
            }
        }
    }
    return sessions, latestSession
}

func getSessionPath(sessionName string) string {
    return filepath.Join("sessions", sessionName)
}

func readTargets(filename string) ([]string, int, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, 0, err
    }
    defer file.Close()
    var ids []string
    scanner := bufio.NewScanner(file)
    progress := 0
    if scanner.Scan() {
        firstLine := strings.TrimSpace(scanner.Text())
        if strings.HasPrefix(firstLine, "Progress:") {
            p := strings.TrimSpace(strings.TrimPrefix(firstLine, "Progress:"))
            if n, err := strconv.Atoi(p); err == nil {
                progress = n
            }
        } else if firstLine != "" {
            ids = append(ids, firstLine)
        }
    }
    for scanner.Scan() {
        id := strings.TrimSpace(scanner.Text())
        if id != "" {
            ids = append(ids, id)
        }
    }
    return ids, progress, scanner.Err()
}

func updateProgress(filename string, progress int, ids []string) error {
    lines := []string{fmt.Sprintf("Progress: %d", progress)}
    lines = append(lines, ids[progress:]...)
    return os.WriteFile(filename, []byte(strings.Join(lines, "\n")+"\n"), 0644)
}

func showSplash() {
    fmt.Println(Red + `
------------------------------------------------------------------------------------------------------------------------

███████╗████████╗███████╗ █████╗ ███╗   ███╗    ██╗██████╗      ██████╗██╗  ██╗███████╗ ██████╗██╗  ██╗███████╗██████╗ 
██╔════╝╚══██╔══╝██╔════╝██╔══██╗████╗ ████║    ██║██╔══██╗    ██╔════╝██║  ██║██╔════╝██╔════╝██║ ██╔╝██╔════╝██╔══██╗
███████╗   ██║   █████╗  ███████║██╔████╔██║    ██║██║  ██║    ██║     ███████║█████╗  ██║     █████╔╝ █████╗  ██████╔╝
╚════██║   ██║   ██╔══╝  ██╔══██║██║╚██╔╝██║    ██║██║  ██║    ██║     ██╔══██║██╔══╝  ██║     ██╔═██╗ ██╔══╝  ██╔══██╗
███████║   ██║   ███████╗██║  ██║██║ ╚═╝ ██║    ██║██████╔╝    ╚██████╗██║  ██║███████╗╚██████╗██║  ██╗███████╗██║  ██║
╚══════╝   ╚═╝   ╚══════╝╚═╝  ╚═╝╚═╝     ╚═╝    ╚═╝╚═════╝      ╚═════╝╚═╝  ╚═╝╚══════╝ ╚═════╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝
                                                                                                                        
------------------------------------------------------------------------------------------------------------------------` + Green + `
STEAM ID AVAILABILITY CHECKER — Template by yTax - modified by v@maakima
Using rotating residential proxy via DataImpulse
` + Red + `
------------------------------------------------------------------------------------------------------------------------` + Reset)
}

func generateRandomIDs() {
    const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
    rand.Seed(time.Now().UnixNano())

    var ids []string
    for i := 0; i < 25000; i++ {
        id := string(chars[rand.Intn(len(chars))]) +
            string(chars[rand.Intn(len(chars))]) +
            string(chars[rand.Intn(len(chars))])
        ids = append(ids, id)
    }
    for i := 0; i < 50000; i++ {
        id := string(chars[rand.Intn(len(chars))]) +
            string(chars[rand.Intn(len(chars))]) +
            string(chars[rand.Intn(len(chars))]) +
            string(chars[rand.Intn(len(chars))])
        ids = append(ids, id)
    }

    content := strings.Join(ids, "\n") + "\n"
    err := os.WriteFile("targets.txt", []byte(content), 0644)
    if err != nil {
        fmt.Println(Red+"Failed to write targets.txt:"+Reset, err)
    } else {
        fmt.Println(Green + "Generated 75,000 random 3 & 4 character IDs → targets.txt" + Reset)
    }
    pauseTerminal()
}

func main() {
    showSplash()

    sessions, latestSession := getAllSessions()
    fmt.Println(Cyan + "\n-> Existing sessions:" + Reset)
    for _, session := range sessions {
        if session == latestSession {
            fmt.Printf(" - %s"+Blue+" (LATEST SESSION)\n"+Reset, session)
        } else {
            fmt.Printf(" - %s\n", session)
        }
    }
    fmt.Println(Red + `
+-----------------------+
|` + Green + ` 1. Start New Session` + Red + `  |
|` + Green + ` 2. Resume Session` + Red + `     |
|` + Green + ` 3. Generate Random IDs` + Red + `|
|` + Green + ` 4. Exit` + Red + `               |
+-----------------------+` + Reset)

    fmt.Print(Cyan + "\n-> Choose an option" + Reset + ": ")
    var choice string
    fmt.Scanln(&choice)

    var sessionPath, targetsPath, outputPath string
    var isNewSession bool

    switch choice {
    case "1":
        newSessionName := "SESSION_" + strconv.Itoa(len(sessions)+1)
        sessionPath = getSessionPath(newSessionName)
        targetsPath = filepath.Join(sessionPath, "targets.txt")
        outputPath = filepath.Join(sessionPath, "output.txt")
        isNewSession = true
    case "2":
        fmt.Print(Cyan + "-> Enter the session name (e.g. SESSION_1): " + Reset)
        var chosenSession string
        fmt.Scanln(&chosenSession)
        sessionPath = getSessionPath(chosenSession)
        targetsPath = filepath.Join(sessionPath, "targets.txt")
        outputPath = filepath.Join(sessionPath, "output.txt")
        isNewSession = false
    case "3":
        generateRandomIDs()
        return
    case "4":
        fmt.Println(Red + "Exiting program. Hope you found some good IDs!" + Reset)
        os.Exit(0)
    default:
        fmt.Println(Red + "Invalid choice. Please restart the program." + Reset)
        return
    }

    if isNewSession {
        if err := os.MkdirAll(sessionPath, os.ModePerm); err != nil {
            fmt.Println(Red+"Error creating session directory:"+Reset, err)
            pauseTerminal()
            return
        }
        input, err := os.ReadFile("targets.txt")
        if err != nil {
            fmt.Println(Red+"Error reading targets.txt:"+Reset, err)
            pauseTerminal()
            return
        }
        if err := os.WriteFile(targetsPath, input, 0644); err != nil {
            fmt.Println(Red+"Error copying targets.txt:"+Reset, err)
            pauseTerminal()
            return
        }
    }

    ids, progress, err := readTargets(targetsPath)
    if err != nil {
        fmt.Println(Red+"Error reading targets file:"+Reset, err)
        pauseTerminal()
        return
    }
    if progress >= len(ids) {
        fmt.Println(Green + "All IDs have already been checked!" + Reset)
        pauseTerminal()
        return
    }

    file, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Println(Red+"Error creating/opening output file:"+Reset, err)
        pauseTerminal()
        return
    }
    defer file.Close()

    // Thread selection
    fmt.Print(Cyan + "\n-> Number of threads (default 15): " + Reset)
    var threadInput string
    fmt.Scanln(&threadInput)
    threads := 15
    if threadInput != "" {
        if t, err := strconv.Atoi(threadInput); err == nil && t > 0 {
            threads = t
        }
    }
    fmt.Printf(Cyan+"Using %d threads with rotating residential proxy\n"+Reset, threads)

    fmt.Println(Cyan + "\nChecking Steam IDs...\n" + Reset)

    // Worker pool
    jobs := make(chan int, len(ids)-progress)
    var wg sync.WaitGroup

    for w := 0; w < threads; w++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for i := range jobs {
                id := ids[i]
                if !checkID(id) {
                    fmt.Printf(Green+"Available: %s\n"+Reset, id)
                    file.WriteString(id + "\n")
                } else {
                    fmt.Printf(Red+"Not available: %s\n"+Reset, id)
                }

                // Update progress
                newProgress := i + 1
                if err := updateProgress(targetsPath, newProgress, ids); err != nil {
                    fmt.Printf(Yellow+"Failed to update progress: %v\n"+Reset, err)
                }
            }
        }()
    }

    // Send remaining jobs
    for i := progress; i < len(ids); i++ {
        jobs <- i
    }
    close(jobs)

    wg.Wait()

    fmt.Println(Green + "\nCheck completed. Available IDs saved to " + outputPath + Reset)
    pauseTerminal()
}
