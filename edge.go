package main

import (
	"image"
	"image/color"
	"math"
)

var (
    Ly = [][]int{
        {1, 0, -1},
        {2, 0, -2},
        {1, 0, -1},
    }
    Lx = [][]int{
        {1, 2, 1},
        {0, 0, 0},
        {-1, -2, -1},
    }
);

func edgeDetection(src image.Image) *image.Gray {
    rgbaImg := imageToRGBA(src);
    b := rgbaImg.Bounds();
    grayImg := image.NewGray(b);

    for y := b.Min.Y; y < b.Max.Y; y++ {
        for x := b.Min.X; x < b.Max.X; x++ {
            oldPixel := rgbaImg.At(x, y)
            pixel := color.GrayModel.Convert(oldPixel)
            grayImg.Set(x, y, pixel)
        }
    }

    var (
        pgray uint8
        grayVal float64
        lyacc, lxacc, weight int
        index int
    );

    outImg := image.NewGray(image.Rect(1, 1, b.Dx() - 1, b.Dy() - 1));
    ob := outImg.Bounds();
    for y := ob.Min.Y; y < ob.Max.Y; y++ {
        for x := ob.Min.X; x < ob.Max.X; x++ {
           index = outImg.PixOffset(x, y);

           lyacc = 0;
           for ky := -1; ky < 2; ky++ {
                for kx := -1; kx < 2; kx++ {
                    weight = Ly[ky + 1][kx + 1];
                    pgray = grayImg.Pix[grayImg.PixOffset(x + kx, y + ky)];

                    lyacc += (int(pgray) * weight);
                }
           }

           lxacc = 0;
           for ky := -1; ky < 2; ky++ {
                for kx := -1; kx < 2; kx++ {
                    weight = Lx[ky + 1][kx + 1];
                    pgray = grayImg.Pix[grayImg.PixOffset(x + kx, y + ky)];

                    lxacc += (int(pgray) * weight);
                }
           }

           grayVal = math.Sqrt(float64(lyacc*lyacc) + float64(lxacc*lxacc));
           outImg.Pix[index] = uint8(grayVal);
        }
    }

    return outImg;
}
