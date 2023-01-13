package main

import (
	"image"
	"math"
)

type pixel struct {
    r, g, b, a uint32;
};

func createPixel(r, g, b, a uint8) pixel {
    return pixel{
        r: uint32(r),
        g: uint32(g),
        b: uint32(b),
        a: uint32(a),
    };
}

type stack []pixel;

func (s *stack) push(p pixel) {
    *s = append(*s, p);
}

func (s *stack) top() pixel {
    return (*s)[len((*s)) - 1];
}

func (s *stack) clear() {
    *s = (*s)[:0];
}

func (s *stack) removeFirst() pixel {
    removed := (*s)[0];
    *s = (*s)[1:];
    return removed;
}

// Box Blur implementation with horizontal and vertical passes. 
// It's a fast algorithm, but it doesn't generate a perfect blur.
// In conjunction with the edge detection algorithm, it does not
// give me the best results.A small radius produces a big blur, and
// that compromises the edge detector a little bit.
func BoxBlur(src image.Image, radius int) *image.RGBA {
    s := new(stack);
    elementsInBox := uint32(2 * radius + 1);

    rgbaImg := imageToRGBA(src);
    b := rgbaImg.Bounds();
    firstPassImg := image.NewRGBA(b);

    var (
        removedPixel, topPixel pixel
        pixelStartIndex int
        avgR, avgG, avgB, avgA uint32 = 0, 0, 0, 0
        pa, pr, pg, pb uint8
    )

    for y := b.Min.Y; y < b.Max.Y; y++ {
        for x := b.Min.X; x < b.Max.X; x++ {
            pixelStartIndex = (y-rgbaImg.Rect.Min.Y)*rgbaImg.Stride + (x-rgbaImg.Rect.Min.X)*4;

            if x == b.Min.X {
                avgR = 0;
                avgG = 0;
                avgB = 0;
                avgA = 0;

                for nx := -radius; nx < radius + 1; nx++ {
                    var xi int;
                    if x + nx < 0 {
                        xi = (x + nx) + (2 * radius);
                    } else {
                        xi = x + nx
                    }

                    index := (y-rgbaImg.Rect.Min.Y)*rgbaImg.Stride + (xi-rgbaImg.Rect.Min.X)*4
                    pr, pg, pb, pa = rgbaImg.Pix[index], rgbaImg.Pix[index + 1], rgbaImg.Pix[index + 2], rgbaImg.Pix[index + 3];

                    s.push(createPixel(pr, pg, pb, pa));
                }

                for _, p := range *s {
                    avgR += p.r;
                    avgG += p.g;
                    avgB += p.b;
                    avgA += p.a;
                }


                avgR /= elementsInBox;
                avgG /= elementsInBox;
                avgB /= elementsInBox;
                avgA /= elementsInBox;

                firstPassImg.Pix[pixelStartIndex] = uint8(avgR);
                firstPassImg.Pix[pixelStartIndex + 1] = uint8(avgG);
                firstPassImg.Pix[pixelStartIndex + 2] = uint8(avgB);
                firstPassImg.Pix[pixelStartIndex + 3] = uint8(avgA);
            } else {
                removedPixel = s.removeFirst();
                avgR -= (removedPixel.r / elementsInBox);
                avgG -= (removedPixel.g / elementsInBox);
                avgB -= (removedPixel.b / elementsInBox);
                avgA -= (removedPixel.a / elementsInBox);

                var nextX int;
                if x + radius >= b.Max.X {
                    nextX = x - radius;
                } else {
                    nextX = x + radius;
                }

                index := (y-rgbaImg.Rect.Min.Y)*rgbaImg.Stride + (nextX-rgbaImg.Rect.Min.X)*4
                pr, pg, pb, pa = rgbaImg.Pix[index], rgbaImg.Pix[index + 1], rgbaImg.Pix[index + 2], rgbaImg.Pix[index + 3];

                s.push(createPixel(pr, pg, pb, pa));
                topPixel = s.top();
                avgR += (topPixel.r / elementsInBox);
                avgG += (topPixel.g / elementsInBox);
                avgB += (topPixel.b / elementsInBox);
                avgA += (topPixel.a / elementsInBox);

                firstPassImg.Pix[pixelStartIndex] = uint8(avgR);
                firstPassImg.Pix[pixelStartIndex + 1] = uint8(avgG);
                firstPassImg.Pix[pixelStartIndex + 2] = uint8(avgB);
                firstPassImg.Pix[pixelStartIndex + 3] = uint8(avgA);
            }
        }
        s.clear();
    }

    outImg := image.NewRGBA(b);

    for x := b.Min.X; x < b.Max.X; x++ {
        for y := b.Min.Y; y < b.Max.Y; y++ {
            pixelStartIndex = (y-firstPassImg.Rect.Min.Y)*firstPassImg.Stride + (x-firstPassImg.Rect.Min.X)*4;

            if y == b.Min.Y {
                avgR = 0;
                avgG = 0;
                avgB = 0;
                avgA = 0;

                for ny := -radius; ny < radius + 1; ny++ {
                    var yi int;
                    if y + ny < 0 {
                        yi = (y + ny) + (2 * radius);
                    } else {
                        yi = y + ny
                    }

                    index := (yi-firstPassImg.Rect.Min.Y)*firstPassImg.Stride + (x-firstPassImg.Rect.Min.X)*4
                    pr, pg, pb, pa = firstPassImg.Pix[index], firstPassImg.Pix[index + 1], firstPassImg.Pix[index + 2], firstPassImg.Pix[index + 3];

                    s.push(createPixel(pr, pg, pb, pa));
                }

                for _, p := range *s {
                    avgR += p.r;
                    avgG += p.g;
                    avgB += p.b;
                    avgA += p.a;
                }


                avgR /= elementsInBox;
                avgG /= elementsInBox;
                avgB /= elementsInBox;
                avgA /= elementsInBox;

                outImg.Pix[pixelStartIndex] = uint8(avgR);
                outImg.Pix[pixelStartIndex + 1] = uint8(avgG);
                outImg.Pix[pixelStartIndex + 2] = uint8(avgB);
                outImg.Pix[pixelStartIndex + 3] = uint8(avgA);
            } else {
                removedPixel = s.removeFirst();
                avgR -= (removedPixel.r / elementsInBox);
                avgG -= (removedPixel.g / elementsInBox);
                avgB -= (removedPixel.b / elementsInBox);
                avgA -= (removedPixel.a / elementsInBox);

                var nextY int;
                if y + radius >= b.Max.Y {
                    nextY = y - radius;
                } else {
                    nextY = y + radius;
                }

                index := (nextY-firstPassImg.Rect.Min.Y)*firstPassImg.Stride + (x-firstPassImg.Rect.Min.X)*4
                pr, pg, pb, pa = firstPassImg.Pix[index], firstPassImg.Pix[index + 1], firstPassImg.Pix[index + 2], firstPassImg.Pix[index + 3];

                s.push(createPixel(pr, pg, pb, pa));
                topPixel = s.top();
                avgR += (topPixel.r / elementsInBox);
                avgG += (topPixel.g / elementsInBox);
                avgB += (topPixel.b / elementsInBox);
                avgA += (topPixel.a / elementsInBox);

                outImg.Pix[pixelStartIndex] = uint8(avgR);
                outImg.Pix[pixelStartIndex + 1] = uint8(avgG);
                outImg.Pix[pixelStartIndex + 2] = uint8(avgB);
                outImg.Pix[pixelStartIndex + 3] = uint8(avgA);
            }
        }
        s.clear();
    }

    return outImg;
}

func GaussianBlur(src image.Image, radius int) *image.RGBA {
    // 1D kernel generation
    kSize := (2 * radius) + 1;
    kernel := make([]float64, kSize);
    stdDev := math.Max(float64(radius / 2), 1);

    var sum float64 = 0;
    for x := -radius; x < radius + 1; x++ {
        exponent := -float64(x*x)/(2*stdDev*stdDev);
        kValue := math.Pow(math.E, exponent);
        kernel[x + radius] = kValue;
        sum += kValue;
    }

    for i := 0; i < len(kernel); i++ {
        kernel[i] /= sum;
    }

    // first pass convolution

    // second pass convolution

    return nil
}
