# EchoFunction | Javascript

### Install the `synctl` binary

### Set your environment
```bash
export SCP_SERVER=https://cloud.synadia.com && \
export SCP_TOKEN=uat_<YOUR_KEY_FROM_SYNADIA_CLOUD_UI> && \
export SCP_ACCOUNT=<YOUR_SYNADIA_CLOUD_ACCOUNT_ID>
```

### Upload your artifact
```bash
synctl app art put -t v8 -f ./echofunction.js -n echofunction
```
To verify the upload
```bash
synctl app art ls
```

### Upload your appspec
```bash
synctl app put --file ./appspec.json
```
To verify the upload
```bash
synctl app ls && \
synctl app versions --name echofunction_js
```

### Deploy your application
```bash
synctl app deploy --name echofunction_js --version 0.0.1
```
To verify deployment, run the `versions` command again. The deployed version should be highlighted
```bash
synctl app versions --name echofunction_js
```

> This step can also be completed via the UI under the `Applications` tab

### Interact with your application
Using a context that is associated with your Synadia Cloud account, you can interact with your application using the `nats` CLI tool.
```bash
nats req my.echotrigger "Hello, World!"
```

Your response should be:
```bash
replace
```
