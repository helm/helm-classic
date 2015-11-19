package action

type TestRunner struct {
	out []byte
	err error
}

func (r TestRunner) ClusterInfo() ([]byte, error) {
	return r.out, r.err
}

func (r TestRunner) Create(stdin []byte, ns string, dryRun bool) ([]byte, error) {
	return r.out, r.err
}

func (r TestRunner) Delete(name, ktype, ns string, dryRun bool) ([]byte, error) {
	return r.out, r.err
}

func (r TestRunner) Get(stdin []byte, ns string, dryRun bool) ([]byte, error) {
	return r.out, r.err
}
