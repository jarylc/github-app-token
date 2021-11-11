# GitHub App Token Generator
Generate a GitHub App integration access token to use for subsequent steps

## Inputs
### `app`
**Required** - App name or ID
### `key`
**Required** - Base64 encoded private key of the GitHub App
### `owner`
**Optional** - owner name for which the integration token will be requested (defaults to current)

## Outputs
### `token`
**Masked** - An installation access token for the GitHub App on the requested organization

## Example usage
```yml
jobs:
  job:
    runs-on: ubuntu-latest
    steps:
      - name: Generate Token
        id: generate_token
        uses: jarylc/github-app-token@v3.0.1
        with:
          app: ${{ secrets.GH_APP_NAME }}
          key: ${{ secrets.GH_APP_KEY }}
          # Optional (defaults to the current owner)
          # owner: jarylc
      - name: Use Token
        env:
          TOKEN: ${{ steps.generate_token.outputs.token }}
        run: |
          echo "The generated token is masked: ${TOKEN}"
```
