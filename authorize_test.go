package gaefire


//// Firebase Storage link :: https://www.googleapis.com/storage/v1/b/${firebase-project}.appspot.com/o
//
//const (
//	_TEST_TOKEN_JSON = `
//		{
//			"type": "service_account",
//			"project_id": "example-project",
//			"private_key_id": "0123456789012345678901234567890123456789",
//			"private_key": "-----BEGIN RSA PRIVATE KEY-----\nMIIEpQIBAAKCAQEArX266ZevGyWIrQPdVfXi9nlYraPvnirUA3kPhVTIVDON20IZ\nMK1I1F8aixdhRzyvjsGVREsPp+ksdN5Tm8z7sMnl3mbVBhYqlIaTyOI1u5/efyuH\nHaeZ4+eNOUiYnZNyhKH+V5cL/hGOOF5YGo86Fd22OFWFBvi1TLIDfitSLXrjeE4V\nQDSYaCTrTaFYwEgRHq/aY26eW9Hf+A9v0ofGJN2rHgeeE2GeB0KwOU3EE1nNC4gB\nG9pOzQWdSTp6Az3cS8EzKnbWInaDu7dEFS3S4QPEhmtxCgwAJStok3Cqw4jTzHEY\nZxdWRbLsG3XWz5mK/FYeWtNsxLydTYMtKg73bwIDAQABAoIBAQCWXW9/f8D7lUdl\nNQFDvI5OsxOBw6Etg88nh2hFVhsep2QekaZFgxqpPhhCTyw30sgRwYQ+oZPbIwQt\n8neihEHsksBRRu0AjnuvKQizRiPuDvyHWdGMiTOezJSg/mOPxzis7b4EalyxgsOx\n5UsVEef1BakmIDDgvvRbmbnfQ2rBEODVLLxf/UpUqQ8AZ6oG+/Ny4KLpK03GeT3l\ncaAnocwNfqN4gyWHm4laZf74VYmYrlNG6mlT1gG+51e+hhhbPnAT11rL3XnOs9Xi\n+j7aeE66Y134rTnara3mbdqmPPHgWblptJgPGlzgEDjelsMCP7AVU4y8UnHsMK5t\nlP/cIGTpAoGBANk2A+UNTka73c/cFQjiym4lpNmwgKLEliO772HMskTzEqUUSQ6h\n+5y+G8H3FjDZAwKCXKMYYqS9gRudMXv/VPtXFLb1au/e/l2vEZ0ELtyrv2zvKtWx\nucQeFhTgyJGnYQ2T0mWiU91htvucUNH9twjjpNomTM8KIJDvgcS0hqWVAoGBAMx5\nBjiGQMJB0lC2uPOVC8MBLFQC0tZsPvAlwTrm5yad47MJgwFP/z8XC+2FxQkUWH62\nzu2nYlkm6tfdIC017ctclHEyGipc6tEO1RWrJZlbZBQcxS+NalSinFZQ5d64j/MF\n4EDOkZ3GDGdXdBTvzINBMruwmgo/zj1UZQztN9/zAoGAb6qRPgQlJcAXPHEMb1EI\neK/pm/BdcVBXT2+ilUjCrSe5ghx3ooor7FzfsEvyoJIwNe4G6eHzdHXoFeYuNm0B\np2URRS3OGBsv8cG68FniLZguBTa/crS3p9c/yuP0uMyv3GcOVymoq7s8cwXdltc0\nbeF6MpxWCGpQa7J1qEaWojECgYEAxb/OKmB8xOKPmov892aQR3oc+ur4KXPqsqpw\n5JxntUtB6ecrEdviSYvqdz7GPm+03mfCXMljLkGbIkWzVsYvQlw5G/iOoaXXW3Ry\n1E//Pv/KHEFu2vxzd4MEm94FUo9AeJKYPVUKM4JUgKVtmMoKCm7FuAumDn+C4IF8\ncTICtc0CgYEAiKjLufrP3XH3fmoRzHqknQpFpTrh33HOuYW/Wxy18tscV+uNyzuf\nHa9gV/672DNmCd5ts1nes5Cny0WQVOTsaAkEVw1O6JEDQw4sigIDyI9Ou2ZTzG4a\nWPqcNscjxv7e7CiywscHXGP6g7A25XsjMQgDPKdSu9UhMGMP2vs1jsU=\n-----END RSA PRIVATE KEY-----\n",
//			"client_email": "firebase-adminsdk@example-project.iam.gserviceaccount.com",
//			"client_id": "123456789012345678901",
//			"auth_uri": "https://accounts.google.com/o/oauth2/auth",
//			"token_uri": "https://accounts.google.com/o/oauth2/token",
//			"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
//			"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk%40example-project.iam.gserviceaccount.com"
//		}
//	`
//	_TEST_MOBILE_TOKEN = "" // replace check token.
//	_TEST_MOBILE_BROKEN_TOKEN = "" // replace check token.
//	_TEST_MOBILE_SERVER_CODE = "" // replace oauth server code(one time code)
//	_TEST_MOBILE_REFRESH_TOKEN = ""	// replace refresh token
//)
//
//func TestFirebaseAuthorizeService_NewToken(t *testing.T) {
//	serviceAccount := ServiceAccountJson{}
//
//	JSON := _TEST_TOKEN_JSON
//
//	if !assert.Nil(t, json.Unmarshal([]byte(JSON), &serviceAccount)) {
//		return
//	}
//
//	service, _ := NewServiceAccount(&serviceAccount)
//	token, err := service.NewToken("example-user").AddClaim("email", "example@example.com").Generate()
//	if !assert.Nil(t, err) {
//		return
//	}
//	if !assert.NotEqual(t, "", token) {
//		return
//	}
//
//	ioutil.WriteFile("generated-token.txt", []byte(token), os.ModePerm)
//}
//
//func TestPublicKeystore_Refresh(t *testing.T) {
//	testContext, dest, _ := aetest.NewContext()
//	defer dest()
//
//	// Google Keystoreは複数の鍵を持つが、時間経過で鍵が自動的に入れ替わる仕組みになっている
//	keystore := NewKeystore("securetoken@system.gserviceaccount.com")
//
//	if !assert.NotNil(t, keystore) {
//		return
//	}
//
//	if !assert.Nil(t, keystore.Refresh(testContext)) {
//		return
//	}
//}
//
//func TestFirebaseAuthorizeService_Verify(t *testing.T) {
//	if len(_TEST_MOBILE_TOKEN) == 0 {
//		return
//	}
//
//	testContext, dest, _ := aetest.NewContext()
//	defer dest()
//
//	keystore := NewKeystore("securetoken@system.gserviceaccount.com")
//	serviceAccount := ServiceAccountJson{}
//	if !assert.Nil(t, json.Unmarshal([]byte(_TEST_TOKEN_JSON), &serviceAccount)) {
//		return
//	}
//
//	service, _ := NewServiceAccount(&serviceAccount)
//	verifiedToken, err := service.Verify(testContext, _TEST_MOBILE_TOKEN, keystore, &VerifyOption{SkipExpireCheck:true})
//
//	if !assert.NotNil(t, verifiedToken) || !assert.Nil(t, err) {
//		log.Errorf(testContext, "Verify Error %v", err.Error())
//		return
//	}
//
//	if !assert.Equal(t, verifiedToken.GetClaim("email", ""), "email@example.com") {
//		return
//	}
//}
//
//func TestFirebaseAuthorizeService_VerifyFailed(t *testing.T) {
//	if len(_TEST_MOBILE_BROKEN_TOKEN) == 0 {
//		return
//	}
//
//	testContext, dest, _ := aetest.NewContext()
//	defer dest()
//
//	keystore := NewKeystore("securetoken@system.gserviceaccount.com")
//	serviceAccount := ServiceAccountJson{}
//	if !assert.Nil(t, json.Unmarshal([]byte(_TEST_TOKEN_JSON), &serviceAccount)) {
//		return
//	}
//
//	service, _ := NewServiceAccount(&serviceAccount)
//	verifiedToken, err := service.Verify(testContext, _TEST_MOBILE_BROKEN_TOKEN, keystore, &VerifyOption{SkipExpireCheck:true})
//
//	if !assert.Nil(t, verifiedToken) || !assert.NotNil(t, err) {
//		return
//	}
//}
//
//func TestFirebaseAuthorizeService_NewUserToken(t *testing.T) {
//	if _TEST_MOBILE_SERVER_CODE == "" {
//		// skip testing
//		return
//	}
//
//	testContext, delFunc, _ := aetest.NewContext()
//	defer delFunc()
//
//	serviceAccount := ServiceAccountJson{}
//	if !assert.Nil(t, json.Unmarshal([]byte(_TEST_TOKEN_JSON), &serviceAccount)) {
//		return
//	}
//
//	service, _ := NewServiceAccount(&serviceAccount)
//	token, err := service.NewUserAuthRequest(testContext,
//		WebApplicationInfo{ClientId:"your-client-id.apps.googleusercontent.com", ClientSecret:"your-client-secret"},
//		"test-user-email@example.com", _TEST_MOBILE_SERVER_CODE).
//		GetToken()
//
//	if !assert.Nil(t, err) {
//		return
//	}
//	if !assert.NotEqual(t, "", token.GetToken()) {
//		return
//	}
//	if !assert.NotEqual(t, "", token.GetRawToken().RefreshToken) {
//		return
//	}
//
//
//	// verify https://www.googleapis.com/oauth2/v3/tokeninfo?access_token=
//	ioutil.WriteFile("oauth2_initial_token.txt", []byte(token.GetToken()), os.ModePerm)
//	ioutil.WriteFile("oauth2_initial_jwt.txt", []byte(token.GetSignedToken()), os.ModePerm)
//	ioutil.WriteFile("oauth2_initial_refresh.txt", []byte(token.GetRefreshToken()), os.ModePerm)
//}
//
//func TestFirebaseAuthorizeService_RefreshUserToken(t *testing.T) {
//	if _TEST_MOBILE_REFRESH_TOKEN == "" {
//		// skip testing
//		return
//	}
//
//	testContext, delFunc, _ := aetest.NewContext()
//	defer delFunc()
//
//	serviceAccount := ServiceAccountJson{}
//	if !assert.Nil(t, json.Unmarshal([]byte(_TEST_TOKEN_JSON), &serviceAccount)) {
//		return
//	}
//
//	service, _ := NewServiceAccount(&serviceAccount)
//	token, err := service.NewUserTokenRefreshRequest(testContext,
//		WebApplicationInfo{ClientId:"your-client-id.apps.googleusercontent.com", ClientSecret:"your-client-secret"},
//		"test-user-email@example.com", _TEST_MOBILE_REFRESH_TOKEN).
//		GetToken()
//
//	if !assert.Nil(t, err) {
//		return
//	}
//	if !assert.NotEqual(t, "", token.GetToken()) {
//		return
//	}
//	if !assert.NotEqual(t, "", token.GetRawToken().RefreshToken) {
//		return
//	}
//
//
//	// verify https://www.googleapis.com/oauth2/v3/tokeninfo?access_token=
//	ioutil.WriteFile("oauth2_refreshed_token.txt", []byte(token.GetToken()), os.ModePerm)
//	ioutil.WriteFile("oauth2_refreshed_jwt.txt", []byte(token.GetSignedToken()), os.ModePerm)
//}
