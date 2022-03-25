package connection

type conn interface {
	Send()
}

type Linux struct{}

type Windows struct{}

func (l *Linux) SendToRemoteHost() {

}

func (w *Windows) SendToRemoteHost() {

}
