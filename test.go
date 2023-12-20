package main

import "fmt"
import "os"
import "image"
import "math"
import "image/color"
import  "image/draw"
import "image/jpeg"
import _"image/png"


func main() {
    test2();
    return ;
}


func test2(){
    // 打开图像文件
    file, err := os.Open("0209.jpg")
    if err!= nil {
            fmt.Println("Error opening image file:", err)
            return
    }
    defer file.Close()
    // 从文件创建一个 image.Image 对象
    img, _, err := image.Decode(file)
    if err!= nil {
            fmt.Println("Error decoding image:", err)
            return
    }
    
      // 打开图像文件
      file1, err := os.Open("0201.jpg")
      if err!= nil {
              fmt.Println("Error opening image file:", err)
              return
      }
      defer file1.Close()
      // 从文件创建一个 image.Image 对象
      tpl, _, err := image.Decode(file1)
      if err!= nil {
              fmt.Println("Error decoding image:", err)
              return
      }


    x,y:=   matchTemplate(img,tpl);
    fmt.Printf("x:%d y:%d\r\n",x,y);

    x1,y1:=   matchTemplate1(img,tpl);
    fmt.Printf("x:%d y:%d\r\n",x1,y1);



  // 创建一个新的RGBA图像，并将原图像复制到新图像中
    b := img.Bounds()
    imgRGBA := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
    draw.Draw(imgRGBA, img.Bounds(), img, b.Min, draw.Src)

    // 创建一个红色的画笔
    red := color.RGBA{255, 0, 0, 255}

    // 在新图像上绘制矩形
    rect := image.Rect(x1, y1, x1+tpl.Bounds().Size().X, y1+tpl.Bounds().Size().Y) // 这是一个例子，你可以根据需要改变矩形的大小和位置
    // 在新图像上绘制矩形的边框
    draw.Draw(imgRGBA, image.Rect(rect.Min.X, rect.Min.Y, rect.Max.X, rect.Min.Y+1), &image.Uniform{red}, image.ZP, draw.Src) // Top line
    draw.Draw(imgRGBA, image.Rect(rect.Min.X, rect.Max.Y-1, rect.Max.X, rect.Max.Y), &image.Uniform{red}, image.ZP, draw.Src) // Bottom line
    draw.Draw(imgRGBA, image.Rect(rect.Min.X, rect.Min.Y, rect.Min.X+1, rect.Max.Y), &image.Uniform{red}, image.ZP, draw.Src) // Left line
    draw.Draw(imgRGBA, image.Rect(rect.Max.X-1, rect.Min.Y, rect.Max.X, rect.Max.Y), &image.Uniform{red}, image.ZP, draw.Src) // Right line


    // 创建一个新文件，并将修改后的图像保存为PNG格式
    outFile, err := os.Create("output.png")
    if err != nil {
        panic(err)
    }
    defer outFile.Close()

    err = jpeg.Encode(outFile, imgRGBA, nil)
    if err != nil {
        panic(err)
    }

}

func toGrayscale(img image.Image) *image.Gray {
    bounds := img.Bounds()
    grayImg := image.NewGray(bounds)
    for y := 0; y < bounds.Max.Y; y++ {
        for x := 0; x < bounds.Max.X; x++ {
            oldColor := img.At(x, y)
            grayColor := color.GrayModel.Convert(oldColor)
            grayImg.Set(x, y, grayColor)
        }
    }
    return grayImg
}


func matchTemplate(img image.Image, tpl image.Image) (bestX, bestY int) {
    grayImg := toGrayscale(img)
    grayTpl := toGrayscale(tpl)

    tplBounds := grayTpl.Bounds()
    width, height := tplBounds.Size().X, tplBounds.Size().Y
    minDiff := math.MaxFloat64

    for y := 0; y <= grayImg.Bounds().Size().Y-height; y++ {
        for x := 0; x <= grayImg.Bounds().Size().X-width; x++ {
            var sumSquareDiff float64
            for j := 0; j < height; j++ {
                for i := 0; i < width; i++ {
                    v1, _, _, _ := grayImg.At(x+i, y+j).RGBA()
                    v2, _, _, _ := grayTpl.At(i, j).RGBA()
                    diff := float64(v1-v2) * float64(v1-v2)
                    sumSquareDiff += diff
                }
            }
            if sumSquareDiff < minDiff {
                minDiff = sumSquareDiff
                bestX = x
                bestY = y
            }
        }
    }
    return bestX, bestY
}


func matchTemplate1(img image.Image, tpl image.Image) (bestX, bestY int) {
    grayImg := toGrayscale(img)
    grayTpl := toGrayscale(tpl)
    scores := crossCorrelationNormed(grayImg, grayTpl)
    bestScore := -math.MaxFloat64
    bestX = -1
    bestY = -1
    dx := img.Bounds().Dx() - tpl.Bounds().Dx() + 1
    dy := img.Bounds().Dy() - tpl.Bounds().Dy() + 1
    for y := 0; y < dy; y++ {
        for x := 0; x < dx; x++ {
            score := scores[y*dx+x]
            if score > bestScore {
                bestScore = score
                bestX = x
                bestY = y
            }
        }
    }
    if bestX != -1 && bestY != -1 {
        bestX += tpl.Bounds().Min.X
        bestY += tpl.Bounds().Min.Y
    }
    return bestX, bestY
}

func crossCorrelationNormed(src, tpl *image.Gray) []float64 {
    sumTpl := 0.0
    sumSqTpl := 0.0
    for y := tpl.Rect.Min.Y; y < tpl.Rect.Max.Y; y++ {
        for x := tpl.Rect.Min.X; x < tpl.Rect.Max.X; x++ {
            val := float64(tpl.GrayAt(x, y).Y)
            sumTpl += val
            sumSqTpl += val * val
        }
    }
    meanTpl := sumTpl / float64(tpl.Rect.Dx()*tpl.Rect.Dy())
    stddevTpl := math.Sqrt(sumSqTpl/float64(tpl.Rect.Dx()*tpl.Rect.Dy()) - meanTpl*meanTpl)
    scores := make([]float64, (src.Rect.Dx()-tpl.Rect.Dx()+1)*(src.Rect.Dy()-tpl.Rect.Dy()+1))
    i := 0
    for y := src.Rect.Min.Y; y <= src.Rect.Max.Y-tpl.Rect.Dy(); y++ {
        for x := src.Rect.Min.X; x <= src.Rect.Max.X-tpl.Rect.Dx(); x++ {
            sumSrc := 0.0
            sumSqSrc := 0.0
            sumTplSrc := 0.0
            for ty := tpl.Rect.Min.Y; ty < tpl.Rect.Max.Y; ty++ {
                for tx := tpl.Rect.Min.X; tx < tpl.Rect.Max.X; tx++ {
                    srcVal := float64(src.GrayAt(x+tx, y+ty).Y)
                    tplVal := float64(tpl.GrayAt(tx, ty).Y)
                    sumSrc += srcVal
                    sumSqSrc += srcVal * srcVal
                    sumTplSrc += srcVal * tplVal
                }
            }
            meanSrc := sumSrc / float64(tpl.Rect.Dx()*tpl.Rect.Dy())
            stddevSrc := math.Sqrt(sumSqSrc/float64(tpl.Rect.Dx()*tpl.Rect.Dy()) - meanSrc*meanSrc)
            if stddevSrc > 1e-8 {
                scores[i] = (sumTplSrc/float64(tpl.Rect.Dx()*tpl.Rect.Dy()) - meanSrc*meanTpl) / (stddevSrc * stddevTpl)
            }
            i++
        }
    }
    return scores
}
