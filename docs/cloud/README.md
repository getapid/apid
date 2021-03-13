# Cloud

Besides being able to run your tests from your local machine, APId also has the functionality to run from around the world. This functionality is powered by the APId cloud offering.

## Scheduled execution

This is the most powerful mode of APId cloud. It allows you to define a set of checks, upload them and let us run them on a predefined schedule and from predefined locations for you! You'll get a notification as soon as when we detect something went wrong with your API!

- [more information](cloud.md)

## On demand remote execution

The remote execution offering works pretty much the same way. You need to define your transactions, as you would with the CLI, and then using the APId CLI issue a `remote` command.

- [more information](remote.md)

## Regions

APId cloud runs in multiple regions worldwide. Below is a list of the current ones.

| Region Name  | Location      |
| :----------- | :------------ |
| montreal     | Montreal      |
| washington   | Washington    |
| sanfrancisco | San Francisco |
| mumbai       | Mumbai        |
| tokyo        | Tokyo         |
| sydney       | Sydney        |
| dublin       | Dublin        |
| stockholm    | Stockholm     |
| frankfurt    | Frankfurt     |
| saopaulo     | Sao Paulo     |

## Billing

We've tried making the billing model as simple as possible. Each account has a free tier quota of units they can use each month for running their tests on the cloud infrastructure, after which there is a flat fee for each unit used.

Each unit corresponds to 100ms of execution time.

In case of on demand remote execution, you will be billed separately for every step.

When using scheduled execution you're billed for the whole duration of the suite execution.

### Examples

| API response time | Units billed |
| :---------------- | :----------- |
| 23                | 1            |
| 99                | 1            |
| 100               | 1            |
| 105               | 2            |
| 1999              | 20           |
