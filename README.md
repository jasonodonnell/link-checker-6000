# Link Checker 6000

> **_NOTE:_**  This was just a fun program to write. Use at your own discretion.

Checks a website for broken links. This crawls live websites and will 
fetch the HTML for every link it finds. The crawler can be configured to 
not leave a specific website using the `allowedDomains` configuration option.

## Usage

```bash
$ link-checker-6000 -config=/path/to/config.yaml <Starting URL> <Base URL>
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
