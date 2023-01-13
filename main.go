package main

import (
	"errors"
	"fmt"
	"os"
)

func root(args []string) error {
    if len(args) < 1 {
        return errors.New("a subcommand is required");
    }

    cmds := []Runner{
        NewResizeCommand(),
        NewBlurCommand(),
        NewBrickCommand(),
    }

    subcommand := args[0];

    for _, cmd := range cmds {
        if cmd.Name() == subcommand {
            cmd.Init(args[1:])
            return cmd.Run();
        }
    }

    return fmt.Errorf("Unknown subcommand: %s", subcommand)
}

func main() {
    if err := root(os.Args[1:]); err != nil {
        fmt.Println(err);
        os.Exit(1);
    }
}
