package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"strings"
)

type Workspace struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Window struct {
	Address   string    `json:"address"`
	Workspace Workspace `json:"workspace"`
}

func sendCommand(command string, socketPath string) ([]byte, int, error) {
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Println("Error closing connection:", cerr)
		}
	}()
	_, err = conn.Write([]byte(command))
	if err != nil {
		return nil, 0, err
	}
	buffer := make([]byte, 8192)
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, 0, err
	}
	return buffer, n, nil
}

func main() {
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	instanceSig := os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")
	if runtimeDir == "" || instanceSig == "" {
		panic("XDG_RUNTIME_DIR or HYPRLAND_INSTANCE_SIGNATURE is not set, are you using Hyprland?")
	}
	socketPath := runtimeDir + "/hypr/" + instanceSig + "/.socket.sock"
	buffer, n, err := sendCommand("j/activeworkspace", socketPath)
	if err != nil {
		panic("failed to connect to Hyprland IPC: " + err.Error())
	}
	var workspace Workspace
	err = json.Unmarshal(buffer[:n], &workspace)
	if err != nil {
		panic("failed to marshal response to JSON: " + err.Error())
	}

	// Prepare the command to get all windows in the current workspace
	buffer, n, err = sendCommand("j/clients", socketPath)
	if err != nil {
		panic("failed to get windows: " + err.Error())
	}
	var windows []Window
	err = json.Unmarshal(buffer[:n], &windows)
	if err != nil {
		panic("failed to marshal response to JSON: " + err.Error())
	}
	var windowInWorkspace []Window
	for _, window := range windows {
		if window.Workspace.ID == workspace.ID {
			windowInWorkspace = append(windowInWorkspace, window)
		}
	}

	_, _, err = sendCommand("dispatch togglegroup", socketPath)
	if err != nil {
		panic("failed to toggle group: " + err.Error())
	}

	if len(windowInWorkspace) == 0 {
		return
	}

	buffer, n, err = sendCommand("j/activewindow", socketPath)
	if err != nil {
		panic("failed to get active window: " + err.Error())
	}
	var activeWindow Window
	err = json.Unmarshal(buffer[:n], &activeWindow)
	if err != nil {
		panic("failed to marshal response to JSON: " + err.Error())
	}
	
	// BUG: Sometime there's a window that somehow not moved to the group. So it't very consistent to move all windows in the workspace to the group.
	var commands []string
	for _, window := range windowInWorkspace {
		// skip the active window
		if window.Address == activeWindow.Address {
			continue
		}
		for range 3 {
			commands = append(commands,
				"dispatch focuswindow address:"+window.Address,
				"dispatch moveintogroup l",
				"dispatch moveintogroup r",
				"dispatch moveintogroup u",
				"dispatch moveintogroup d",
			)
		}
	}

	// focus the active window last
	commands = append(commands, "dispatch focuswindow address:"+activeWindow.Address)

	dispatchBatch := "[[BATCH]]" + strings.Join(commands, ";")
	
	// TODO: Split the command into smaller chunks if it exceeds the maximum length
	_, _, err = sendCommand(dispatchBatch, socketPath)
	if err != nil {
		panic("failed to move windows to group: " + err.Error())
	}
}
