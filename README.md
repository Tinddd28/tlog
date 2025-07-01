## Tlog
My view on the used logger.
This logger based on the log/slog package.
All logs are saved to a file in the specified dir. The recording takes place in a file corresponding to the current date, and the file changes on the following day. Log files are stored for the specified time - savingDays.

## Usage
```go
package main

import (
    "log/slog"
    "github.com/Tinddd28/tlog"
)

func main() {
    level := "debug" 
    dir := "./logs"
    format := ".log"
    savingDays := 7

    logOpts := tlog.NewLogOpts(level, dir, format, savingDays)

    logger, err := tlog.SetupLogger(slogOpts)
    if err != nil {
        os.Exit(1)
    }

    logger.Info("Hello!", "message level", level)
    /*
       [DD:MM:YY HH:MM:SS] INFO Hello!: 
       {
            "message level": "debug"
       }
    */

    lg := logger.With(
        slog.String("Point", "some point")
    )

    lg.Warn("hello!", "message of warn", "some msg")
    /*
        [DD:MM:YY HH:MM:SS] WARN hello!:
        {
            "Point": "some point",
            "message of warn": "some msg"
        }

    */
}

```
