# Godocit

Create documentation needed issues across different repositories

 [![state](https://img.shields.io/badge/state-unstable-blue.svg)]() [![release](https://img.shields.io/github/release/okkur/actions-godocit.svg)](https://godocit.okkur.org/releases) [![license](https://img.shields.io/github/license/okkur/actions-godocit.svg)](LICENSE)

**NOTE: This is a work-in-progress, we do not consider it production ready. Use at your own risk.**

When a PR introduces a new feature that requires documentation, this action
can automatically create an issue on your website's repository to add
documentation for that new feature.

## Using Godocit
For using godocit you first need to create a Github App and install it on your
user or organization that has **write** access to the target repository.
For more information read the ["Authenticating with GitHub Apps"](https://developer.github.com/apps/building-github-apps/authenticating-with-github-apps/) page.

Note: The generated private key for the GitHub Action should be stored in a 
secret named `PrivateKey` on the source repository. For more information
about secrets, read the [Creating and using encrypted secrets](https://help.github.com/en/actions/automating-your-workflow-with-github-actions/creating-and-using-encrypted-secrets) page.


After creating the app and installing it on your user or organization, write the
following config in `.github/workflows/godocit.yml` file in your source repository:
```
name: CI

on: 
  pull_request:
    types: [labeled]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - name: Create the documentation issues
      uses: okkur/actions-godocit@v0.1.0
      with:
        targetRepo: 'USER/REPO'
      env:
        PRIVATE_KEY: ${{ secrets.PrivateKey }}
        INSTALLATION_ID: ${{ secrets.InstallationID }}

```

Replace the `targetRepo` input with the target repository's name. 
Then after activating the action on the source repository, create the
`PrivateKey` and `InstallationID` secrets on the source repository and fill 
them with the GitHub App's **Private Key** and **Installation ID**.

Now that GitHub Action, workflow and the secrets are setup, you can add the
"documentation needed" label. Then a job will get
activated and creates a "Documentation Needed" issue on the target repository.


## Support
For detailed information on support options see our [support guide](/SUPPORT.md).

## Helping out
Best place to start is our [contribution guide](/CONTRIBUTING.md).

----

*Code is licensed under the [Apache License, Version 2.0](/LICENSE).*  
*Documentation/examples are licensed under [Creative Commons BY-SA 4.0](/docs/LICENSE).*  
*Illustrations, trademarks and third-party resources are owned by their respective party and are subject to different licensing.*

---

Copyright 2019 - The Godocit authors
