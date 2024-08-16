# EchoFunction | Javascript

### Set your environment via Environment
```bash
export SCP_SERVER=https://cloud.synadia.com && \
export SCP_TOKEN=uat_<YOUR_KEY_FROM_SYNADIA_CLOUD_UI> && \
export SCP_ACCOUNT=<YOUR_SYNADIA_CLOUD_ACCOUNT_ID>
```
You will have to include the `--no-context` to all further `synctl` commands to use the environment

### Set your `synctl` context
```bash
synctl context add synadia-test-user <SC_ACCT_ID> \
  --server https://cloud.nas-test.synadia.io \
  --token uat_.....

synctl context select --name synadia-test-user
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
synctl app versions --name echofunctionjs
```

### Deploy your application
```bash
synctl app deploy --name echofunctionjs --version 0.0.1
```
To verify deployment, run the `versions` command again. The deployed version should be highlighted green
```bash
synctl app versions --name echofunctionjs
```

> This step can also be completed via the UI under the `Applications` tab

### Interact with your application
Using a context that is associated with your Synadia Cloud account, you can interact with your application using the `nats` CLI tool.
```bash
nats req my.echotrigger "Hello, World"
```

Your response should be:
```bash
└─❯ nats req "Hello, World"
10:45:31 Sending request on "my.echotrigger"
10:45:31 Received with rtt 157.497898ms
{"triggered_on":"my.echotrigger","payload":"Hello, World"}
```

### Undeploy your application
```bash
synctl app undeploy --name echofunctionjs
```
