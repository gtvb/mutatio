package main

import (
	"image"
	"image/color"
)

type blendFunc func(src, dst color.Color) color.Color;

// blend overlays source with destination
func Blend(src, dst image.Image, bf blendFunc) *image.RGBA {
    srcRgba := imageToRGBA(src); 
    srcBounds := srcRgba.Bounds();

    dstRgba := imageToRGBA(dst); 
    dstBounds := dstRgba.Bounds();

    outImg := image.NewRGBA(dstBounds);
    inter := dstBounds.Intersect(srcBounds);

    for y := dstBounds.Min.Y; y < dstBounds.Max.Y; y++ {
        for x := dstBounds.Min.X; x < dstBounds.Max.X; x++ {
            if p := image.Pt(x, y); p.In(inter) {
                outImg.Set(x, y, bf(srcRgba.At(x, y), dstRgba.At(x, y)));
            } else {
                outImg.Set(x, y, dstRgba.At(x, y));
            }
        }
    }

    return outImg;
}

const (
	max = 65535.0
	mid = max / 2.0
)

func OverlayFunc(src, dst color.Color) color.Color {
    sr, sg, sb, sa := colorToFloat64(src.RGBA()); 
    dr, dg, db, da := colorToFloat64(dst.RGBA()); 

    or := toUint16(overlay(sr, dr));
    og := toUint16(overlay(sg, dg));
    ob := toUint16(overlay(sb, db));
    oa := toUint16(overlay(sa, da));

    return color.RGBA64{or, og, ob, oa};
}

func colorToFloat64(r, g, b, a uint32) (float64, float64, float64, float64) {
    return float64(r), float64(g), float64(b), float64(a);
}

func toUint16(x float64) uint16 {
	if x < 0 {
		return 0
	}
	if x > 65535 {
		return 65535
	}
	return uint16(int(x + 0.5))
}

func overlay(d, s float64) float64 {
    if d < mid {
		return 2 * s * d / max
	}
	return max - 2*(max-s)*(max-d)/max
}

