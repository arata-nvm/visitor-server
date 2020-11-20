package main

import (
	"bytes"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"image/png"
	"mime"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", handle)

	e.Logger.Fatal(e.Start(":" + port))
}

func handle(c echo.Context) error {
	imgData, err := generateImage()
	if err != nil {
		return err
	}
	c.Response().Header().Set("Cache-Control", "no-cache")
	return c.Blob(http.StatusOK, mime.TypeByExtension(".png"), imgData)
}

var count = 1

func generateImage() ([]byte, error) {
	text := fmt.Sprintf("You are %dth vistor.", count)
	count++

	const S = 1024
	dc := gg.NewContext(S, S)
	dc.SetRGB(0, 0, 0)
	if err := dc.LoadFontFace("./Roboto-Regular.ttf", 96); err != nil {
		return nil, err
	}
	dc.DrawStringAnchored(text, S/2, S/2, 0.5, 0.5)

	var buf bytes.Buffer
	if err := png.Encode(&buf, dc.Image()); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
