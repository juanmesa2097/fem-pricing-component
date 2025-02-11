package cmd

import (
	"fmt"
	"os/exec"

	"github.com/lector-org/lector/internal/errx"
	"github.com/lector-org/lector/internal/sysinfo"
	"github.com/lector-org/lector/pkg/inspector"
	"github.com/urfave/cli/v2"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type OpenCommandOptions struct {
	App *application.App
}

// NewOpenCommand is a CLI command that opens the specified file or website in a new graphical user interface (GUI) lector window.
func NewOpenCommand(options OpenCommandOptions) *cli.Command {
	return &cli.Command{
		Name:      "open",
		Usage:     "Opens the specified file in a new graphical user interface (GUI) lector window",
		UsageText: "lector open [OPTIONS] [INPUT_SOURCE]",
		Flags:     []cli.Flag{},
		Action: func(cCtx *cli.Context) error {
			filePath := cCtx.Args().Get(0)
			if filePath == "" {
				return errx.NewCliError(cCtx).GetNoArgError()
			}

			fileInspector := inspector.NewFileInspector(filePath)
			if !fileInspector.IsValid() {
				return errx.NewCliError(cCtx).GetCustomError(
					fmt.Sprintf("requires a valid file path. %q was not found", filePath),
				)
			}

			options.App.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
				Title: "Lector",
				Linux: application.LinuxWindow{
					WebviewGpuPolicy: func() application.WebviewGpuPolicy {
						if sysinfo.HasNvidiaGPU(exec.Command) {
							return application.WebviewGpuPolicyNever // Workaround for https://github.com/wailsapp/wails/issues/2977
						}

						return application.WebviewGpuPolicyAlways
					}(),
				},
			})

			return options.App.Run()
		},
	}
}
