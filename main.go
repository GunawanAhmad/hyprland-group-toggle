package main

import (
	"encoding/json"
	"net"
	"os"
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
	defer conn.Close()
	_, err = conn.Write([]byte(command))
	if err != nil {
		return nil, 0, err
	}
	buffer := make([]byte, 10000000)
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, 0, err
	}
	return buffer, n, nil
}

func main() {
	// Step to change current workspace to group in hyprland
	// 1. Create socket connection to hyprland IPC
	// 2. Get current workspace
	// 3. Get all window in current workspace
	// 4. Set current active window into a group
	// 5. check if current active window is in group
	// 6. If not, change current workspace to group workspace and move all window to active group
	// 7. If yes, remove the group

	// Step 1: create socket connection to hyprland IPC
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

	// Step 3: Get all windows in current workspace
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

	// TODO: Step 4: Set current active window into a group
	buffer, n, err = sendCommand("j/activewindow", socketPath)
	if err != nil {
		panic("failed to get active window: " + err.Error())
	}
	var activeWindow Window
	err = json.Unmarshal(buffer[:n], &activeWindow)
	if err != nil {
		panic("failed to marshal response to JSON: " + err.Error())
	}

	// TODO: Step 5: Check if current active window is in group

	// TODO: Step 6: If not, change current workspace to group workspace and move all window to active group

	// move all windows to the group workspace
	for _, window := range windowInWorkspace {
		if window.Address == activeWindow.Address {
			continue // skip the active window, it is already in the group
		}
		// focus the window
		_, _, err = sendCommand("dispatch focuswindow address:"+window.Address, socketPath)
		if err != nil {
			panic("failed to focus window: " + err.Error())
		}
		// move the window to the group window
		_, _, err = sendCommand("dispatch moveintogroup l", socketPath)
		if err != nil {
			panic("failed to move window into group: " + err.Error())
		}
		_, _, err = sendCommand("dispatch moveintogroup r", socketPath)
		if err != nil {
			panic("failed to move window into group: " + err.Error())
		}
		_, _, err = sendCommand("dispatch moveintogroup u", socketPath)
		if err != nil {
			panic("failed to move window into group: " + err.Error())
		}
		_, _, err = sendCommand("dispatch moveintogroup b", socketPath)
		if err != nil {
			panic("failed to move window into group: " + err.Error())
		}

		// set active as last current active window
		_, _, err = sendCommand("dispatch focuswindow address:"+activeWindow.Address, socketPath)
		if err != nil {
			panic("failed to focus window: " + err.Error())
		}
	}
}
