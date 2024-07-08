package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"ypeskov/go-orgfin/internal/config"
	"ypeskov/go-orgfin/internal/logger"
	"ypeskov/go-orgfin/internal/server"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		fmt.Sprintf(fmt.Sprintf("cannot read config: %s", err))
		panic(err)
	}

	appLogger := logger.New(cfg)

	appServer := server.New(cfg, appLogger)

	openBrowser(fmt.Sprintf("http://localhost:%s", cfg.Port))

	err = appServer.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}

}

func openBrowser(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = append(args, "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = "open"
		args = append(args, url)
	case "linux":
		cmd = "xdg-open"
		args = append(args, url)
	default:
		fmt.Printf("Unsupported platform: %s\n", runtime.GOOS)
		return
	}
	fmt.Printf("Opening browser at %s\n", url)
	err := exec.Command(cmd, args...).Start()
	if err != nil {
		fmt.Printf("Failed to open browser: %v\n", err)
	}
}
