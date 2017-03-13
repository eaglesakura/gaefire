# GAE/Fire

GAE/FireはGAE/Goで Firebase Service Account の機能を利用するためのライブラリです。

# 機能

GAE/FireはFirebaseとの連携をサポートするため、次の機能を提供します。

 * [Json Web Token(JWT)](https://jwt.io/) 管理
  - [x] [Firebase Custom Authentication](https://firebase.google.com/docs/auth/) のためのToken Generate / Verify
  - [x] Firebaseが発行したTokenのVerify
  - [x] [Google Id Token](https://developers.google.com/identity/sign-in/web/backend-auth) のVerify
 * OAuth2 Token管理
  - [x] Service Account権限を持ったOAuth2 Token Generate / Refresh / Verify
  - [x] Google Play Service(Android)にて発行された access_code からOAuth2 TokenのGenerate / Refresh / Verify
  - [x] MemcacheによるToken Cache
 * [Cloud Endpoints 2.0相当](https://cloud.google.com/endpoints/)の認証サポート
  - [x] [Service Control API](https://cloud.google.com/service-control/how-to) を経由した [Google Cloud Console API Key](https://console.developers.google.com/apis/dashboard) のVerify
  - [x] `Authorization` HeaderのVerify / Google Cloud Endpoints 2.0互換 `X-Endpoint-API-UserInfo` Header Generate
  - [x] `swagger.json` に基づいたOAuth2 audience verify
  - [x] `swagger.json` に基づいたGoogle Id Token audience verify
  - [ ] TODO: `swagger.json` に基づいたPath単位での認証ON/OFF

# Install

```
# Install GAE/Fire
$ go get -u -f github.com/eaglesakura/gaefire
```

# Samples

## Setup Service Account

## Generate Firebase custom auth token

## Verify Firebase auth token

## Verify Google id token

## Generate Service account oauth2 token

## Generate User oauth2 token

## Refresh oauth2 token

## Cloud Endpoints 2.0 compat

### Setup

### Verify API Key & Authorization header


# LICENSE
