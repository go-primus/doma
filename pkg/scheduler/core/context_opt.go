package core

type OptionFunc func(*taskCtx)

func WithProgressFunc(progressFunc ProgressFunc) OptionFunc {
	return func(tc *taskCtx) {
		tc.updateProgressFunc = progressFunc
	}
}

func WithSaveCheckPointFunc(saveCheckPointFunc SaveCheckPointFunc) OptionFunc {
	return func(tc *taskCtx) {
		tc.saveCheckPointFunc = saveCheckPointFunc
	}
}

func WithLoadCheckPointFunc(loadCheckPointFunc LoadCheckPointFunc) OptionFunc {
	return func(tc *taskCtx) {
		tc.loadCheckPointFunc = loadCheckPointFunc
	}
}
