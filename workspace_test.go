package aerospace

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	mock_client "github.com/cristianoliveira/aerospace-ipc/internal/mocks"
	"github.com/cristianoliveira/aerospace-ipc/pkg/client"
	"go.uber.org/mock/gomock"
)

func TestWorkspace(t *testing.T) {
	t.Run("Happy path", func(tt *testing.T) {
		t.Run("GetFocusedWorkspace", func(tt *testing.T) {
			ctrl := gomock.NewController(tt)
			defer ctrl.Finish()

			mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
			aeroSpaceWM := AeroSpaceWM{Conn: mockConn}

			workspaces := []Workspace{
				{Workspace: "42"},
			}

			dataJSON, err := json.Marshal(workspaces)
			if err != nil {
				t.Fatalf("failed to marshal windows response: %v", err)
			}

			mockConn.EXPECT().
				SendCommand(
					"list-workspaces",
					[]string{
						"--focused",
						"--json",
					},
				).
				Return(
					&client.Response{
						StdOut: string(dataJSON),
					},
					nil,
				)

			workspace, err := aeroSpaceWM.GetFocusedWorkspace()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if workspace.Workspace != "42" {
				t.Fatalf("expected workspace '42', got '%s'", workspace.Workspace)
			}
		})

		t.Run("MoveWindowToWorkspace", func(tt *testing.T) {
			ctrl := gomock.NewController(tt)
			defer ctrl.Finish()

			mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
			aeroSpaceWM := AeroSpaceWM{Conn: mockConn}

			windowID := "12345"
			workspace := "42"

			mockConn.EXPECT().
				SendCommand(
					"move-node-to-workspace",
					[]string{
						workspace,
						"--window-id", windowID,
					},
				).
				Return(
					&client.Response{
						StdOut: "",
					},
					nil,
				)

			intWindowID, err := strconv.Atoi(windowID)
			if err != nil {
				t.Fatalf("failed to convert window ID to int: %v", err)
			}

			err = aeroSpaceWM.MoveWindowToWorkspace(intWindowID, workspace)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	})

	t.Run("Error cases", func(tt *testing.T) {
		t.Run("GetFocusedWorkspace return error", func(tt *testing.T) {
			ctrl := gomock.NewController(tt)
			defer ctrl.Finish()

			mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
			aeroSpaceWM := AeroSpaceWM{Conn: mockConn}

			mockConn.EXPECT().
				SendCommand(
					"list-workspaces",
					[]string{
						"--focused",
						"--json",
					},
				).
				Return(nil, fmt.Errorf("no focused workspace found")).
				Times(1)

			_, err := aeroSpaceWM.GetFocusedWorkspace()
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})

		t.Run("GetFocusedWorkspace return empty", func(tt *testing.T) {
			ctrl := gomock.NewController(tt)
			defer ctrl.Finish()

			mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
			aeroSpaceWM := AeroSpaceWM{Conn: mockConn}

			mockConn.EXPECT().
				SendCommand(
					"list-workspaces",
					[]string{
						"--focused",
						"--json",
					},
				).
				Return(&client.Response{StdOut: "[]"}, nil).
				Times(1)

			_, err := aeroSpaceWM.GetFocusedWorkspace()
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	})
}
