# Link Checker 6000

> **_NOTE:_**  This was just a fun program to write. Use at your own discretion.

This crawls live websites and reports links that are broken. The crawler can 
be configured to not leave a specific website using the `allowedDomains` 
configuration option.

## Usage

```bash
Usage of link-checker-6000:
        link-checker-6000 [flags] [starting url]
Flags:
  -baseURL string
        the URL to use when rewriting relative links
  -config string
        path to a config.yaml file
```

## Config

The following configurations are available in `config.yaml`:

| Key            | Type    | Description                                                                 |
|----------------|---------|-----------------------------------------------------------------------------|
| workerPool     | Integer | Number of threads that will concurrently fetch HTML.                        |
| maxDepth       | Integer | The maximum number of levels this will crawl.                               |
| timeout        | Integer | The timeout before the client errors.                                       |
| allowedDomains | Array   | This configuration prevents the crawler from leaving the configured domains.|
| deniedDomains  | Array   | Prevents the crawler from fetching HTML from the configured domains.        |
