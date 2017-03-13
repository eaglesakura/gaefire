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

 1. [Google Cloud Platform Console](https://console.developers.google.com)からプロジェクトを作成する
 1. [Firebase Console](https://console.firebase.google.com/) にアクセスし、GCP ProjectとFirebase Projectをリンクする
 1. Firebase Consoleで `Overview` > `プロジェクトの設定` > `サービスアカウント`からFirebase Admin SDKをセットアップ
 1. Firebase Service Accountのメールアドレスを確認する
  * ex) firebase-adminsdk@your-gcp-project-name.iam.gserviceaccount.com
 1. `新しい秘密鍵を生成` でサービスアカウント情報を取得する

```
# service-account.json
{
  "type": "service_account",
  "project_id": "your-gcp-project-name",
  "private_key_id": "fa9....3d",
  "private_key": "-----BEGIN PRIVATE KEY----- ... -----END PRIVATE KEY-----\n",
  "client_email": "firebase-adminsdk@your-gcp-project-name.iam.gserviceaccount.com",
  "client_id": "12345678901234567890",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://accounts.google.com/o/oauth2/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk%40your-gcp-project-name.iam.gserviceaccount.com"
}
```

## Generate Firebase custom auth token

GAE/Fireの`ServiceAccount`オブジェクトはJWTの生成・管理を行います。

生成されたTokenはMobile等の認証に使用することができます。
トークンはFirebase Service Accountにより署名され、3600 secの有効期限があります。

```go
import "github.com/eaglesakura/gaefire"
import "github.com/eaglesakura/gaefire/factory"

func NewCustomAuthToken(newUserId string) string {
  var serviceAccountJson []byte
  serviceAccountJson = ... // load from service-account.json

  serviceAccount := factory.NewServiceAccount(serviceAccountJson)
  token, err :=  serviceAccount.NewFirebaseAuthTokenGenerator(newUserId).Generate()

  if err != nil {
    return ""
  } else {
    return token
  }
}
```

```java
// Login Custom Token(for Android)
token = ... // load from GAE/Go
FirebaseAuth.getInstance().signInWithCustomToken(token);
```

## Verify Firebase auth token

Firebase Custom Authでログインを行った後は、Firebase SDKが自動的にJWTトークン管理を行います。
Firebase SDKが持っているTokenを `Authentication` 等のHeaderに付与することで、ユーザー認証を行えます。

```java
FirebaseUser user = ... // wait login
user.getToken(true);

//$ curl -H "Authentication: Bearer ${token}" https://your-gcp-project-name.appspot.com/path/to/api
```

```go
import "github.com/eaglesakura/gaefire"

func VerifyFirebaseAuthToken(ctx context.Context, token string) bool {
  serviceAccount := ... // NewServiceAccount()

  verifiedToken, err := serviceAccount.NewFirebaseAuthTokenVerifier(ctx, token).Valid()
  if err != nil {
    return false
  } else {
    // check user...
    user := gaefire.FirebaseUser{}
    verifiedToken.GetUser(&user)

    // verified ok
    return true
  }
}
```

## Verify Google id token

[Google Play Service](https://developers.google.com/identity/sign-in/android/start-integrating) 等のGoogleログイン機能を利用して取得したGoogle Id TokenをVerifyし、認証に利用することができます。
Google Id Tokenは3600 secでExpireされます。

```java
// using Google Play Service
public static GoogleApiClient.Builder newFullPermissionClient(Context context) {
    GoogleSignInOptions options = new GoogleSignInOptions.Builder(GoogleSignInOptions.DEFAULT_SIGN_IN)
            .requestIdToken(context.getString(R.string.default_web_client_id))
            .requestEmail()
            .build();
    return new GoogleApiClient.Builder(context)
            .addApi(Auth.GOOGLE_SIGN_IN_API, options)
            ;
}

@OnActivityResult(REQUEST_GOOGLE_AUTH)
void resultGoogleAuth(int result, Intent data) {
    GoogleSignInResult signInResult = Auth.GoogleSignInApi.getSignInResultFromIntent(data);
    if (signInResult.isSuccess()) {
        GoogleSignInAccount signInAccount = signInResult.getSignInAccount() ;

        String googleIdToken = signInAccount.getIdToken();
        // curl -H "Authorization: Bearer {token}"  https://your-gcp-project-name.appspot.com/path/to/api
    } else {
      // login failed.
    }
}
```

```go

import "github.com/eaglesakura/gaefire"

func VerifyGoogleIdToken(ctx context.Context, token string) bool {
  serviceAccount := ... // factory.NewServiceAccount()

  verifiedToken, err := serviceAccount.NewGoogleAuthTokenVerifier(ctx, token).Valid()
  if err != nil {
    return false
  } else {
    // check user...
    email, err := verifiedToken.GetClaim("email")
    if err == nil {
        userEmailAddress := fmt.Sptintf("%v", email)
    }

    // verified ok
    return true
  }
}
```

## Generate Service account oauth2 token

Firebase Service AccountからREST APIでFirebase DatabaseやFirebase Storageにアクセスする場合、ServiceAccountのOAuth2 Tokenが必要になります。

FirebaseのRead/Writeを行う場合、Scopeには　"https://www.googleapis.com/auth/firebase" "https://www.googleapis.com/auth/userinfo.email" が最低限必要となります。

```go

import "github.com/eaglesakura/gaefire"

func GetServiceAccountOAuth2Token(ctx context.Context, token string) string {
  serviceAccount := ... // factory.NewServiceAccount()

  token, err := serviceAccount.GetServiceAccountToken(ctx, "https://www.googleapis.com/auth/firebase", "https://www.googleapis.com/auth/userinfo.email")
  if err != nil {
    return ""
  } else {
    return token.AccessToken
  }
}
```

## Generate User oauth2 token

Firebase Service Account / Firebase Authのセットアップを行うと、`OAuth 2.0 クライアント ID` (ex, "Web client (auto created by Google Service)")が登録されます。
Google Play Service等の認証SDKと組み合わせることで、サーバーサイドでもユーザー権限のOAuth2 tokenを取得することができます。

自動されたWeb Clientを開き、`JSONをダウンロード` メニューから次のようなJSONを取得します。

```json
# web-application.json
{
  "web": {
    "client_id": "your-client-id.apps.googleusercontent.com",
    "project_id": "your-gcp-project-name",
    "auth_uri": "https://accounts.google.com/o/oauth2/auth",
    "token_uri": "https://accounts.google.com/o/oauth2/token",
    "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
    "client_secret": "your-client-secret",
    "redirect_uris": [
      "https://your-gcp-project-name.firebaseapp.com/__/auth/handler"
    ],
    "javascript_origins": [
      "http://localhost",
      "http://localhost:5000",
      "https://your-gcp-project-name.firebaseapp.com"
    ]
  }
}
```

例として、サーバーサイドでユーザー権限を用いたFirebase Read/WriteのOAuth2 Tokenが必要な場合を示します。

```java
// using Google Play Service
public static GoogleApiClient.Builder newFullPermissionClient(Context context) {
  GoogleSignInOptions options = new GoogleSignInOptions.Builder(GoogleSignInOptions.DEFAULT_SIGN_IN)
          .requestIdToken(context.getString(R.string.default_web_client_id))
          .requestEmail()
          .requestScopes(new Scope("https://www.googleapis.com/auth/firebase")) // Request scopes...
          .requestServerAuthCode("your-client-id.apps.googleusercontent.com", true)
          .build();
  return new GoogleApiClient.Builder(context)
          .addApi(Auth.GOOGLE_SIGN_IN_API, options)
          ;
}

@OnActivityResult(REQUEST_GOOGLE_AUTH)
void resultGoogleAuth(int result, Intent data) {
  GoogleSignInResult signInResult = Auth.GoogleSignInApi.getSignInResultFromIntent(data);
  if (signInResult.isSuccess()) {
      GoogleSignInAccount signInAccount = signInResult.getSignInAccount() ;

      String googleIdToken = signInAccount.getIdToken();
      String serverAuthCode = signInAccount.getServerAuthCode();
      // curl -H "Authorization: Bearer {token}"  https://your-gcp-project-name.appspot.com/path/to/api?serverAuthCode={serverAuthCode}
  } else {
    // login failed.
  }
}
```

```go

import "github.com/eaglesakura/gaefire"
import "github.com/eaglesakura/gaefire/factory"

func NewUserOAuth2Token(ctx context.Context, serverAuthCode string) (gaefire.OAuth2Token, error) {
  var webApplicationJson []byte
  webApplicationJson = ... // load from web-application.json

  webApplication := factory.NewWebApplication(webApplicationJson)
  return serviceAccount.NewUserAccountToken(ctx, serverAuthCode)
}
```

## Refresh oauth2 token

`NewUserAccountToken()`にて取得したユーザー権限のOAuth2 Tokenは3600 secでExpireされます。
ExpireされたOAuth2トークンはRefresh tokenを利用してリフレッシュしなければなりません。

GAE/Fireでは２通りの方法を用意しています。
AccessToken及びRefreshTokenはGAE/Fireはストレージへの保存を行いません。
それらはGAE/Go Appが適切な手段で保存を行ってください。

```go

import "github.com/eaglesakura/gaefire"

// Refresh Tokenから再度取得する
// この方法の場合、OAuth2 TokenはMemcacheにキャッシュされ、Expireされない限りは前回のトークンを再利用する
func GetUserAccountToken(ctx context.Context, refreshCode string) (gaefire.OAuth2Token, error) {
  webApplication := ...//  factory.NewWebApplication()
  return serviceAccount.GetUserAccountToken(ctx, refreshCode)
}
```

```go
import "github.com/eaglesakura/gaefire"

// OAuth2Token オブジェクトを直接利用する場合
func RefreshOAuth2Token(ctx context.Context, token *gaefire.OAuth2Token, clientId string, clientSecret string) error {
  return token.Refresh(ctx, clientId, clientSecret)
}
```

## Authentication from swagger.json(Cloud Endpoints 2.0 compat)

Cloud Endpoints 2.0の仕様に従って生成された `swagger.json(openapi.json)` を元に、認証を自動的に行うUtilを提供しています。
GAE/Fireが検証するのは次の `Authorization` ヘッダとAPI Keyです。

 * `Authorization: Bearer {google-id-token(Google sign-in)}`
 * `Authorization: Bearer {firebase-auth-json-web-token(Firebase SDK)}`
 * `Authorization: Bearer {json-web-token(GAE/Fire)}`
 * `Authorization: Bearer {google-oauth2-token}`
 * API Key(query|header)

 Http Requestに対して正しい認証が行われているかを確認することができます。
 例として、次のような認証はerrorとして扱います。

  * Authorization()を行う前に `X-Endpoint-API-UserInfo` Headerが存在している場合
  * `x-google-audiences` にて指定されないaudienceに対するOAuth2 Token
  * `x-google-issuer` に指定されないissuerが発行したJson Web Token
  * Service AccountもしくはGoogle以外が署名したjson-web-token
  * Google Cloud Platform Consoleに登録されないAPI Key
  * expireされているjson-web-token
  * expireされているOAuth2 Token

 下記はエラーとして扱いません。

  * API Keyを指定していない場合
  * `Authentication` ヘッダが存在しない場合


### Setup

事前にswagger.yaml等に下記の情報を追記します。
この仕様は [Google Cloud Endpoints 2.0](https://cloud.google.com/endpoints/) 及びSwaggerに従います。

また、Firebase Service Accountが "https://www.googleapis.com/auth/service.management" "https://www.googleapis.com/auth/servicecontrol" の権限を持つことを確認します。

```
host: "your-gcp-project-name.appspot.com"

securityDefinitions:
  api_key:
    type: "apiKey"
    name: "key"
    in: "query"
  google_id_token:
    type: oauth2
    authorizationUrl: ""
    flow: implicit
    x-google-issuer: "https://accounts.google.com"
    x-google-jwks-uri: "https://www.googleapis.com/oauth2/v1/certs"
    x-google-audiences:
      - "your-client-id-1.apps.googleusercontent.com"
      - "your-client-id-2.apps.googleusercontent.com"
  firebase:
    authorizationUrl: ""
    flow: "implicit"
    type: "oauth2"
    x-google-issuer: "https://securetoken.google.com/your-gcp-project-name"
    x-google-jwks_uri: "https://www.googleapis.com/service_accounts/v1/metadata/x509/securetoken@system.gserviceaccount.com"
```

作成したyaml/jsonは[gcloudコマンドでプロジェクトへDeploy](https://cloud.google.com/endpoints/docs/deploy-an-api)します。

```
$ gcloud service-management deploy path/to/openapi.yaml
```

### Verify API Key & Authorization header

 ```
 $ curl "Authentication: Bearer {your-token}" https://your-gcp-project-name.appspot.com/path/to/api?key={your-api-key}
 ```

```go
import "github.com/eaglesakura/gaefire"
import "github.com/eaglesakura/gaefire/factory"

func VerifyHttpRequest(ctx context.Context, request *http.Request) (*gaefire.AuthenticationInfo, error) {
  serviceAccountJson := // load service-account.json
  swaggerJson := // load swagger.json

  serviceAccount := factory.NewServiceAccount(serviceAccountJson)
  authenticationProxy := factory.NewAuthenticationProxy(serviceAccount, swaggerJson)

  authInfo, err := authenticationProxy.Authorization(ctx, request)
  if err == nil {
      if authInfo.User == nil {
        // Anonymous user.
      } else {
        // auth user
      }
  }

  return authInfo, err
}

```

# LICENSE

```
The MIT License (MIT)

Copyright (c) 2017 @eaglesakura

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
```
