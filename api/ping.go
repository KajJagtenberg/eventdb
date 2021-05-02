package api

func Ping(ctx *Ctx) error {
	ctx.Conn.WriteString("PONG")
	return nil
}
