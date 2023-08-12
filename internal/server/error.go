package server

type HttpError struct {
	// エラーが発生した行番号などを持たせておく
	Status int
	Err    error
}

// fmt.Println など、出力するための関数は、内部的に error型だった場合に Error を使う
// 行番号など、エラーメッセージに付け足すための処理を書いておく
func (e *HttpError) Error() string {
	return e.Err.Error()
}

// Unwrap を定義しておけば、errors.Unwrap が使える
func (e *HttpError) Unwrap() error {
	return e.Err
}
