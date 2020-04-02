# Cloud Reference

There are some fields in the config (suite) that only make sense when used for APId Cloud. They do not affect the
local execution for a config. However, some of them are required for uploading a suite to the cloud.

## Schedule

`schedule` is a root-level field that specifies how often to run the suite in the cloud. When loading
configs (suites) from a directory and multiple files contain `schedule`, it's not deterministic which
schedule will be used; it's not a problem if they are the same.

| Field    | Type   | Required          | Description                                                                   |
| :------  | :----- | :-------          | :-------------------------------------------------                            |
| schedule | string | no; yes for cloud | A [valid cron expression](https://en.wikipedia.org/wiki/Cron#CRON_expression) |
 
We support `?` `/` `*` `,` `-` as well as the standard macros  `@yearly` `@annually` `@monthly` `@weekly` `@daily` `@midnight` `@hourly` `@every <duration>` (`<duration>` is a sequence of numbers with time units: "ns", "us", "ms", "s", "m", "h") |

All times will be interpreted in GMT+0.

```yaml
schedule: "@every 1h20m"
```

```yaml
schedule: "0 0 * * *"
```

You may use [crontab.guru](https://crontab.guru/) for more examples.

## Locations

`locations` is a root-level fields that specifies from where to run the tests in APId Cloud. When loading
configs (suites) from a directory and multiple files contain `locations`, it's not deterministic which
set of locations will be used. It's not a problem if they are the same.

| Field    | Type     | Required          | Description                                              |
| :------  | :-----   | :-------          | :-------------------------------------------------       |
| schedule | []string | no; yes for cloud | A list of locations. Valid elements are to be announced. |

```yaml
locations: ['us-east', 'europe-london']
```

```yaml
locations: 
  - 'us-east'
  - 'europe-london'
```
