# Scheduled execution

Scheduled execution works pretty much the same way. You need to define your transactions, as you would with the CLI, and then using the APId CLI upload them to APId cloud. This works as follows:

- The CLI reads the transactions from the specified directory or file, as per ususal\
- The CLI will marshal the transactions and upload them to the APId cloud with the provided suite name. A suite is a bundle of transactions that will be executed together.
- APId cloud will read the configuration you've uploaded and will execute the transactions from each of the specified locations as per the schedule.

There are some caveats in doing this, mainly that shell commands DO NOT work, since, well, it's a remote machine.

## Usage

In order to use the power of the cloud you will need a personal access key. To generate one, you will have to:

- Head over to https://console.getapid.com and sign up
- Go to the settings page and create a new access key

Once you have your key you will need to [install the APId CLI](../installation/cli.md) (if you haven't already) or use our [official docker image](../installation/docker.md).

A reference on how use the CLI after installation for uploading suites can be found [here](../cli/cloud/upload.md).
