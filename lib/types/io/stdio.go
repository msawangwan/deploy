package io

import "bytes"

type Stdio struct {
    Stdout bytes.Buffer
    Stderr bytes.Buffer
}
