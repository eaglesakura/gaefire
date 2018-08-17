package gaefire

/*
 * Firebaseの１ユーザーを示す。
 */
type FirebaseUser struct {
	/*
	 * 一意に示されるユーザーID
	 *
	 * 英数1-36文字以内
	 */
	UniqueId string
}
