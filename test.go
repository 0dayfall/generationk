package main

type Common struct {
  Id string
}

type X struct {
  Common
  XProperty string
}

type Y struct {
  Common
  YProperty string
}

func (c *Common) AutoFill() {
  //set field on Common struct
}

func main() {
  x := &X{}
  x.AutoFill()
}
