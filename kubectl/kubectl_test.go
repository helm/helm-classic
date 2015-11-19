package kubectl

type TestRunner struct {
	Runner

	out []byte
	err error
}

func (r TestRunner) Get(stdin []byte, ns string, dryRun bool) ([]byte, error) {
	return r.out, r.err
}
