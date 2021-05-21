package main

import (
  "encoding/json"
  "image"
  "image/png"
  "image/color"
  "log"
  "os"
  "fmt"
  // "io"
  // "html"

  "bytes"
  // "flag"
  // "fmt"
  // "net/http"
  // "net/http/cookiejar"
  // "net/url"
  // "os"
  // "strings"

  tf "github.com/galeone/tensorflow/tensorflow/go"
  
  "github.com/otiai10/gosseract/v2"
  
  "github.com/srwiley/oksvg"
  "github.com/srwiley/rasterx"
)

import _ "image/png"    //register PNG decoder

type CaptchaSVG struct {
  Captcha string `json:"captcha"`
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
    log.Printf("Error while parsing:%s",err.Error())
    return
  }
  svgFile, err := os.Create("captcha.svg")

  if err != nil {
     fmt.Println(err)
     return
  }

  defer svgFile.Close()
  // log.Print("Captcha Parsed:\n",ctcha.Captcha)
  svgBytes := []byte(ctcha.Captcha)
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
  // for y := 0; y < b.Max.Y; y++ {
  //   for x := 0; x < b.Max.X; x++ {
  //     // oldPixel := img.At(x, y)
  //     // r, g, b, _ := oldPixel.RGBA()
  //     // y := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
  //     // pixel := color.Gray{uint8(y / 256)}
  //     // imgSet.Set(x, y, pixel)
  //     imgSet.Set(x, y, color.Gray16Model.Convert(img.At(x, y)))
  //    }
  //   }
  icon.Draw(rasterx.NewDasher(w, h, rasterx.NewScannerGV(w, h, rgba, rgba.Bounds())), 2)

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

func getPngImage(pngFileName string) (image.Image, string, error) {
  f, err := os.Open(pngFileName)
  if err != nil {
    log.Printf("Error while loading PNG image:%s",err.Error())
  }
  defer f.Close()

  img, fmtName, err := image.Decode(f)
  if err != nil {
    log.Printf("Error while decoding PNG image:%s",err.Error())
  }
  
  // `fmtName` contains the name used during format registration
  // Work with `img` ...
  return img, fmtName, err
}

func convertToGrayScale(img image.Image) (image.Image, string, error) {
b := img.Bounds()
imgSet := image.NewRGBA(b)
for y := 0; y < b.Max.Y; y++ {
    for x := 0; x < b.Max.X; x++ {
      // oldPixel := img.At(x, y)
      // r, g, b, _ := oldPixel.RGBA()
      // lum := (19595*r + 38470*g + 7471*b + 1<<15) >> 24
      // imgSet.Set(x, y, color.Gray{uint8(lum)})
      imgSet.Set(x, y, color.Gray16Model.Convert(img.At(x, y)))
     }
    }

    outFile, err := os.Create("captcha_grayscale.png")
    if err != nil {
      log.Fatal(err)
    }
    defer outFile.Close()
    png.Encode(outFile, imgSet)
    return getPngImage("captcha_grayscale.png")
}

func breakCaptchaTensorflow(img image.Image) string {
  // img_bytes := buf.Bytes()
// load tensorflow model
  savedModel, err := tf.LoadSavedModel("./tensorflow_savedmodel_captcha", []string{"serve"}, nil)
  if err != nil {
    log.Println("failed to load model", err)
    return ""
  }
  buf := new(bytes.Buffer)
  err = png.Encode(buf, img)

  // run captcha through tensorflow model
  feedsOutput := tf.Output{
    Op:    savedModel.Graph.Operation("CAPTCHA/input_image_as_bytes"),
    Index: 0,
  }
  feedsTensor, err := tf.NewTensor(string(buf.String()))
  if err != nil {
    log.Fatal(err)
    return ""
  }
  feeds := map[tf.Output]*tf.Tensor{feedsOutput: feedsTensor}

  fetches := []tf.Output{
    {
      Op:    savedModel.Graph.Operation("CAPTCHA/prediction"),
      Index: 0,
    },
  }

  captchaText, err := savedModel.Session.Run(feeds, fetches, nil)
  if err != nil {
    log.Fatal(err)
    return ""
  }
  captchaString := captchaText[0].Value().(string)
  return captchaString
}

func resolveCaptchaTesseract(fileName string) {
  client := gosseract.NewClient()
  defer client.Close()
  client.SetImage(fileName)
  text, _ := client.Text()
  fmt.Println("CAPTCHA from tesseract:",text)
  // Hello, World!
}