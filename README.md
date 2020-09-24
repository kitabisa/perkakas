# perkakas

perkakas/per·ka·kas/ - segala yang dapat dipakai sebagai alat (seperti untuk makan, bekerja di dapur, perang, ngoding)

Library for supporting common backend tasks. If you have function that considered common, please move it here.

How to develop on perkakas:

1. Any common function that **does not have any business logic** can be moved here.
2. You should group your function into a folder as a package.
3. If the folder not exist, just create it. If the folder exist, put your code there.

## Installation

```bash
$ go get github.com/kitabisa/perkakas/v2
```

## How To Use

See `README.md` in each folders to see how to use these modules.

### Logger

nb: logger from v2.14.6 is deprecated.
Perkakas logger is based on Zerolog since version v2.15.0
If you want to use logging from perkakas, follow this:

1.  use middleware: `RequestIDToContextAndLogMiddleware`
1.  then call the logger: `Zlogger(context)`, e.g.:

        Zlogger(context).Err(err).Msg("your-message")
        Zlogger(context).Info().Msg("your-message")
        Zlogger(context).Err(err).Send()

        //see zerolog github for more usage

There's also another middleware that logging the request header & body, only when error happen: `RequestLogger`
