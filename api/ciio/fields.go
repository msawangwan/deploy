package ciio

type Command struct {
    Exec string
    Args []string
}

type Addr struct {
    IP string
    PortOut string
    PortIn string
}
