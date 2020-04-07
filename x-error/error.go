package xerror

func DoIfErrorNotNil(err error,work func(err error )) {
	if err != nil {
		work(err)
	}
}
