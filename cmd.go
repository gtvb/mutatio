package main

import (
	"errors"
	"flag"
	"image"

	"golang.org/x/image/draw"
)

type Runner interface {
    Init(args []string) 
    Name() string
    Run() error
}

type ResizeCommand struct {
    fs *flag.FlagSet;

    in, out string
    width, height int;
}

func NewResizeCommand() *ResizeCommand {
    rc := &ResizeCommand{
        fs: flag.NewFlagSet("resize", flag.PanicOnError),
    }

    rc.fs.StringVar(&rc.in, "in", "", "source image path");
    rc.fs.StringVar(&rc.out, "out", "out.jpg", "destination image path");
    rc.fs.IntVar(&rc.width, "w", 0, "new image width");
    rc.fs.IntVar(&rc.height, "h", 0, "new image height");

    return rc;
}

func (rc *ResizeCommand) Init(args []string) {
    rc.fs.Parse(args);
}

func (rc *ResizeCommand) Name() string {
    return rc.fs.Name();
}

func (rc *ResizeCommand) Run() error {
    if rc.in == "" {
        return errors.New("no source image path provided");
    }

    if rc.width == 0 || rc.height == 0 {
        return errors.New("no destination image width or height provided");
    }

    src, err := LoadImage(rc.in);
    if err != nil {
        return err
    }

    dst := ResizeImage(src, rc.width, rc.height);

    if err = SaveJPG(rc.out, dst, 100); err != nil {
        return err
    }

    return nil;
}

type BlurCommand struct {
    fs *flag.FlagSet;

    in, out string
    radius int
}

func NewBlurCommand() *BlurCommand {
    bc := &BlurCommand{
        fs: flag.NewFlagSet("blur", flag.PanicOnError),
    }

    bc.fs.StringVar(&bc.in, "in", "", "source image path");
    bc.fs.StringVar(&bc.out, "out", "out.jpg", "destination image path");
    bc.fs.IntVar(&bc.radius, "radius", 1, "kernel radius");

    return bc;
}

func (bc *BlurCommand) Init(args []string) {
    bc.fs.Parse(args);
}

func (bc *BlurCommand) Name() string {
    return bc.fs.Name();
}

func (bc *BlurCommand) Run() error {
    if bc.in == "" {
        return errors.New("no source image path provided");
    }

    if bc.radius <= 0 {
        return errors.New("radius must be bigger than or equal to one");
    }

    src, err := LoadImage(bc.in);
    if err != nil {
        return err
    }

    dst := BoxBlur(src, bc.radius);

    if err = SaveJPG(bc.out, dst, 100); err != nil {
        return err
    }

    return nil;
}

type BrickCommand struct {
    fs *flag.FlagSet;

    legoIn string
    in, out string
}

func NewBrickCommand() *BrickCommand {
    bc := &BrickCommand{
        fs: flag.NewFlagSet("brick", flag.PanicOnError),
    }

    // TODO: if the user provides a lego brick on its own,
    // we need to resize it, so, remenber to make the resize operation 
    // a util.
    bc.fs.StringVar(&bc.in, "in", "", "source image path");
    bc.fs.StringVar(&bc.legoIn, "lin", "images/lego-25.jpg", "lego brick image path");
    bc.fs.StringVar(&bc.out, "out", "out.jpg", "destination image path");

    return bc;
}

func (bc *BrickCommand) Init(args []string) {
    bc.fs.Parse(args);
}

func (bc *BrickCommand) Name() string {
    return bc.fs.Name();
}

func (bc *BrickCommand) Run() error {
    if bc.in == "" {
        return errors.New("no source image path provided");
    }

    src, err := LoadImage(bc.in);
    if err != nil {
        return err
    }
    srcb := src.Bounds();

    lego, err := LoadImage(bc.legoIn);
    if err != nil {
        return err
    }
    legob := lego.Bounds();

    remainderW := srcb.Max.X % legob.Max.X;
    remainderH := srcb.Max.Y % legob.Max.Y;
    if remainderW != 0 && remainderH != 0 {
        src = ResizeImage(src, srcb.Dx() - remainderW, srcb.Dy() - remainderH)     
    } else if remainderW != 0 {
        src = ResizeImage(src, srcb.Dx() - remainderW, srcb.Dy())     
    } else if remainderH != 0 {
        src = ResizeImage(src, srcb.Dx(), srcb.Dy() - remainderH)     
    }

    srcb = src.Bounds();
    brickLayer := image.NewRGBA(srcb);
    var r image.Rectangle;
    for y := srcb.Min.Y; y < srcb.Max.Y; y += legob.Dy() {
        for x := srcb.Min.X; x < srcb.Max.X; x += legob.Dx() {
            r = image.Rect(x, y, x + legob.Dx(), y + legob.Dy());
            draw.Draw(brickLayer, r, lego, image.Pt(0, 0), draw.Src);
        }
    }

    outImg := Blend(src, brickLayer, OverlayFunc)
    if err = SaveJPG(bc.out, outImg, 100); err != nil {
        return err;
    }

    return nil;
}
