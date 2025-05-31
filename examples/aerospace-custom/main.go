package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/cristianoliveira/aerospace-ipc"
)

func main() {
	socketPath := fmt.Sprintf("/tmp/bobko.%s-%s.sock", "aerospace", os.Getenv("USER"))
	client, err := aerospace.NewAeroSpaceCustomConnection(
		aerospace.AeroSpaceCustomConnectionOpts{
			SocketPath:      socketPath,
			ValidateVersion: true,
		},
	)
	if err != nil {
		if errors.Is(err, aerospace.ErrVersionMismatch) {
			fmt.Printf("[WARN] %s\n", err)
		} else {
			log.Fatalf("Failed to connect: %v", err)
		}
	}
	defer client.CloseConnection()

	windows, err := client.GetAllWindows()
	if err != nil {
		log.Fatalf("Failed to get windows: %v", err)
	}

	for i, window := range windows {
		fmt.Printf("%d) %s\n", i, window)
	}

	if len(windows) > 0 {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("What window to focus? (empty to cancel): ")
		text, _ := reader.ReadString('\n')
		input := strings.TrimSpace(text)
		if input == "" {
			fmt.Println("No window selected, exiting.")
		} else {
			index, err := strconv.Atoi(input)
			if err != nil {
				log.Fatalf("Invalid input: %v", err)
			}

			if index < 0 || index >= len(windows) {
				log.Fatalf("No window with index %d", index)
			}

			err = client.SetFocusByWindowID(windows[index].WindowID)
			if err != nil {
				log.Fatalf("Failed to focus on window: %v", err)
			}
		}
	}

	fmt.Println("Listed all windows successfully.")
}
