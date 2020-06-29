# env-to-flags

Convert environment variables to cli flags for a command. Example:

```bash
export CURL_REMOTE_NAME=google.html
export CURL_LOCATION=''
env-to-flags curl google.com
```

## Why Do this?

I wanted to create an AWS Batch job to run a command. I thought it would be nice if the caller could specify arbitrary flags for the commands. However, some of the command comes from user generated input. To avoid any injection attacks while maintaining flexibility I created this tool.
