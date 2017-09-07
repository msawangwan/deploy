package ci

func ProcessPushEvent() {
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
		err error
	)

	cmd := exec.Command("/bin/sh", "./bin/webhook", payload.Repository.FullName, payload.Ref)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	if err != nil {
		printErr(err, stderr.String())
	}

	log.Printf("command executed with result:\n%s\n", stdout.String())
}