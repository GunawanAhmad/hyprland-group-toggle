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

	conn, err := net.Dial("unix", socketPath)

	if err != nil {
		panic("failed to connect to Hyprland IPC socket: " + err.Error())
	}
	// errcheck: ignore error on close, as we are just cleaning up
	// defer conn.Close()

	// Step 2: Get current workspace
	// Prepare the command to get the active workspace
	command := "j/activeworkspace"
	_, err = conn.Write([]byte(command))
	if err != nil {
		panic("failed to send command to Hyprland IPC: " + err.Error())
	}
	// Read the response from the socket
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		panic("failed to read resonse from Hyprland IPC: " + err.Error())
	}
	var workspace Workspace
	err = json.Unmarshal(buffer[:n], &workspace)
	if err != nil {
		panic("failed to marshal response to JSON: " + err.Error())
	}

	// Step 3: Get all windows in current workspace
	// Prepare the command to get all windows in the current workspace
	conn, err = net.Dial("unix", socketPath)
	if err != nil {
		panic("failed to connect to Hyprland IPC socket: " + err.Error())
	}
	command = "j/clients"
	_, err = conn.Write([]byte(command))
	if err != nil {
		panic("failed to send command to Hyprland IPC: " + err.Error())
	}
	// Read the response from the socket
	buffer = make([]byte, 4096)
	n, err = conn.Read(buffer)
	if err != nil {
		panic("failed to read response from Hyprland IPC: " + err.Error())
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

	// Step 4: Set current active window into a group
	command = "dispatch togglegroup"
	conn, err = net.Dial("unix", socketPath)
	if err != nil {
		panic("failed to connect to Hyprland IPC socket: " + err.Error())
	}
	_, err = conn.Write([]byte(command))

	if err != nil {
		panic("failed to send command to Hyprland IPC: " + err.Error())
	}

	// TODO: Step 5: Check if current active window is in group

	// TODO: Step 6: If not, change current workspace to group workspace and move all window to active group

	// move all windows to the group workspace
	for _, window := range windowInWorkspace {
		// focus the window
		command = "dispatch focuswindow address:" + window.Address
		conn, err = net.Dial("unix", socketPath)
		if err != nil {
			panic("failed to connect to Hyprland IPC socket: " + err.Error())
		}
		_, err = conn.Write([]byte(command))
		if err != nil {
			panic("failed to send command to Hyprland IPC: " + err.Error())
		}
		// move the window to the group window
		command = "dispatch moveintogroup l"
		conn, err = net.Dial("unix", socketPath)
		if err != nil {
			panic("failed to connect to Hyprland IPC socket: " + err.Error())
		}
		_, err = conn.Write([]byte(command))
		if err != nil {
			panic("failed to send command to Hyprland IPC: " + err.Error())
		}

		command = "dispatch moveintogroup r"
		conn, err = net.Dial("unix", socketPath)
		if err != nil {
			panic("failed to connect to Hyprland IPC socket: " + err.Error())
		}
		_, err = conn.Write([]byte(command))
		if err != nil {
			panic("failed to send command to Hyprland IPC: " + err.Error())
		}
		
		command = "dispatch moveintogroup u"
		conn, err = net.Dial("unix", socketPath)
		if err != nil {
			panic("failed to connect to Hyprland IPC socket: " + err.Error())
		}
		_, err = conn.Write([]byte(command))
		if err != nil {
			panic("failed to send command to Hyprland IPC: " + err.Error())
		}

		command = "dispatch moveintogroup b"
		conn, err = net.Dial("unix", socketPath)
		if err != nil {
			panic("failed to connect to Hyprland IPC socket: " + err.Error())
		}
		_, err = conn.Write([]byte(command))
		if err != nil {
			panic("failed to send command to Hyprland IPC: " + err.Error())
		}
	}

}
