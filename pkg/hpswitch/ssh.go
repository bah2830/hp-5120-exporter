package hpswitch

import "io"

func (s *Switch) getSSHResponse(cmd string) (string, error) {
	session, err := s.client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	outputBuf, err := session.StdoutPipe()
	if err != nil {
		return "", err
	}

	inputBuf, err := session.StdinPipe()
	if err != nil {
		return "", err
	}

	if err := session.Shell(); err != nil {
		return "", err
	}

	if _, err := inputBuf.Write([]byte(cmd + "\nquit\n")); err != nil {
		return "", err
	}

	var output string
	for {
		buf := make([]byte, 1024)
		size, err := outputBuf.Read(buf)
		if err != nil {
			if err == io.EOF {
				return output, nil
			}

			return "", err
		}

		data := string(buf[:size])
		output += data
	}

	return "", nil
}
