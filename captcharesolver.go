package main

import (
  "encoding/json"
  "image"
  "image/png"
  "log"
  "os"
  "fmt"
  // "io"
  // "html"

  "github.com/srwiley/oksvg"
  "github.com/srwiley/rasterx"
)

type CaptchaSVG struct {
  captcha string `json:"captcha"`
}

func getCaptchaSVG() {
  log.Print("Generating SVG CAPTCHA")
  postBody := map[string]interface{}{}
  response, err := queryServer(captchaURLFormat, "POST", postBody)
  
  if err != nil {
     fmt.Println(err)
     return
  }

  ctcha := CaptchaSVG{}
  if err := json.Unmarshal(response, &ctcha); err != nil {
    return
  }
  svgFile, err := os.Create("captcha.svg")

  if err != nil {
     fmt.Println(err)
     return
  }

  defer svgFile.Close()
  log.Print("Captcha Parsed:\n",ctcha.captcha)
  svgBytes := []byte(ctcha.captcha)
  _ , err = svgFile.Write(svgBytes)
  if err != nil {
     fmt.Println(err)
     return
  }
  // save response body into a file
  // io.Copy(svgFile, html.UnescapeString(string(response)))

  fmt.Println("SVG data saved into captcha.svg")
}

func exportToPng(svgFileName string) {
  w, h := 512, 512

  in, err := os.Open(svgFileName)
  if err != nil {
    panic(err)
  }
  defer in.Close()

  icon, _ := oksvg.ReadIconStream(in)
  w = int(icon.ViewBox.W) 
  h = int(icon.ViewBox.H)
  icon.SetTarget(0, 0, float64(w), float64(h))
  rgba := image.NewRGBA(image.Rect(0, 0, w, h))
  icon.Draw(rasterx.NewDasher(w, h, rasterx.NewScannerGV(w, h, rgba, rgba.Bounds())), 1)

  out, err := os.Create("captcha.png")
  if err != nil {
    panic(err)
  }
  defer out.Close()

  err = png.Encode(out, rgba)
  if err != nil {
    panic(err)
  }
}