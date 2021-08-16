package gsh

type Builtin func(ctx *Context, argv []string) error

func bbEcho(ctx *Context, argv []string) error {
	argv = argv[1:] // skip command name itself
	if len(argv) == 0 {
		// only print a '\n'
		_, err := ctxOut(ctx).Write([]byte{'\n'})
		return err
	}

	// if argv[0] starts with a '-' followed by any of [eEn] then it will affect options
	newline := true
	escapes := false

	if len(argv[0]) > 2 && argv[0][0] == '-' {
		valid := true
	loop:
		for _, c := range argv[0][1:] {
			switch c {
			case 'e':
				escapes = true
			case 'E':
				escapes = false
			case 'n':
				newline = false
			default:
				valid = false
				break loop
			}
		}
		if valid {
			// drop first argv
			argv = argv[1:]
		} else {
			// return to initial value
			newline = true
			escapes = false
		}
	}

	out := ctxOut(ctx)
	for _, arg := range argv {
		stop := false
		if escapes {
			arg, stop = handleEscapes(arg)
		}
		_, err := out.Write([]byte(arg))
		if err != nil {
			return err
		}
		if stop {
			return nil
		}
	}
	if newline {
		_, err := out.Write([]byte{'\n'})
		if err != nil {
			return err
		}
	}
	return nil
}
