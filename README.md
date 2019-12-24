# APId - CLI

Powerful declarative end-to-end testing for APIs that works for you! No coding required. Simple to run on any continuous integration tool.

## Documentation
You can find all APId documentation [here](https://www.getapid.com/docs/).

## Instalation
1. Head to our [downloads page](https://www.getapid.com/download/) and select the download link for your operating system.
2. Once downloaded open the archive and place the executable in a durectory on your PATH.

## First steps

### Using APId
```bash
apid check --config path/to/apid_config.yml
```

### Generating shell completion
Currently `apid` can generate shell completion for `bash`, `zsh` and `powershell`
```bash
apid completion bash -f /etc/bash.completion.d/apid.sh
```

## Development
All useful development commands can be found in the Makefile. Follow these simple steps to build and test the CLI locally:

### 1. Install Mockgen
Mockgen is used to generate mock implementations for testing
```bash
go get github.com/golang/mock/mockgen
```

### 2. Building and running CLI tests
```bash
make
```

### 3. Running CLI end-to-end tests
```bash
make e2e
```

### 4. Running the site locally

```bash
cd site
npm install
zola serve
```

## Contributing
To contribute to APId, please see [CONTRIBUTING](CONTRIBUTING.md).

For questions and discussion join our [FAQ page](https://faq.getapid.com).
