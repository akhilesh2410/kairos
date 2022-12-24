package webui

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type FormData struct {
	CloudConfig        string `form:"cloud-config" json:"cloud-config" query:"cloud-config"`
	Reboot             string `form:"reboot" json:"reboot" query:"reboot"`
	PowerOff           string `form:"power-off" json:"power-off" query:"power-off"`
	InstallationDevice string `form:"installation-device" json:"installation-device" query:"installation-device"`
}

//go:embed public
var embededFiles embed.FS

func getFileSystem() http.FileSystem {
	fsys, err := fs.Sub(embededFiles, "public")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}

func Start(ctx context.Context, l string) error {

	ec := echo.New()
	assetHandler := http.FileServer(getFileSystem())
	ec.GET("/*", echo.WrapHandler(http.StripPrefix("/", assetHandler)))

	ec.POST("/install", func(c echo.Context) error {
		formData := new(FormData)
		if err := c.Bind(formData); err != nil {
			return err
		}

		// Process the form data as necessary
		cloudConfig := formData.CloudConfig
		reboot := formData.Reboot
		powerOff := formData.PowerOff
		installationDevice := formData.InstallationDevice

		fmt.Println(cloudConfig, reboot, powerOff, installationDevice)

		return c.String(http.StatusOK, "Form data received!")
	})

	if err := ec.Start(l); err != nil && err != http.ErrServerClosed {
		return err
	}

	go func() {
		<-ctx.Done()
		ct, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		ec.Shutdown(ct)
		cancel()
	}()

	return nil
}
